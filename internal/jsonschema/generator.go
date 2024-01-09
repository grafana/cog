package jsonschema

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	schemaparser "github.com/santhosh-tekuri/jsonschema"
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
)

type Config struct {
	// Package name used to generate code into.
	Package string

	SchemaMetadata ast.SchemaMeta
}

type generator struct {
	schema *ast.Schema
}

func GenerateAST(schemaReader io.Reader, c Config) (*ast.Schema, error) {
	g := &generator{
		schema: &ast.Schema{
			Package:  c.Package,
			Metadata: c.SchemaMetadata,
		},
	}

	compiler := schemaparser.NewCompiler()
	compiler.ExtractAnnotations = true
	if err := compiler.AddResource("schema", schemaReader); err != nil {
		return nil, fmt.Errorf("[%s] %w", c.Package, err)
	}

	schema, err := compiler.Compile("schema")
	if err != nil {
		return nil, fmt.Errorf("[%s] %w", c.Package, err)
	}

	// The root of the schema is an actual type/object
	if schema.Ref == nil {
		if err := g.declareDefinition(c.Package, schema); err != nil {
			return nil, fmt.Errorf("[%s] %w", c.Package, err)
		}
	} else {
		definitionName := g.definitionNameFromRef(schema)

		// The root of the schema contains definitions, and a reference to the "main" object
		if err := g.declareDefinition(definitionName, schema.Ref); err != nil {
			return nil, fmt.Errorf("[%s] %w", c.Package, err)
		}
	}

	// To ensure consistent outputs
	sort.Slice(g.schema.Objects, func(i, j int) bool {
		return g.schema.Objects[i].Name < g.schema.Objects[j].Name
	})

	return g.schema, nil
}

func (g *generator) declareDefinition(definitionName string, schema *schemaparser.Schema) error {
	if _, found := g.schema.LocateObject(definitionName); found {
		return nil
	}

	def, err := g.walkDefinition(schema)
	if err != nil {
		return fmt.Errorf("%s: %w", definitionName, err)
	}

	g.schema.Objects = append(g.schema.Objects, ast.Object{
		Name: definitionName,
		Type: def,
		SelfRef: ast.RefType{
			ReferredPkg:  g.schema.Package,
			ReferredType: definitionName,
		},
	})

	return nil
}

func (g *generator) walkDefinition(schema *schemaparser.Schema) (ast.Type, error) {
	var def ast.Type
	var err error

	if len(schema.Types) == 0 {
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

		if schema.Properties != nil || schema.PatternProperties != nil {
			return g.walkObject(schema)
		}

		if schema.Enum != nil {
			return g.walkEnum(schema)
		}

		return ast.Type{}, errUndescriptiveSchema
	}

	//nolint: gocritic
	if len(schema.Types) > 1 {
		def, err = g.walkScalarDisjunction(schema.Types)
	} else if schema.Enum != nil {
		def, err = g.walkEnum(schema)
	} else {
		switch schema.Types[0] {
		case typeNull:
			def = ast.Null()
		case typeBoolean:
			def = ast.Bool()
		case typeString:
			def, err = g.walkString(schema)
		case typeObject:
			def, err = g.walkObject(schema)
		case typeNumber, typeInteger:
			def, err = g.walkNumber(schema)
		case typeArray:
			def, err = g.walkList(schema)
		default:
			return ast.Type{}, fmt.Errorf("unexpected schema with type '%s'", schema.Types[0])
		}
	}

	return def, err
}

func (g *generator) walkScalarDisjunction(types []string) (ast.Type, error) {
	branches := make([]ast.Type, 0, len(types))

	for _, typeName := range types {
		switch typeName {
		case typeNull:
			branches = append(branches, ast.Null())
		case typeBoolean:
			branches = append(branches, ast.Bool())
		case typeString:
			branches = append(branches, ast.String())
		case typeNumber, typeInteger:
			branches = append(branches, ast.NewScalar(ast.KindInt64))
		default:
			return ast.Type{}, fmt.Errorf("unexpected type in scalar disjunction '%s'", typeName)
		}
	}

	return ast.NewDisjunction(branches), nil
}

