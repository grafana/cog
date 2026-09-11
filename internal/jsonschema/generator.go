package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/santhosh-tekuri/jsonschema/v6"
	schemaparser "github.com/santhosh-tekuri/jsonschema/v6"
)

var errUndescriptiveSchema = fmt.Errorf("the schema does not appear to be describing anything")

const (
	typeNull    = "null"
	typeBoolean = "boolean"
	typeObject  = "object"
	typeArray   = "array"
	typeString  = "string"
	typeNumber  = "number"
	typeInteger = "integer"

	formatDateTime = "date-time"
)

type Config struct {
	// Package name used to generate code into.
	Package string

	SchemaMetadata ir.SchemaMeta

	// Path to the schema file, if any.
	SchemaPath string
}

type generator struct {
	schema *ir.Schema
	seen   map[string]struct{}
}

func GenerateAST(schemaReader io.Reader, c Config) (*ir.Schema, error) {
	g := &generator{
		seen:   make(map[string]struct{}),
		schema: ir.NewSchema(c.Package, c.SchemaMetadata),
	}

	schemaResourceURL := "schema"
	if c.SchemaPath != "" {
		absSchemaPath, err := filepath.Abs(c.SchemaPath)
		if err != nil {
			return nil, fmt.Errorf("[%s] could not resolve base directory: %w", c.Package, err)
		}

		schemaResourceURL = "file://" + absSchemaPath
	}

	loader := jsonschema.SchemeURLLoader{
		"file":  jsonschema.FileLoader{},
		"http":  &HTTPURLLoader{},
		"https": &HTTPURLLoader{},
	}

	compiler := schemaparser.NewCompiler()
	compiler.UseLoader(loader)

	unmarshalledSchema, err := jsonschema.UnmarshalJSON(schemaReader)
	if err != nil {
		return nil, fmt.Errorf("[%s] %w", c.Package, err)
	}

	if err := compiler.AddResource(schemaResourceURL, unmarshalledSchema); err != nil {
		return nil, fmt.Errorf("[%s] %w", c.Package, err)
	}

	schema, err := compiler.Compile(schemaResourceURL)
	if err != nil {
		return nil, fmt.Errorf("[%s] %w", c.Package, err)
	}

	rootObjectName := c.Package

	// The root of the schema is an actual type/object
	if schema.Ref == nil {
		if err := g.declareDefinition(rootObjectName, schema); err != nil {
			return nil, fmt.Errorf("[%s] %w", c.Package, err)
		}
	} else {
		rootObjectName = g.definitionNameFromRef(schema)

		// The root of the schema contains definitions, and a reference to the "main" object
		if err := g.declareDefinition(rootObjectName, schema.Ref); err != nil {
			return nil, fmt.Errorf("[%s] %w", c.Package, err)
		}
	}

	if c.SchemaMetadata.Variant != "" {
		g.schema.Objects.Get(rootObjectName).Type.Hints[ir.HintImplementsVariant] = string(c.SchemaMetadata.Variant)
	}

	g.schema.EntryPoint = rootObjectName
	g.schema.EntryPointType = g.schema.Objects.Get(rootObjectName).SelfRef.AsType()

	// To ensure a consistent output, since github.com/santhosh-tekuri/jsonschema
	// doesn't guarantee the order of the definitions it parses.
	g.schema.Objects.Sort(orderedmap.SortStrings)

	return g.schema, nil
}

func (g *generator) declareDefinition(definitionName string, schema *schemaparser.Schema) error {
	if _, found := g.seen[definitionName]; found {
		return nil
	}

	g.seen[definitionName] = struct{}{}

	def, err := g.walkDefinition(schema)
	if err != nil {
		return fmt.Errorf("%s: %w", definitionName, err)
	}

	g.schema.AddObject(ir.Object{
		Name: definitionName,
		Type: def,
		SelfRef: ir.RefType{
			ReferredPkg:  g.schema.Package,
			ReferredType: definitionName,
		},
	})

	return nil
}

func (g *generator) walkDefinition(schema *schemaparser.Schema) (ir.Type, error) {
	var def ir.Type
	var err error

	if schema.Ref != nil {
		return g.walkRef(schema)
	}

	if schema.OneOf != nil {
		return g.walkOneOf(schema)
	}

	if schema.AnyOf != nil {
		return g.walkAnyOf(schema)
	}

	if schema.AllOf != nil {
		return g.walkAllOf(schema)
	}

	if schema.Enum != nil {
		return g.walkEnum(schema)
	}

	if schema.Types == nil || len(schema.Types.ToStrings()) == 0 {
		if schema.Properties != nil || schema.PatternProperties != nil || schema.AdditionalProperties != nil {
			return g.walkObject(schema)
		}

		if schema.Const != nil {
			return g.walkUntypedConstant(schema)
		}

		return ir.Any(), nil
	}

	// nolint: gocritic
	if len(schema.Types.ToStrings()) > 1 {
		def, err = g.walkScalarDisjunction(schema.Types.ToStrings())
	} else if schema.Enum != nil {
		def, err = g.walkEnum(schema)
	} else {
		switch schema.Types.ToStrings()[0] {
		case typeNull:
			def = ir.Null()
		case typeBoolean:
			def, err = g.walkBool(schema)
		case typeString:
			def, err = g.walkString(schema)
		case typeObject:
			def, err = g.walkObject(schema)
		case typeNumber, typeInteger:
			def, err = g.walkNumber(schema)
		case typeArray:
			def, err = g.walkList(schema)
		default:
			return ir.Type{}, fmt.Errorf("unexpected schema with type '%s'", schema.Types.String())
		}
	}

	return def, err
}

