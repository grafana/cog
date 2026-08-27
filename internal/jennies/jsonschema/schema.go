package jsonschema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type Definition = *orderedmap.Map[string, any]

type Schema struct {
	Config             Config
	ReferenceFormatter func(ref ir.RefType) string
	// OpenAPI3Compatible dictates whether the generated JSONSchema will be compatible with OpenAPI 3.0,
	// rather than JSONSchema (OpenAPI 3.1 is fully compatible with JSON Schema)
	OpenAPI3Compatible bool

	foreignObjects     *orderedmap.Map[string, ir.Object]
	referenceResolver  func(ref ir.RefType) (ir.Object, bool)
	isForeignReference func(ref ir.RefType) bool
}

func (jenny Schema) JennyName() string {
	return "JSONSchema"
}

func (jenny Schema) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	if jenny.ReferenceFormatter == nil {
		jenny.ReferenceFormatter = jenny.defaultRefFormatter
	}

	for _, schema := range context.Schemas {
		output, err := jenny.toJSON(jenny.GenerateSchema(context, schema))
		if err != nil {
			return nil, err
		}

		files = append(files, *codejen.NewFile(schema.Package+".jsonschema.json", output, jenny))
	}

	return files, nil
}

func (jenny Schema) toJSON(input any) ([]byte, error) {
	if !jenny.Config.Compact {
		return json.MarshalIndent(input, "", "  ")
	}

	return json.Marshal(input)
}

func (jenny Schema) GenerateSchema(context languages.Context, schema *ir.Schema) Definition {
	jenny.foreignObjects = orderedmap.New[string, ir.Object]()

	jenny.isForeignReference = func(ref ir.RefType) bool {
		return ref.ReferredPkg != schema.Package
	}
	jenny.referenceResolver = func(ref ir.RefType) (ir.Object, bool) {
		return context.GetObject(ref.ReferredPkg, ref.ReferredType)
	}

	jsonSchema := orderedmap.New[string, any]()
	jsonSchema.Set("$schema", "http://json-schema.org/draft-07/schema#")

	if schema.EntryPoint != "" {
		jsonSchema.Set("$ref", jenny.ReferenceFormatter(ir.RefType{
			ReferredPkg:  schema.Package,
			ReferredType: schema.EntryPoint,
		}))
	}

	definitions := orderedmap.New[string, Definition]()
	schema.Objects.Iterate(func(_ string, object ir.Object) {
		definitions.Set(object.Name, jenny.objectToDefinition(object))
	})

	for {
		if jenny.foreignObjects.Len() == 0 {
			break
		}

		foreignObjects := jenny.foreignObjects
		jenny.foreignObjects = orderedmap.New[string, ir.Object]()

		foreignObjects.Iterate(func(_ string, foreignObject ir.Object) {
			definitions.Set(foreignObject.Name, jenny.objectToDefinition(foreignObject))
		})
	}

	jsonSchema.Set("definitions", definitions)

	return jsonSchema
}

func (jenny Schema) objectToDefinition(object ir.Object) Definition {
	definition := jenny.formatType(object.Type)

	if comments := jenny.objectComments(object); len(comments) != 0 {
		definition.Set("description", comments)
	}
	if object.DeprecationMessage != "" {
		definition.Set("deprecated", true)
		definition.Set("x-deprecation-message", object.DeprecationMessage)
	}

	return definition
}

func (jenny Schema) formatType(typeDef ir.Type) Definition {
	var definition Definition
	switch typeDef.Kind {
	case ir.KindStruct:
		definition = jenny.formatStruct(typeDef)
	case ir.KindScalar:
		definition = jenny.formatScalar(typeDef)
	case ir.KindRef:
		definition = jenny.formatRef(typeDef)
	case ir.KindEnum:
		definition = jenny.formatEnum(typeDef)
	case ir.KindArray:
		definition = jenny.formatArray(typeDef)
	case ir.KindMap:
		definition = jenny.formatMap(typeDef)
	case ir.KindDisjunction:
		definition = jenny.formatDisjunction(typeDef)
	case ir.KindIntersection:
		definition = jenny.formatIntersection(typeDef)
	case ir.KindComposableSlot:
		definition = jenny.formatComposableSlot()
	case ir.KindConstantRef:
		definition = jenny.formatConstantRef(typeDef)
	default:
		definition = orderedmap.New[string, any]()
	}

	if typeDef.Nullable {
		definition = jenny.applyNullable(typeDef.Kind, definition)
	}

	return definition
}

