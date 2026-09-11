package builders

import (
	"github.com/grafana/cog/pkg/ir"
)

type Generator struct {
}

func (generator *Generator) FromSchemas(schemas ir.Schemas) []ir.Builder {
	builders := make([]ir.Builder, 0, len(schemas))

	for _, schema := range schemas {
		schema.Objects.Iterate(func(_ string, object ir.Object) {
			resolvedType := schemas.Resolve(object.Type)
			if !resolvedType.IsAnyOf(ir.KindStruct, ir.KindRef) {
				return
			}

			builders = append(builders, generator.structObjectToBuilder(schemas, schema, object))
		})
	}

	return builders
}

func (generator *Generator) structObjectToBuilder(schemas ir.Schemas, schema *ir.Schema, object ir.Object) ir.Builder {
	builder := ir.Builder{
		Package:            schema.Package,
		For:                object,
		Name:               object.Name,
		DeprecationMessage: object.DeprecationMessage,
	}

	structType := schemas.Resolve(object.Type).AsStruct()
	for _, field := range structType.Fields {
		if field.Type.IsScalar() && field.Type.AsScalar().IsConcrete() {
			constantAssignment := ir.ConstantAssignment(ir.PathFromStructField(field), field.Type.AsScalar().Value)

			builder.Constructor.Assignments = append(builder.Constructor.Assignments, constantAssignment)
			continue
		}
		if field.Required && !field.Type.Nullable && generator.fieldIsRefToConcrete(schemas, field) {
			resolvedType := schemas.Resolve(field.Type)

			constantAssignment := ir.ConstantAssignment(ir.PathFromStructField(field), resolvedType.AsScalar().Value)
			builder.Constructor.Assignments = append(builder.Constructor.Assignments, constantAssignment)

			continue
		}
		if field.Required && !field.Type.Nullable && field.Type.IsConstantRef() {
			continue
		}

		builder.Options = append(builder.Options, generator.structFieldToOption(field))
	}

	return builder
}

func (generator *Generator) fieldIsRefToConcrete(schemas ir.Schemas, field ir.StructField) bool {
	if !field.Type.IsRef() {
		return false
	}

	return schemas.Resolve(field.Type).IsConcreteScalar()
}

func (generator *Generator) structFieldToOption(field ir.StructField) ir.Option {
	opt := ir.Option{
		Name:     field.Name,
		Comments: field.Comments,
		Args: []ir.Argument{
			{Name: field.Name, Type: field.Type},
		},
		Assignments: []ir.Assignment{
			ir.FieldAssignment(field),
		},
	}

	if field.Type.Default != nil {
		opt.Default = &ir.OptionDefault{
			ArgsValues: []any{field.Type.Default},
		}
	}

	return opt
}
