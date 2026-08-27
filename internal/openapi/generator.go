package openapi

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
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
	SchemaMetadata ir.SchemaMeta
	Validate       bool
}

type generator struct {
	schema *ir.Schema
}

func GenerateAST(ctx context.Context, oapi *openapi3.T, cfg Config) (*ir.Schema, error) {
	if cfg.Validate {
		if err := oapi.Validate(ctx, openapi3.DisableExamplesValidation()); err != nil {
			return nil, fmt.Errorf("[%s] %w", cfg.Package, err)
		}
	}

	g := &generator{
		schema: ir.NewSchema(cfg.Package, cfg.SchemaMetadata),
	}

	if oapi.Components == nil {
		return g.schema, nil
	}

	if err := g.declareDefinition(oapi.Components.Schemas); err != nil {
		return nil, fmt.Errorf("[%s] %w", cfg.Package, err)
	}

	// To ensure a consistent output, since github.com/getkin/kin-openapi/openapi3
	// doesn't guarantee the order of the definitions it parses.
	g.schema.Objects.Sort(orderedmap.SortStrings)

	return g.schema, nil
}

func (g *generator) declareDefinition(schemas openapi3.Schemas) error {
	for name, schemaRef := range schemas {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return err
		}

		g.schema.AddObject(ir.Object{
			Name:     name,
			Comments: schemaComments(schemaRef.Value),
			Type:     def,
			SelfRef: ir.RefType{
				ReferredPkg:  g.schema.Package,
				ReferredType: name,
			},
		})
	}

	return nil
}

func (g *generator) walkSchemaRef(schemaRef *openapi3.SchemaRef) (ir.Type, error) {
	if isRef(schemaRef.Ref) {
		return g.walkRef(schemaRef)
	}

	return g.walkDefinitions(schemaRef.Value)
}

func (g *generator) walkDefinitions(schema *openapi3.Schema) (ir.Type, error) {
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

	switch {
	case schema.Type.Is(openapi3.TypeString):
		return g.walkString(schema)
	case schema.Type.Is(openapi3.TypeObject):
		return g.walkObject(schema)
	case schema.Type.Is(openapi3.TypeArray):
		return g.walkArray(schema)
	case schema.Type.Is(openapi3.TypeBoolean):
		return g.walkBoolean(schema)
	case schema.Type.Is(openapi3.TypeInteger):
		return g.walkInteger(schema)
	case schema.Type.Is(openapi3.TypeNumber):
		return g.walkNumber(schema)
	}

	return g.walkAny(schema)
}

func (g *generator) walkRef(schema *openapi3.SchemaRef) (ir.Type, error) {
	pkg, referredKindName := g.getRefName(schema.Ref)

	return ir.NewRef(pkg, referredKindName), nil
}

func (g *generator) walkObject(schema *openapi3.Schema) (ir.Type, error) {
	if len(schema.Properties) == 0 {
		if schema.AdditionalProperties.Schema == nil {
			return ir.Any(), nil
		}

		valueType, err := g.walkSchemaRef(schema.AdditionalProperties.Schema)
		if err != nil {
			return ir.Type{}, err
		}

		return ir.NewMap(ir.String(), valueType), nil
	}

	fields := make([]ir.StructField, 0, len(schema.Properties))
	for name, schemaRef := range schema.Properties {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return ir.Type{}, err
		}

		field := ir.NewStructField(name, def, ir.Comments(schemaComments(schemaRef.Value)))
		field.Required = tools.ItemInList(name, schema.Required)

		fields = append(fields, field)
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	return ir.NewStruct(fields...), nil
}

func (g *generator) walkArray(schema *openapi3.Schema) (ir.Type, error) {
	def, err := g.walkSchemaRef(schema.Items)
	if err != nil {
		return ir.Type{}, err
	}

	return ir.NewArray(def, ir.Default(schema.Default)), nil
}

func (g *generator) walkString(schema *openapi3.Schema) (ir.Type, error) {
	var t ir.Type
	switch schema.Format {
	case FormatDateTime:
		t = ir.String(ir.Hints(ir.JenniesHints{
			ir.HintStringFormatDateTime: true,
		}))
	case FormatByte:
		t = ir.Bytes()
	default:
		t = ir.String()
	}

	if schema.Pattern != "" && tools.RegexMatchesConstantString(schema.Pattern) {
		t.Scalar.Value = tools.ConstantStringFromRegex(schema.Pattern)
	}

	t.Scalar.Constraints = getConstraints(schema)
	t.Nullable = schema.Nullable
	t.Default = schema.Default
	return t, nil
}