// applyNullable wraps a definition to express "type | null".
// Ref-like kinds need special treatment in OpenAPI 3.0 because $ref replaces the whole object
// and cannot be combined with other keywords — those use allOf as a wrapper.
// All other kinds can carry nullable: true inline (OpenAPI 3.0) or use anyOf (JSON Schema).
func (jenny Schema) applyNullable(kind ir.Kind, definition Definition) Definition {
	nullDef := orderedmap.New[string, any]()
	nullDef.Set("type", "null")

	refLike := kind == ir.KindRef || kind == ir.KindConstantRef

	if jenny.OpenAPI3Compatible {
		nullable := orderedmap.New[string, any]()
		if refLike {
			nullable.Set("nullable", true)
			nullable.Set("allOf", []Definition{definition})
		} else {
			// Copy the existing definition and add nullable: true inline
			definition.Iterate(func(key string, value any) {
				nullable.Set(key, value)
			})
			nullable.Set("nullable", true)
		}
		return nullable
	}

	// JSON Schema: anyOf: [{...}, {type: "null"}]
	wrapper := orderedmap.New[string, any]()
	wrapper.Set("anyOf", []Definition{definition, nullDef})
	return wrapper
}

func (jenny Schema) formatScalar(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()

	switch typeDef.AsScalar().ScalarKind {
	case ir.KindNull:
		definition.Set("type", "null")
	case ir.KindAny:
		definition.Set("type", "object")
		definition.Set("additionalProperties", map[string]any{})
	case ir.KindBytes:
		definition.Set("type", "string")
		jenny.addStringConstraints(definition, typeDef)
	case ir.KindString:
		definition.Set("type", "string")
		jenny.addStringConstraints(definition, typeDef)
		if typeDef.HasHint(ir.HintStringFormatDateTime) {
			definition.Set("format", "date-time")
		}
	case ir.KindBool:
		definition.Set("type", "boolean")
	case ir.KindFloat32, ir.KindFloat64:
		definition.Set("type", "number")
		jenny.addNumberConstraints(definition, typeDef)
	case ir.KindUint8, ir.KindUint16, ir.KindUint32, ir.KindUint64,
		ir.KindInt8, ir.KindInt16, ir.KindInt32, ir.KindInt64:
		definition.Set("type", "integer")
		jenny.addNumberConstraints(definition, typeDef)
	}

	// constant value?
	if typeDef.AsScalar().IsConcrete() {
		definition.Set("const", typeDef.AsScalar().Value)
	}

	return definition
}

func (jenny Schema) addStringConstraints(definition *orderedmap.Map[string, any], typeDef ir.Type) {
	for _, constraint := range typeDef.AsScalar().Constraints {
		switch constraint.Op {
		case ir.MinLengthOp:
			definition.Set("minLength", constraint.Args[0])
		case ir.MaxLengthOp:
			definition.Set("maxLength", constraint.Args[0])
		case ir.RegexMatchOp:
			definition.Set("pattern", constraint.Args[0])
		case ir.NotRegexMatchOp:
			notDef := orderedmap.New[string, any]()
			notDef.Set("pattern", constraint.Args[0])
			definition.Set("not", notDef)
		}
	}
}

func (jenny Schema) addNumberConstraints(definition *orderedmap.Map[string, any], typeDef ir.Type) {
	for _, constraint := range typeDef.AsScalar().Constraints {
		switch constraint.Op {
		case ir.LessThanOp:
			if jenny.OpenAPI3Compatible {
				definition.Set("maximum", constraint.Args[0])
				definition.Set("exclusiveMaximum", true)
			} else {
				definition.Set("exclusiveMaximum", constraint.Args[0])
			}
		case ir.LessThanEqualOp:
			definition.Set("maximum", constraint.Args[0])
		case ir.GreaterThanOp:
			if jenny.OpenAPI3Compatible {
				definition.Set("minimum", constraint.Args[0])
				definition.Set("exclusiveMinimum", true)
			} else {
				definition.Set("exclusiveMinimum", constraint.Args[0])
			}
		case ir.GreaterThanEqualOp:
			definition.Set("minimum", constraint.Args[0])
		case ir.MultipleOfOp:
			definition.Set("multipleOf", constraint.Args[0])
		}
	}
}

func (jenny Schema) formatStruct(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()

	definition.Set("type", "object")
	definition.Set("additionalProperties", false)
	if typeDef.HasHint(ir.HintOpenStruct) {
		if val, _ := typeDef.Hints[ir.HintOpenStruct].(string); strings.ToLower(val)[0] == 't' {
			definition.Set("additionalProperties", map[string]any{})
		}
	}

	properties := orderedmap.New[string, any]()
	var required []string

	for _, field := range typeDef.AsStruct().Fields {
		fieldDef := jenny.formatType(field.Type)

		if comments := jenny.fieldComments(field); len(comments) != 0 {
			fieldDef.Set("description", comments)
		}

		properties.Set(field.Name, fieldDef)

		if field.Required {
			required = append(required, field.Name)
		}

		// TODO: review defaults management
		if field.Type.Default != nil {
			fieldDef.Set("default", field.Type.Default)
		}
	}

	if len(required) != 0 {
		definition.Set("required", required)
	}

	definition.Set("properties", properties)

	return definition
}

