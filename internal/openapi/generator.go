package openapi

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

const (
	FormatFloat    = "float"
	FormatDouble   = "double"
	FormatInt32    = "int32"
	FormatInt64    = "int64"
	FormatByte     = "byte"
	FormatDate     = "date"
	FormatDateTime = "date-time"
	FormatPassword = "password"
)

type Config struct {
	Package        string
	SchemaMetadata ast.SchemaMeta
}

type generator struct {
	schema *ast.Schema
}

func GenerateAST(filePath string, cfg Config) (*ast.Schema, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	oapi, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := oapi.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("[%s] %w", cfg.Package, err)
	}

	g := &generator{
		schema: &ast.Schema{
			Package:  cfg.Package,
			Metadata: cfg.SchemaMetadata,
		},
	}

	if oapi.Components == nil {
		return g.schema, nil
	}

	if err := g.declareDefinition(oapi.Components.Schemas); err != nil {
		return nil, fmt.Errorf("[%s] %w", cfg.Package, err)
	}

	return g.schema, nil
}

func (g *generator) declareDefinition(schemas openapi3.Schemas) error {
	for name, schemaRef := range schemas {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return err
		}

		g.schema.Objects = append(g.schema.Objects, ast.Object{
			Name:     name,
			Comments: schemaComments(schemaRef.Value),
			Type:     def,
			SelfRef: ast.RefType{
				ReferredPkg:  g.schema.Package,
				ReferredType: name,
			},
		})
	}

	sort.Slice(g.schema.Objects, func(i, j int) bool {
		return g.schema.Objects[i].Name < g.schema.Objects[j].Name
	})

	return nil
}

func (g *generator) walkSchemaRef(schemaRef *openapi3.SchemaRef) (ast.Type, error) {
	if isRef(schemaRef.Ref) {
		return g.walkRef(schemaRef)
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
	default:
		// No type defined
		return g.walkAny(schema)
	}
}

func (g *generator) walkRef(schema *openapi3.SchemaRef) (ast.Type, error) {
	pkg, referredKindName := g.getRefName(schema.Ref)

	return ast.NewRef(pkg, referredKindName), nil
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
		})
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	return ast.NewStruct(fields...), nil
}

func (g *generator) walkArray(schema *openapi3.Schema) (ast.Type, error) {
	def, err := g.walkSchemaRef(schema.Items)
	if err != nil {
		return ast.Type{}, err
	}

	return ast.NewArray(def), nil
}

func (g *generator) walkString(schema *openapi3.Schema) (ast.Type, error) {
	var t ast.Type
	switch schema.Format {
	case FormatDate, FormatDateTime, FormatPassword:
		t = ast.String()
	case FormatByte:
		t = ast.Bytes()
	default:
		t = ast.String()
	}

	t.Scalar.Constraints = getConstraints(schema)
	t.Nullable = schema.Nullable
	t.Default = schema.Default
	return t, nil
}

func (g *generator) walkNumber(schema *openapi3.Schema) (ast.Type, error) {
	var t ast.Type
	switch schema.Format {
	case FormatFloat:
		t = ast.NewScalar(ast.KindFloat32)
	case FormatDouble:
		t = ast.NewScalar(ast.KindFloat64)
	default:
		t = ast.NewScalar(ast.KindFloat32)
	}
	t.Scalar.Constraints = getConstraints(schema)
	t.Nullable = schema.Nullable
	t.Default = schema.Default
	return t, nil
}

func (g *generator) walkInteger(schema *openapi3.Schema) (ast.Type, error) {
	var t ast.Type
	switch schema.Format {
	case FormatInt32:
		t = ast.NewScalar(ast.KindInt32)
	case FormatInt64:
		t = ast.NewScalar(ast.KindInt64)
	default:
		t = ast.NewScalar(ast.KindInt64)
	}

	t.Scalar.Constraints = getConstraints(schema)
	t.Nullable = schema.Nullable
	t.Default = schema.Default
	return t, nil
}

func (g *generator) walkBoolean(_ *openapi3.Schema) (ast.Type, error) {
	return ast.Bool(), nil
}

func (g *generator) walkAny(_ *openapi3.Schema) (ast.Type, error) {
	return ast.Any(), nil
}

func (g *generator) walkAllOf(schema *openapi3.Schema) (ast.Type, error) {
	branches := make([]ast.Type, len(schema.AllOf))
	for i, sch := range schema.AllOf {
		def, err := g.walkSchemaRef(sch)
		if err != nil {
			return ast.Type{}, err
		}

		branches[i] = def
	}

	return ast.NewIntersection(branches), nil
}

func (g *generator) walkOneOf(schema *openapi3.Schema) (ast.Type, error) {
	discriminator, mapping := g.getDiscriminator(schema)
	return g.walkDisjunctions(schema.OneOf, discriminator, mapping)
}

func (g *generator) walkAnyOf(schema *openapi3.Schema) (ast.Type, error) {
	discriminator, mapping := g.getDiscriminator(schema)
	return g.walkDisjunctions(schema.AnyOf, discriminator, mapping)
}

func (g *generator) walkEnum(schema *openapi3.Schema) (ast.Type, error) {
	// Nullable enums? https://swagger.io/docs/specification/data-models/enums/
	enums := make([]ast.EnumValue, 0, len(schema.Enum))
	format := "%#v"
	if schema.Type == openapi3.TypeString {
		format = "%s"
	}

	enumType, err := getEnumType(schema.Type)
	if err != nil {
		return ast.Type{}, err
	}

	for _, value := range schema.Enum {
		enums = append(enums, ast.EnumValue{
			Type:  enumType,
			Name:  fmt.Sprintf(format, value),
			Value: value,
		})
	}

	return ast.NewEnum(enums, ast.Default(schema.Default)), nil
}

func (g *generator) walkDisjunctions(schemaRefs []*openapi3.SchemaRef, discriminator string, mapping map[string]string) (ast.Type, error) {
	typeDefs := make([]ast.Type, 0, len(schemaRefs))
	for _, schemaRef := range schemaRefs {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return ast.Type{}, err
		}

		typeDefs = append(typeDefs, def)
	}

	return ast.NewDisjunction(typeDefs, ast.Discriminator(discriminator, mapping)), nil
}

func (g *generator) getDiscriminator(schema *openapi3.Schema) (string, map[string]string) {
	name := ""
	mapping := make(map[string]string)
	if schema.Discriminator != nil {
		name = schema.Discriminator.PropertyName
		for k, v := range schema.Discriminator.Mapping {
			mapping[k] = v
		}
	}

	return name, mapping
}

func (g *generator) getRefName(value string) (string, string) {
	rgx := regexp.MustCompile(`(../)*(\w*/)*(.*).(json|yml)`)
	group := rgx.FindStringSubmatch(value)

	parts := strings.Split(value, "/")
	schemaName := parts[len(parts)-1]

	// Reference in the same file
	if len(group) == 0 {
		return g.schema.Package, schemaName
	}

	return group[3], schemaName
}