func (g *generator) walkNumber(schema *openapi3.Schema) (ir.Type, error) {
	var t ir.Type
	switch schema.Format {
	case FormatFloat:
		t = ir.NewScalar(ir.KindFloat32)
	case FormatDouble:
		t = ir.NewScalar(ir.KindFloat64)
	default:
		t = ir.NewScalar(ir.KindFloat32)
	}
	t.Scalar.Constraints = getConstraints(schema)
	t.Nullable = schema.Nullable
	t.Default = schema.Default
	return t, nil
}

func (g *generator) walkInteger(schema *openapi3.Schema) (ir.Type, error) {
	var t ir.Type
	switch schema.Format {
	case FormatInt32:
		t = ir.NewScalar(ir.KindInt32)
	case FormatInt64:
		t = ir.NewScalar(ir.KindInt64)
	default:
		t = ir.NewScalar(ir.KindInt64)
	}

	t.Scalar.Constraints = getConstraints(schema)
	t.Nullable = schema.Nullable
	t.Default = schema.Default
	return t, nil
}

func (g *generator) walkBoolean(schema *openapi3.Schema) (ir.Type, error) {
	return ir.Bool(ir.Default(schema.Default)), nil
}

func (g *generator) walkAny(_ *openapi3.Schema) (ir.Type, error) {
	return ir.Any(), nil
}

func (g *generator) walkAllOf(schema *openapi3.Schema) (ir.Type, error) {
	branches := make([]ir.Type, len(schema.AllOf))
	for i, sch := range schema.AllOf {
		def, err := g.walkSchemaRef(sch)
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

func (g *generator) walkOneOf(schema *openapi3.Schema) (ir.Type, error) {
	discriminator, mapping := g.getDiscriminator(schema)
	return g.walkDisjunctions(schema.OneOf, discriminator, mapping)
}

func (g *generator) walkAnyOf(schema *openapi3.Schema) (ir.Type, error) {
	discriminator, mapping := g.getDiscriminator(schema)
	return g.walkDisjunctions(schema.AnyOf, discriminator, mapping)
}

func (g *generator) walkEnum(schema *openapi3.Schema) (ir.Type, error) {
	// Nullable enums? https://swagger.io/docs/specification/data-models/enums/
	enums := make([]ir.EnumValue, 0, len(schema.Enum))
	format := "%#v"
	if schema.Type.Is(openapi3.TypeString) {
		format = "%s"
	}

	enumType, err := getEnumType(schema.Type.Slice()[0])
	if err != nil {
		return ir.Type{}, err
	}

	for _, value := range schema.Enum {
		enums = append(enums, ir.EnumValue{
			Type:  enumType,
			Name:  fmt.Sprintf(format, value),
			Value: value,
		})
	}

	return ir.NewEnum(enums, ir.Default(schema.Default)), nil
}

func (g *generator) walkDisjunctions(schemaRefs []*openapi3.SchemaRef, discriminator string, mapping map[string]string) (ir.Type, error) {
	typeDefs := make([]ir.Type, 0, len(schemaRefs))
	for _, schemaRef := range schemaRefs {
		def, err := g.walkSchemaRef(schemaRef)
		if err != nil {
			return ir.Type{}, err
		}

		typeDefs = append(typeDefs, def)
	}

	return ir.NewDisjunction(typeDefs, ir.Discriminator(discriminator, mapping)), nil
}

func (g *generator) getDiscriminator(schema *openapi3.Schema) (string, map[string]string) {
	name := ""
	mapping := make(map[string]string)
	if schema.Discriminator != nil {
		name = schema.Discriminator.PropertyName
		// schema.Discriminator.Mapping is a openapi3.StringMap[openapi3.MappingRef],
		// where openapi3.MappingRef is an alias for openapi3.SchemaRef that marshals to/from a string, using the `Ref` field.
		// As such, as can use v.Ref to get the value, instead of having to cast between []byte and string (and handle an always-nil error) with v.MarshalText()
		for k, v := range schema.Discriminator.Mapping {
			mapping[k] = v.Ref
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