func (jenny Schema) formatRef(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()
	ref := typeDef.AsRef()

	if jenny.isForeignReference(ref) {
		referredObject, found := jenny.referenceResolver(ref)

		if found {
			jenny.foreignObjects.Set(referredObject.SelfRef.String(), referredObject)
		}
	}

	// TODO: handle foreign refs
	definition.Set("$ref", jenny.ReferenceFormatter(ref))

	return definition
}

func (jenny Schema) formatConstantRef(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()
	constRef := typeDef.AsConstantRef()
	ref := ir.NewRef(constRef.ReferredPkg, constRef.ReferredType).AsRef()

	if jenny.isForeignReference(ref) {
		referredObject, found := jenny.referenceResolver(ref)

		if found {
			jenny.foreignObjects.Set(referredObject.SelfRef.String(), referredObject)
		}
	}

	obj, ok := jenny.referenceResolver(ref)
	if !ok {
		definition.Set("$ref", jenny.ReferenceFormatter(ref))
		return definition
	}

	// TODO: handle foreign refs
	if obj.Type.IsEnum() {
		definition.Set("allOf", []Definition{jenny.formatType(ref.AsType())})
	} else {
		definition.Set("$ref", jenny.ReferenceFormatter(ref))
	}

	if obj, ok := jenny.referenceResolver(ref); ok && obj.Type.IsEnum() {
		definition.Set("default", constRef.ReferenceValue)
	}

	return definition
}

func (jenny Schema) defaultRefFormatter(ref ir.RefType) string {
	return fmt.Sprintf("#/definitions/%s", ref.ReferredType)
}

func (jenny Schema) formatEnum(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()

	values := tools.Map(typeDef.AsEnum().Values, func(value ir.EnumValue) any {
		return value.Value
	})

	definition.Set("enum", values)
	// Make an educated guess about the enum type by looking at the first element in the values set
	if len(typeDef.AsEnum().Values) > 0 {
		def := jenny.formatType(typeDef.AsEnum().Values[0].Type)
		definition.Set("type", def.Get("type"))
	}

	return definition
}

func (jenny Schema) formatArray(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()

	definition.Set("type", "array")
	definition.Set("items", jenny.formatType(typeDef.AsArray().ValueType))
	jenny.addArrayConstraints(definition, typeDef)

	return definition
}

func (jenny Schema) addArrayConstraints(definition *orderedmap.Map[string, any], typeDef ir.Type) {
	for _, constraint := range typeDef.AsArray().Constraints {
		switch constraint.Op {
		case ir.MinItemsOp:
			definition.Set("minItems", constraint.Args[0])
		case ir.MaxItemsOp:
			definition.Set("maxItems", constraint.Args[0])
		case ir.UniqueItemsOp:
			definition.Set("uniqueItems", true)
		}
	}
}

func (jenny Schema) formatMap(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()

	definition.Set("type", "object")
	valueType := typeDef.AsMap().ValueType
	if jenny.OpenAPI3Compatible && valueType.IsAny() {
		definition.Set("additionalProperties", true)
		return definition
	}

	definition.Set("additionalProperties", jenny.formatType(valueType))

	return definition
}

func (jenny Schema) formatDisjunction(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()

	branches := tools.UniqueFormattedBy(
		typeDef.AsDisjunction().Branches,
		jenny.formatType,
		func(d Definition) string {
			key, _ := json.Marshal(d)
			return string(key)
		},
	)

	definition.Set("oneOf", branches)

	return definition
}

func (jenny Schema) formatIntersection(typeDef ir.Type) Definition {
	definition := orderedmap.New[string, any]()
	branches := tools.Map(typeDef.AsIntersection().Branches, jenny.formatType)

	definition.Set("allOf", branches)

	return definition
}

func (jenny Schema) formatComposableSlot() Definition {
	definition := orderedmap.New[string, any]()

	// Same as "any"
	definition.Set("type", "object")
	definition.Set("additionalProperties", map[string]any{})

	return definition
}

func (jenny Schema) objectComments(object ir.Object) string {
	comments := object.Comments
	if jenny.Config.Debug {
		passesTrail := tools.Map(object.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	return strings.Join(comments, "\n")
}

func (jenny Schema) fieldComments(field ir.StructField) string {
	comments := field.Comments
	if jenny.Config.Debug {
		comments = append(comments, tools.Map(field.PassesTrail, jenny.passTrailFormatter)...)
		comments = append(comments, tools.Map(field.Type.PassesTrail, jenny.passTrailFormatter)...)
	}

	return strings.Join(comments, "\n")
}

func (jenny Schema) passTrailFormatter(trail string) string {
	return fmt.Sprintf("Modified by compiler pass '%s'", trail)
}
