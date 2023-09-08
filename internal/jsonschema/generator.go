package jsonschema

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/grafana/cog/internal/ast"
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
}

type newGenerator struct {
	file *ast.File
}

func GenerateAST(schemaReader io.Reader, c Config) (*ast.File, error) {
	g := &newGenerator{
		file: &ast.File{
			Package: c.Package,
		},
	}

	compiler := schemaparser.NewCompiler()
	compiler.ExtractAnnotations = true
	if err := compiler.AddResource("schema", schemaReader); err != nil {
		return nil, err
	}

	schema, err := compiler.Compile("schema")
	if err != nil {
		return nil, err
	}

	// The root of the schema is an actual type/object
	if schema.Ref == nil {
		if err := g.declareDefinition(c.Package, schema); err != nil {
			return nil, err
		}
	} else {
		// The root of the schema contains definitions, and a reference to the "main" object
		if err := g.declareDefinition(c.Package, schema.Ref); err != nil {
			return nil, err
		}
	}

	return g.file, nil
}

func (g *newGenerator) declareDefinition(definitionName string, schema *schemaparser.Schema) error {
	def, err := g.walkDefinition(schema)
	if err != nil {
		return fmt.Errorf("%s: %w", definitionName, err)
	}

	g.file.Definitions = append(g.file.Definitions, ast.Object{
		Name: definitionName,
		Type: def,
	})

	return nil
}

func (g *newGenerator) walkDefinition(schema *schemaparser.Schema) (ast.Type, error) {
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
			return g.walkOneOf(schema)
		}

		if schema.Properties != nil || schema.PatternProperties != nil {
			return g.walkObject(schema)
		}

		if schema.Enum != nil {
			return g.walkEnum(schema)
		}

		return nil, errUndescriptiveSchema
	}

	if len(schema.Types) > 1 {
		def, err = g.walkDisjunction(schema)
	} else if schema.Enum != nil {
		def, err = g.walkEnum(schema)
	} else {
		switch schema.Types[0] {
		case typeNull:
			def = ast.ScalarType{ScalarKind: ast.KindNull}
		case typeBoolean:
			def = ast.ScalarType{ScalarKind: ast.KindBool}
		case typeString:
			def, err = g.walkString(schema)
		case typeObject:
			def, err = g.walkObject(schema)
		case typeNumber, typeInteger:
			def, err = g.walkNumber(schema)
		case typeArray:
			def, err = g.walkList(schema)

		default:
			return nil, fmt.Errorf("unexpected schema with type '%s'", schema.Types[0])
		}
	}

	return def, err
}

func (g *newGenerator) walkDisjunction(schema *schemaparser.Schema) (ast.DisjunctionType, error) {
	// TODO: finish implementation
	return ast.DisjunctionType{}, nil
}

func (g *newGenerator) walkDisjunctionBranches(branches []*schemaparser.Schema) ([]ast.Type, error) {
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

func (g *newGenerator) walkOneOf(schema *schemaparser.Schema) (ast.DisjunctionType, error) {
	if len(schema.OneOf) == 0 {
		return ast.DisjunctionType{}, fmt.Errorf("oneOf with no branches")
	}

	branches, err := g.walkDisjunctionBranches(schema.OneOf)
	if err != nil {
		return ast.DisjunctionType{}, err
	}

	return ast.DisjunctionType{
		Branches: branches,
	}, nil
}

// TODO: what's the difference between oneOf and anyOf?
func (g *newGenerator) walkAnyOf(schema *schemaparser.Schema) (ast.DisjunctionType, error) {
	if len(schema.AnyOf) == 0 {
		return ast.DisjunctionType{}, fmt.Errorf("anyOf with no branches")
	}

	branches, err := g.walkDisjunctionBranches(schema.AnyOf)
	if err != nil {
		return ast.DisjunctionType{}, err
	}

	return ast.DisjunctionType{
		Branches: branches,
	}, nil
}

func (g *newGenerator) walkAllOf(schema *schemaparser.Schema) (ast.DisjunctionType, error) {
	// TODO: finish implementation and use correct type
	return ast.DisjunctionType{}, nil
}

func (g *newGenerator) walkRef(schema *schemaparser.Schema) (ast.RefType, error) {
	parts := strings.Split(schema.Ref.Ptr, "/")
	referredKindName := parts[len(parts)-1] // Very naive

	if err := g.declareDefinition(referredKindName, schema.Ref); err != nil {
		return ast.RefType{}, err
	}

	return ast.RefType{
		ReferredType: referredKindName,
		//Comments: schemaComments(schema),
	}, nil
}

func (g *newGenerator) walkString(schema *schemaparser.Schema) (ast.ScalarType, error) {
	def := ast.ScalarType{ScalarKind: ast.KindString}

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

func (g *newGenerator) walkNumber(schema *schemaparser.Schema) (ast.ScalarType, error) {
	// TODO: finish implementation
	return ast.ScalarType{ScalarKind: ast.KindInt64}, nil
}

func (g *newGenerator) walkList(schema *schemaparser.Schema) (ast.ArrayType, error) {
	var itemsDef ast.Type
	var err error

	if schema.Items == nil {
		itemsDef = ast.ScalarType{
			ScalarKind: ast.KindAny,
		}
	} else {
		// TODO: schema.Items might not be a schema?
		itemsDef, err = g.walkDefinition(schema.Items.(*schemaparser.Schema))
		// items contains an empty schema: `{}`
		if errors.Is(err, errUndescriptiveSchema) {
			itemsDef = ast.ScalarType{
				ScalarKind: ast.KindAny,
			}
		} else if err != nil {
			return ast.ArrayType{}, err
		}
	}

	return ast.ArrayType{
		ValueType: itemsDef,
	}, nil
}

func (g *newGenerator) walkEnum(schema *schemaparser.Schema) (ast.EnumType, error) {
	if len(schema.Enum) == 0 {
		return ast.EnumType{}, fmt.Errorf("enum with no values")
	}

	values := make([]ast.EnumValue, 0, len(schema.Enum))
	for _, enumValue := range schema.Enum {
		values = append(values, ast.EnumValue{
			Type: ast.ScalarType{ScalarKind: ast.KindString}, // TODO: identify that correctly

			// Simple mapping of all enum values (which we are assuming are in
			// lowerCamelCase) to corresponding CamelCase
			Name:  enumValue.(string),
			Value: enumValue.(string),
		})
	}

	return ast.EnumType{
		Values: values,
		// TODO: default value?
	}, nil
}

func (g *newGenerator) walkObject(schema *schemaparser.Schema) (ast.StructType, error) {
	// TODO: finish implementation
	fields := make([]ast.StructField, 0, len(schema.Properties))
	for name, property := range schema.Properties {
		fieldDef, err := g.walkDefinition(property)
		if err != nil {
			return ast.StructType{}, err
		}

		fields = append(fields, ast.StructField{
			Name:     name,
			Comments: schemaComments(schema),
			Required: stringInList(schema.Required, name),
			Type:     fieldDef,
		})
	}

	return ast.StructType{
		Fields: fields,
	}, nil
}