func (g *generator) walkScalarDisjunction(types []string) (ir.Type, error) {
	branches := make([]ir.Type, 0, len(types))

	for _, typeName := range types {
		switch typeName {
		case typeNull:
			branches = append(branches, ir.Null())
		case typeBoolean:
			branches = append(branches, ir.Bool())
		case typeString:
			branches = append(branches, ir.String())
		case typeNumber:
			branches = append(branches, ir.NewScalar(ir.KindFloat64))
		case typeInteger:
			branches = append(branches, ir.NewScalar(ir.KindInt64))
		default:
			return ir.Type{}, fmt.Errorf("unexpected type in scalar disjunction '%s'", typeName)
		}
	}

	return ir.NewDisjunction(branches), nil
}

func (g *generator) walkDisjunctionBranches(branches []*schemaparser.Schema) ([]ir.Type, error) {
	definitions := make([]ir.Type, 0, len(branches))
	for _, oneOf := range branches {
		branch, err := g.walkDefinition(oneOf)
		if err != nil {
			return nil, err
		}

		definitions = append(definitions, branch)
	}

	return definitions, nil
}

func (g *generator) walkUntypedConstant(schema *schemaparser.Schema) (ir.Type, error) {
	value := *schema.Const

	switch constant := value.(type) {
	case json.Number:
		if val, err := constant.Int64(); err == nil {
			return ir.NewScalar(ir.KindInt64, ir.Value(val)), nil
		} else if val, err := constant.Float64(); err == nil {
			return ir.NewScalar(ir.KindFloat64, ir.Value(val)), nil
		} else {
			return ir.Type{}, fmt.Errorf("could not parse json.Number %v", constant)
		}
	case bool:
		return ir.Bool(ir.Value(constant)), nil
	case string:
		return ir.String(ir.Value(constant)), nil
	case nil:
		return ir.Null(), nil
	default:
		return ir.Type{}, fmt.Errorf("unhandled constant type %T", value)
	}
}

func (g *generator) walkOneOf(schema *schemaparser.Schema) (ir.Type, error) {
	if len(schema.OneOf) == 0 {
		return ir.Type{}, fmt.Errorf("oneOf with no branches")
	}

	branches, err := g.walkDisjunctionBranches(schema.OneOf)
	if err != nil {
		return ir.Type{}, err
	}

	return ir.NewDisjunction(branches), nil
}

// TODO: what's the difference between oneOf and anyOf?
func (g *generator) walkAnyOf(schema *schemaparser.Schema) (ir.Type, error) {
	if len(schema.AnyOf) == 0 {
		return ir.Type{}, fmt.Errorf("anyOf with no branches")
	}

	branches, err := g.walkDisjunctionBranches(schema.AnyOf)
	if err != nil {
		return ir.Type{}, err
	}

	return ir.NewDisjunction(branches), nil
}

func (g *generator) walkAllOf(schema *schemaparser.Schema) (ir.Type, error) {
	branches := make([]ir.Type, len(schema.AllOf))
	for i, sch := range schema.AllOf {
		def, err := g.walkDefinition(sch)
		if err != nil {
			return ir.Type{}, err
		}

		branches[i] = def
	}

	if len(branches) == 1 {
		return branches[0], nil
	}

	return ir.NewIntersection(branches), nil
}

func (g *generator) definitionNameFromRef(schema *schemaparser.Schema) string {
	parts := strings.Split(schema.Ref.Location, "/")

	return parts[len(parts)-1] // Very naive
}

func (g *generator) walkRef(schema *schemaparser.Schema) (ir.Type, error) {
	referredKindName := g.definitionNameFromRef(schema)

	if err := g.declareDefinition(referredKindName, schema.Ref); err != nil {
		return ir.Type{}, err
	}

	// TODO: get the correct package for the referred type
	return ir.NewRef(g.schema.Package, referredKindName), nil
}

func (g *generator) walkString(schema *schemaparser.Schema) (ir.Type, error) {
	def := ir.String(ir.Default(maybeDefault(schema.Default)))

	if schema.Const != nil {
		def.Scalar.Value = *schema.Const
	}

	// to handle constant values defined as a string with a "static" regex:
	// ```
	// "someField": {
	// 	  "type": "string",
	// 	  "pattern": "^math$"
	// }
	// ```
	if schema.Pattern != nil && tools.RegexMatchesConstantString(schema.Pattern.String()) {
		def.Scalar.Value = tools.ConstantStringFromRegex(schema.Pattern.String())
	}

	if schema.Format != nil && schema.Format.Name == formatDateTime {
		def.Hints[ir.HintStringFormatDateTime] = true
	}

	if schema.MinLength != nil {
		def.Scalar.Constraints = append(def.Scalar.Constraints, ir.TypeConstraint{
			Op:   ir.MinLengthOp,
			Args: []any{*schema.MinLength},
		})
	}
	if schema.MaxLength != nil {
		def.Scalar.Constraints = append(def.Scalar.Constraints, ir.TypeConstraint{
			Op:   ir.MaxLengthOp,
			Args: []any{*schema.MaxLength},
		})
	}

	return def, nil
}

