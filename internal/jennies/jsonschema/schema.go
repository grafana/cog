package jsonschema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type Definition = *orderedmap.Map[string, any]

type Schema struct {
	Config             Config
	ReferenceFormatter func(ref ast.RefType) string
}

func (jenny Schema) JennyName() string {
	return "JSONSchema"
}

func (jenny Schema) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	if jenny.ReferenceFormatter == nil {
		jenny.ReferenceFormatter = jenny.defaultRefFormatter
	}

	for _, schema := range context.Schemas {
		output, err := json.Marshal(jenny.GenerateSchema(schema))
		if err != nil {
			return nil, err
		}

		files = append(files, *codejen.NewFile(schema.Package+".jsonschema.json", output, jenny))
	}

	return files, nil
}

func (jenny Schema) GenerateSchema(schema *ast.Schema) Definition {
	jsonSchema := orderedmap.New[string, any]()
	jsonSchema.Set("$schema", "http://json-schema.org/draft-07/schema#")

	if schema.EntryPoint != "" {
		jsonSchema.Set("$ref", jenny.ReferenceFormatter(ast.RefType{
			ReferredPkg:  schema.Package,
			ReferredType: schema.EntryPoint,
		}))
	}

	definitions := orderedmap.New[string, Definition]()

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		definitions.Set(object.Name, jenny.objectToDefinition(object))
	})

	jsonSchema.Set("definitions", definitions)

	return jsonSchema
}

func (jenny Schema) objectToDefinition(object ast.Object) Definition {
	definition := jenny.formatType(object.Type)

	if comments := jenny.objectComments(object); len(comments) != 0 {
		definition.Set("description", comments)
	}

	return definition
}

func (jenny Schema) formatType(typeDef ast.Type) Definition {
	switch typeDef.Kind {
	case ast.KindStruct:
		return jenny.formatStruct(typeDef)
	case ast.KindScalar:
		return jenny.formatScalar(typeDef)
	case ast.KindRef:
		return jenny.formatRef(typeDef)
	case ast.KindEnum:
		return jenny.formatEnum(typeDef)
	case ast.KindArray:
		return jenny.formatArray(typeDef)
	case ast.KindMap:
		return jenny.formatMap(typeDef)
	case ast.KindDisjunction:
		return jenny.formatDisjunction(typeDef)
	case ast.KindComposableSlot:
		return jenny.formatComposableSlot()
	}

	return orderedmap.New[string, any]()
}

func (jenny Schema) formatScalar(typeDef ast.Type) Definition {
	definition := orderedmap.New[string, any]()

	switch typeDef.AsScalar().ScalarKind {
	case ast.KindNull:
		definition.Set("type", "null")
	case ast.KindAny:
		definition.Set("type", "object")
		definition.Set("additionalProperties", map[string]any{})
	case ast.KindBytes, ast.KindString:
		definition.Set("type", "string")
		// TODO: constraints
	case ast.KindBool:
		definition.Set("type", "boolean")
	case ast.KindFloat32, ast.KindFloat64:
		definition.Set("type", "number")
		// TODO: constraints
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64,
		ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		definition.Set("type", "integer")
		// TODO: constraints
	}

	// constant value?
	if typeDef.AsScalar().IsConcrete() {
		definition.Set("const", typeDef.AsScalar().Value)
	}

	return definition
}

func (jenny Schema) formatStruct(typeDef ast.Type) Definition {
	definition := orderedmap.New[string, any]()

	definition.Set("type", "object")
	definition.Set("additionalProperties", false)

	properties := orderedmap.New[string, any]()
	var required []string

	for _, field := range typeDef.AsStruct().Fields {
		fieldDef := jenny.formatType(field.Type)

		// TODO: correctly handle passes trail

		if comments := strings.Join(field.Comments, "\n"); len(comments) != 0 {
			definition.Set("description", comments)
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

func (jenny Schema) formatRef(typeDef ast.Type) Definition {
	definition := orderedmap.New[string, any]()

	// TODO: handle foreign refs
	definition.Set("$ref", jenny.ReferenceFormatter(typeDef.AsRef()))

	return definition
}

func (jenny Schema) defaultRefFormatter(ref ast.RefType) string {
	return fmt.Sprintf("#/definitions/%s", ref.ReferredType)
}

func (jenny Schema) formatEnum(typeDef ast.Type) Definition {
	definition := orderedmap.New[string, any]()

	values := tools.Map(typeDef.AsEnum().Values, func(value ast.EnumValue) any {
		return value.Value
	})

	definition.Set("enum", values)

	return definition
}

func (jenny Schema) formatArray(typeDef ast.Type) Definition {
	definition := orderedmap.New[string, any]()

	definition.Set("type", "array")
	definition.Set("items", jenny.formatType(typeDef.AsArray().ValueType))

	return definition
}

func (jenny Schema) formatMap(typeDef ast.Type) Definition {
	definition := orderedmap.New[string, any]()

	definition.Set("type", "object")
	definition.Set("additionalProperties", jenny.formatType(typeDef.AsMap().ValueType))

	return definition
}

func (jenny Schema) formatDisjunction(typeDef ast.Type) Definition {
	definition := orderedmap.New[string, any]()
	branches := tools.Map(typeDef.AsDisjunction().Branches, jenny.formatType)

	definition.Set("anyOf", branches)

	return definition
}

func (jenny Schema) formatComposableSlot() Definition {
	definition := orderedmap.New[string, any]()

	// Same as "any"
	definition.Set("type", "object")
	definition.Set("additionalProperties", map[string]any{})

	return definition
}

func (jenny Schema) objectComments(object ast.Object) string {
	comments := object.Comments
	if jenny.Config.Debug {
		passesTrail := tools.Map(object.PassesTrail, func(trail string) string {
			return fmt.Sprintf("Modified by compiler pass '%s'", trail)
		})
		comments = append(comments, passesTrail...)
	}

	return strings.Join(comments, "\n")
}
