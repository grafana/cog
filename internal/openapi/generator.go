package openapi

import (
	"context"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"strings"
)

const (
	FormatFloat    = "float"
	FormatDouble   = "double"
	FormatInt32    = "int32"
	FormatInt64    = "int64"
	FormatString   = ""
	FormatByte     = "byte"
	FormatDate     = "date"
	FormatDateTime = "date-time"
	FormatPassword = "password"
)

type Config struct {
	Package string
}

type generator struct {
	file *ast.File
}

func GenerateAST(filePath string, cfg Config) (*ast.File, error) {
	loader := openapi3.NewLoader()
	oapi, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := oapi.Validate(context.Background()); err != nil {
		return nil, err
	}

	g := &generator{
		file: &ast.File{Package: cfg.Package},
	}

	if oapi.Components == nil {
		return g.file, nil
	}

	if err := g.declareDefinition(cfg.Package, oapi.Components.Schemas); err != nil {
		return nil, err
	}

	return g.file, nil
}

func (g *generator) declareDefinition(name string, schemas openapi3.Schemas) error {
	for _, schemaRef := range schemas {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return err
		}

		g.file.Definitions = append(g.file.Definitions, ast.Object{
			Name: name,
			Type: def,
		})
	}

	return nil
}

func (g *generator) walkSchemaRef(schemaRef *openapi3.SchemaRef) (ast.Type, error) {
	if schemaRef.Ref != "" {
		return g.walkRef(schemaRef.Ref)
	}

	return g.walkDefinitions(schemaRef.Value)
}

func (g *generator) walkDefinitions(schema *openapi3.Schema) (ast.Type, error) {
	if schema.AllOf != nil {
		return g.walkAllOf(schema)
	}
	if schema.AnyOf != nil {
		return g.walkAnyOf(schema)
	}
	if schema.OneOf != nil {
		return g.walkOneOf(schema)
	}
	if schema.Enum != nil {
		return g.walkEnum(schema)
	}
	if schema.Not != nil {
		return g.walkNot(schema)
	}

	switch schema.Type {
	case openapi3.TypeString:
		return g.walkString(schema)
	case openapi3.TypeBoolean:
		return g.walkBoolean(schema)
	case openapi3.TypeNumber:
		return g.walkNumber(schema)
	case openapi3.TypeInteger:
		return g.walkInteger(schema)
	case openapi3.TypeObject:
		return g.walkObject(schema)
	case openapi3.TypeArray:
		return g.walkArray(schema)
	}

	return ast.Type{}, nil
}

func (g *generator) walkRef(ref string) (ast.Type, error) {
	// TODO: Its assuming that we are referencing an object in the same document
	// See https://swagger.io/specification/v3/#reference-object
	parts := strings.Split(ref, "/")
	referredKindName := parts[len(parts)-1]
	return ast.NewRef(referredKindName), nil
}

func (g *generator) walkObject(schema *openapi3.Schema) (ast.Type, error) {
	fields := make([]ast.StructField, 0, len(schema.Properties))
	for name, schemaRef := range schema.Properties {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return ast.Type{}, err
		}
		fields = append(fields, ast.StructField{
			Name:     name,
			Comments: schemaComments(schema),
			Type:     def,
			Required: tools.ItemInList(name, schema.Required),
			Default:  schema.Default,
		})
	}
	return ast.NewStruct(fields), nil
}

func (g *generator) walkArray(schema *openapi3.Schema) (ast.Type, error) {
	def, err := g.walkSchemaRef(schema.Items)
	if err != nil {
		return ast.Type{}, err
	}

	return ast.NewArray(def), nil
}

func (g *generator) walkString(schema *openapi3.Schema) (ast.Type, error) {
	switch schema.Format {
	case FormatString, FormatDate, FormatDateTime, FormatPassword:
		return ast.String(), nil
	case FormatByte:
		return ast.Bytes(), nil
	default:
		return ast.Type{}, errors.New("unhandled string format")
	}
}

func (g *generator) walkNumber(schema *openapi3.Schema) (ast.Type, error) {
	switch schema.Format {
	case FormatFloat:
		return ast.NewScalar(ast.KindFloat32), nil
	case FormatDouble:
		return ast.NewScalar(ast.KindFloat64), nil
	default:
		return ast.Type{}, errors.New("unhandled number format")
	}
}

func (g *generator) walkInteger(schema *openapi3.Schema) (ast.Type, error) {
	switch schema.Format {
	case FormatInt32:
		return ast.NewScalar(ast.KindInt32), nil
	case FormatInt64:
		return ast.NewScalar(ast.KindInt64), nil
	default:
		return ast.Type{}, errors.New("unhandled integer format")
	}
}

func (g *generator) walkBoolean(_ *openapi3.Schema) (ast.Type, error) {
	return ast.Bool(), nil
}

func (g *generator) walkAllOf(schema *openapi3.Schema) (ast.Type, error) {
	return g.walkDisjunctions(schema.AllOf)
}

func (g *generator) walkOneOf(schema *openapi3.Schema) (ast.Type, error) {
	return g.walkDisjunctions(schema.OneOf)
}

func (g *generator) walkAnyOf(schema *openapi3.Schema) (ast.Type, error) {
	return g.walkDisjunctions(schema.AnyOf)
}

func (g *generator) walkEnum(schema *openapi3.Schema) (ast.Type, error) {
	// Nullable enums? https://swagger.io/docs/specification/data-models/enums/
	enums := make([]ast.EnumValue, 0, len(schema.Enum))
	for _, value := range schema.Enum {
		enumType, err := getEnumType(schema.Format)
		if err != nil {
			return ast.Type{}, err
		}

		enums = append(enums, ast.EnumValue{
			Type:  enumType,
			Name:  value.(string),
			Value: value,
		})
	}
	return ast.NewEnum(enums), nil
}

func (g *generator) walkDisjunctions(schemaRefs []*openapi3.SchemaRef) (ast.Type, error) {
	typeDefs := make([]ast.Type, 0, len(schemaRefs))
	for _, schemaRef := range schemaRefs {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return ast.Type{}, err
		}

		typeDefs = append(typeDefs, def)
	}
	return ast.NewDisjunction(typeDefs), nil
}

func (g *generator) walkNot(schema *openapi3.Schema) (ast.Type, error) {
	// TODO: Iterate not
	return ast.Type{}, nil
}