func (g *generator) walkBool(schema *schemaparser.Schema) (ir.Type, error) {
	def := ir.Bool(ir.Default(maybeDefault(schema.Default)))

	if schema.Const != nil {
		def.Scalar.Value = *schema.Const
	}

	return def, nil
}

func (g *generator) walkNumber(schema *schemaparser.Schema) (ir.Type, error) {
	scalarKind := ir.KindInt64
	if schema.Types.ToStrings()[0] == typeNumber {
		scalarKind = ir.KindFloat64
	}

	def := ir.NewScalar(scalarKind, ir.Default(maybeDefault(schema.Default)))

	if schema.Const != nil {
		def.Scalar.Value = unwrapJSONNumber(*schema.Const)
	}

	if schema.Minimum != nil {
		value, _ := schema.Minimum.Float64()
		def.Scalar.Constraints = append(def.Scalar.Constraints, ir.TypeConstraint{
			Op:   ir.GreaterThanEqualOp,
			Args: []any{value},
		})
	}
	if schema.ExclusiveMinimum != nil {
		value, _ := schema.ExclusiveMinimum.Float64()
		def.Scalar.Constraints = append(def.Scalar.Constraints, ir.TypeConstraint{
			Op:   ir.GreaterThanOp,
			Args: []any{value},
		})
	}
	if schema.Maximum != nil {
		value, _ := schema.Maximum.Float64()
		def.Scalar.Constraints = append(def.Scalar.Constraints, ir.TypeConstraint{
			Op:   ir.LessThanEqualOp,
			Args: []any{value},
		})
	}
	if schema.ExclusiveMaximum != nil {
		value, _ := schema.ExclusiveMaximum.Float64()
		def.Scalar.Constraints = append(def.Scalar.Constraints, ir.TypeConstraint{
			Op:   ir.LessThanOp,
			Args: []any{value},
		})
	}

	return def, nil
}

func (g *generator) walkList(schema *schemaparser.Schema) (ir.Type, error) {
	var itemsDef ir.Type
	var err error

	switch {
	case schema.Items == nil && schema.Items2020 == nil:
		itemsDef = ir.Any()
	case schema.Items2020 != nil:
		itemsDef, err = g.walkDefinition(schema.Items2020)
	default:
		// TODO: schema.Items might not be a schema?
		itemsDef, err = g.walkDefinition(schema.Items.(*schemaparser.Schema))
	}

	// items contains an empty schema: `{}`
	if errors.Is(err, errUndescriptiveSchema) {
		itemsDef = ir.Any()
	} else if err != nil {
		return ir.Type{}, err
	}

	return ir.NewArray(itemsDef, ir.Default(maybeDefault(schema.Default))), nil
}

func (g *generator) walkEnum(schema *schemaparser.Schema) (ir.Type, error) {
	if schema.Enum == nil || len(schema.Enum.Values) == 0 {
		return ir.Type{}, fmt.Errorf("enum with no values")
	}

	// we only want to deal with string or int enums
	enumType := ir.String()
	if _, ok := schema.Enum.Values[0].(string); !ok {
		enumType = ir.NewScalar(ir.KindInt64)
	}

	values := make([]ir.EnumValue, 0, len(schema.Enum.Values))
	for _, enumValue := range schema.Enum.Values {
		values = append(values, ir.EnumValue{
			Type:  enumType,
			Name:  fmt.Sprintf("%v", enumValue),
			Value: unwrapJSONNumber(enumValue),
		})
	}

	return ir.NewEnum(values), nil
}

func (g *generator) walkObject(schema *schemaparser.Schema) (ir.Type, error) {
	if len(schema.Properties) == 0 {
		// `schema.AdditionalProperties` is nil or false or *schemaparser.Schema
		_, ok := schema.AdditionalProperties.(bool)
		if schema.AdditionalProperties == nil || ok {
			return ir.Any(), nil
		}

		valueType, err := g.walkDefinition(schema.AdditionalProperties.(*schemaparser.Schema))
		if err != nil {
			return ir.Type{}, err
		}

		return ir.NewMap(ir.String(), valueType), nil
	}

	// TODO: finish implementation
	fields := make([]ir.StructField, 0, len(schema.Properties))
	for name, property := range schema.Properties {
		fieldDef, err := g.walkDefinition(property)
		if err != nil {
			return ir.Type{}, fmt.Errorf("%s: %w", name, err)
		}

		field := ir.NewStructField(name, fieldDef, ir.Comments(schemaComments(property)))
		field.Required = tools.ItemInList(name, schema.Required)

		fields = append(fields, field)
	}

	// To ensure consistent outputs
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	return ir.NewStruct(fields...), nil
}
