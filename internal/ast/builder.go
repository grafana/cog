package ast

type Builder struct {
	// Original data used to derive the builder, stored for read-only access
	// for the jennies and veneers.
	Schema *Schema
	For    Object

	// The builder itself
	// These fields are completely derived from the fields above and can be freely manipulated
	// by veneers.
	RootPackage     string // ie: dashboard, alert, ... // TODO: better names and docs
	Package         string // ie: panel, link, ...
	Options         []Option
	Initializations []Assignment
}

type Builders []Builder

func (builders Builders) LocateByObject(pkg string, name string) (Builder, bool) {
	for _, builder := range builders {
		if builder.For.SelfRef.ReferredPkg == pkg && builder.For.SelfRef.ReferredType == name {
			return builder, true
		}
	}

	return Builder{}, false
}

type Option struct {
	Name             string
	Comments         []string
	Args             []Argument
	Assignments      []Assignment
	Default          *OptionDefault
	IsConstructorArg bool
}

type OptionDefault struct {
	ArgsValues []any
}

type Argument struct {
	Name string
	Type Type
}

type Assignment struct {
	// Where
	Path string

	// What
	ValueType    Type   // type of the value being assigned
	ArgumentName string // if empty, then use `Value`
	Value        any

	Constraints []TypeConstraint

	// Some more context on the what
	IntoNullableField bool
}

type BuilderGenerator struct {
}

func (generator *BuilderGenerator) FromAST(schemas []*Schema) []Builder {
	builders := make([]Builder, 0, len(schemas))

	for _, schema := range schemas {
		for _, object := range schema.Objects {
			// we only want builders for structs
			if object.Type.Kind != KindStruct {
				continue
			}

			builders = append(builders, generator.structObjectToBuilder(schema, object))
		}
	}

	return builders
}

func (generator *BuilderGenerator) structObjectToBuilder(schema *Schema, object Object) Builder {
	builder := Builder{
		RootPackage: schema.Package,
		Package:     object.Name,
		Schema:      schema,
		For:         object,
	}
	structType := object.Type.AsStruct()

	for _, field := range structType.Fields {
		if generator.fieldHasStaticValue(field) {
			builder.Initializations = append(builder.Initializations, generator.structFieldToStaticInitialization(field))

			continue
		}

		builder.Options = append(builder.Options, generator.structFieldToOption(field))
	}

	return builder
}

func (generator *BuilderGenerator) fieldHasStaticValue(field StructField) bool {
	if field.Type.Kind != KindScalar {
		return false
	}

	return field.Type.AsScalar().Value != nil
}

func (generator *BuilderGenerator) structFieldToStaticInitialization(field StructField) Assignment {
	return Assignment{
		Path:              field.Name,
		Value:             field.Type.AsScalar().Value,
		ValueType:         field.Type,
		IntoNullableField: field.Type.Nullable,
	}
}

func (generator *BuilderGenerator) structFieldToOption(field StructField) Option {
	var constraints []TypeConstraint
	if field.Type.Kind == KindScalar {
		constraints = field.Type.AsScalar().Constraints
	}

	opt := Option{
		Name:     field.Name,
		Comments: field.Comments,
		Args: []Argument{
			{
				Name: field.Name,
				Type: field.Type,
			},
		},
		Assignments: []Assignment{
			{
				Path:              field.Name,
				ArgumentName:      field.Name,
				ValueType:         field.Type,
				Constraints:       constraints,
				IntoNullableField: field.Type.Nullable,
			},
		},
	}

	if field.Type.Default != nil {
		opt.Default = &OptionDefault{
			ArgsValues: []any{field.Type.Default},
		}
	}

	return opt
}