func (g *generator) walkDisjunctionBranches(branches []*schemaparser.Schema) ([]ast.Type, error) {
	definitions := make([]ast.Type, 0, len(branches))
	for _, oneOf := range branches {
		branch, err := g.walkDefinition(oneOf)
		if err != nil {
			return nil, err
		}

		definitions = append(definitions, branch)
	}

	return definitions, nil
}

func (g *generator) walkOneOf(schema *schemaparser.Schema) (ast.Type, error) {
	if len(schema.OneOf) == 0 {
		return ast.Type{}, fmt.Errorf("oneOf with no branches")
	}

	branches, err := g.walkDisjunctionBranches(schema.OneOf)
	if err != nil {
		return ast.Type{}, err
	}

	return ast.NewDisjunction(branches), nil
}

// TODO: what's the difference between oneOf and anyOf?
func (g *generator) walkAnyOf(schema *schemaparser.Schema) (ast.Type, error) {
	if len(schema.AnyOf) == 0 {
		return ast.Type{}, fmt.Errorf("anyOf with no branches")
	}

	branches, err := g.walkDisjunctionBranches(schema.AnyOf)
	if err != nil {
		return ast.Type{}, err
	}

	return ast.NewDisjunction(branches), nil
}

func (g *generator) walkAllOf(_ *schemaparser.Schema) (ast.Type, error) {
	// TODO: finish implementation and use correct type
	return ast.Type{}, nil
}

func (g *generator) definitionNameFromRef(schema *schemaparser.Schema) string {
	parts := strings.Split(schema.Ref.Ptr, "/")

	return parts[len(parts)-1] // Very naive
}

func (g *generator) walkRef(schema *schemaparser.Schema) (ast.Type, error) {
	referredKindName := g.definitionNameFromRef(schema)

	if err := g.declareDefinition(referredKindName, schema.Ref); err != nil {
		return ast.Type{}, err
	}

	// TODO: get the correct package for the referred type
	return ast.NewRef(g.schema.Package, referredKindName), nil
}

func (g *generator) walkString(_ *schemaparser.Schema) (ast.Type, error) {
	def := ast.String()

	/*
		if len(schema.Enum) != 0 {
			def.Constraints = append(def.Constraints, ast.TypeConstraint{
				Op:   "in",
				Args: []any{schema.Enum},
			})
		}
	*/

	return def, nil
}

func (g *generator) walkNumber(_ *schemaparser.Schema) (ast.Type, error) {
	// TODO: finish implementation
	return ast.NewScalar(ast.KindInt64), nil
}

func (g *generator) walkList(schema *schemaparser.Schema) (ast.Type, error) {
	var itemsDef ast.Type
	var err error

	if schema.Items == nil {
		itemsDef = ast.Any()
	} else {
		// TODO: schema.Items might not be a schema?
		itemsDef, err = g.walkDefinition(schema.Items.(*schemaparser.Schema))
		// items contains an empty schema: `{}`
		if errors.Is(err, errUndescriptiveSchema) {
			itemsDef = ast.Any()
		} else if err != nil {
			return ast.Type{}, err
		}
	}

	return ast.NewArray(itemsDef), nil
}

func (g *generator) walkEnum(schema *schemaparser.Schema) (ast.Type, error) {
	if len(schema.Enum) == 0 {
		return ast.Type{}, fmt.Errorf("enum with no values")
	}

	values := make([]ast.EnumValue, 0, len(schema.Enum))
	for _, enumValue := range schema.Enum {
		values = append(values, ast.EnumValue{
			Type: ast.String(), // TODO: identify that correctly

			// Simple mapping of all enum values (which we are assuming are in
			// lowerCamelCase) to corresponding CamelCase
			Name:  enumValue.(string),
			Value: enumValue.(string),
		})
	}

	return ast.NewEnum(values), nil
}

func (g *generator) walkObject(schema *schemaparser.Schema) (ast.Type, error) {
	// TODO: finish implementation
	fields := make([]ast.StructField, 0, len(schema.Properties))
	for name, property := range schema.Properties {
		fieldDef, err := g.walkDefinition(property)
		if err != nil {
			return ast.Type{}, err
		}

		fields = append(fields, ast.StructField{
			Name:     name,
			Comments: schemaComments(property),
			Required: tools.ItemInList(name, schema.Required),
			Type:     fieldDef,
		})
	}

	// To ensure consistent outputs
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	return ast.NewStruct(fields...), nil
}
