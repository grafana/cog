package ast

type Builder struct {
	Package         string
	For             Object
	Options         []Option
	Initializations []Assignment
}

type Builders []Builder

func (builders Builders) LocateByObject(pkg string, name string) (Builder, bool) {
	for _, builder := range builders {
		if builder.Package == pkg && builder.For.Name == name {
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
	IntoOptionalField bool
}

type BuilderGenerator struct {
}

func (generator *BuilderGenerator) FromAST(files []*File) []Builder {
	builders := make([]Builder, 0, len(files))

	for _, file := range files {
		for _, object := range file.Definitions {
			// we only want builders for structs
			if object.Type.Kind() != KindStruct {
				continue
			}

			builders = append(builders, generator.structObjectToBuilder(file, object))
		}
	}

	return builders
}

func (generator *BuilderGenerator) structObjectToBuilder(file *File, object Object) Builder {
	builder := Builder{
		Package: file.Package,
		For:     object,
		Options: nil,
	}
	structType := object.Type.(StructType)

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
	scalarType, ok := field.Type.(ScalarType)

	return ok && scalarType.Value != nil
}

func (generator *BuilderGenerator) structFieldToStaticInitialization(field StructField) Assignment {
	scalarType, _ := field.Type.(ScalarType)

	return Assignment{
		Path:              field.Name,
		Value:             scalarType.Value,
		ValueType:         field.Type,
		IntoOptionalField: !field.Required,
	}
}

func (generator *BuilderGenerator) structFieldToOption(field StructField) Option {
	var constraints []TypeConstraint
	if scalarType, ok := field.Type.(ScalarType); ok {
		constraints = scalarType.Constraints
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
				IntoOptionalField: !field.Required,
			},
		},
	}

	if field.Default != nil {
		opt.Default = &OptionDefault{
			ArgsValues: []any{field.Default},
		}
	}

	return opt
}
