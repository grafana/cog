package ast

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	// Original data used to derive the builder, stored for read-only access
	// for the jennies and veneers.
	Schema *Schema
	For    Object

	// The builder itself
	// These fields are completely derived from the fields above and can be freely manipulated
	// by veneers.
	Package         string
	Name            string
	Properties      []StructField
	Options         []Option
	Initializations []Assignment `json:",omitempty"`
	VeneerTrail     []string     `json:",omitempty"`
}

func (builder *Builder) AddToVeneerTrail(veneerName string) {
	builder.VeneerTrail = append(builder.VeneerTrail, veneerName)
}

func (builder *Builder) MakePath(builders Builders, pathAsString string) (Path, error) {
	if pathAsString == "" {
		return nil, fmt.Errorf("can not make path from empty input")
	}

	resolveRef := func(ref RefType) (Builder, error) {
		referredObjBuilder, found := builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)
		if !found {
			return Builder{}, fmt.Errorf("could not make path '%s': reference '%s' could not be resolved", pathAsString, ref.String())
		}

		return referredObjBuilder, nil
	}

	currentType := builder.For.Type

	var path Path

	pathParts := strings.Split(pathAsString, ".")
	for _, part := range pathParts {
		if currentType.Kind == KindRef {
			referredObjBuilder, err := resolveRef(currentType.AsRef())
			if err != nil {
				return nil, err
			}

			currentType = referredObjBuilder.For.Type
		}

		if currentType.Kind != KindStruct {
			return nil, fmt.Errorf("could not make path '%s': type at path '%s' is not a struct or a ref", pathAsString, path.String())
		}

		field, found := currentType.AsStruct().FieldByName(part)
		if !found {
			return nil, fmt.Errorf("could not make path '%s': field '%s' not found under path '%s'", pathAsString, part, path.String())
		}

		path = append(path, PathItem{
			Identifier: part,
			Type:       field.Type,
		})

		currentType = field.Type
	}

	return path, nil
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
	Comments         []string `json:",omitempty"`
	VeneerTrail      []string `json:",omitempty"`
	Args             []Argument
	Assignments      []Assignment
	Default          *OptionDefault `json:",omitempty"`
	IsConstructorArg bool
}

func (opt *Option) AddToVeneerTrail(veneerName string) {
	opt.VeneerTrail = append(opt.VeneerTrail, veneerName)
}

type OptionDefault struct {
	ArgsValues []any
}

type Argument struct {
	Name string
	Type Type
}

type PathItem struct {
	Identifier string
	Type       Type // any
	// useful mostly for composability purposes, when a field Type is "any"
	// and we're trying to "compose in" something of a known type.
	TypeHint *Type `json:",omitempty"`
}

type Path []PathItem

func PathFromStructField(field StructField) Path {
	return Path{
		{
			Identifier: field.Name,
			Type:       field.Type,
		},
	}
}

func (path Path) Append(suffix Path) Path {
	var newPath Path
	newPath = append(newPath, path...)
	newPath = append(newPath, suffix...)

	return newPath
}

func (path Path) Last() PathItem {
	return path[len(path)-1]
}

func (path Path) String() string {
	return strings.Join(tools.Map(path, func(t PathItem) string {
		return t.Identifier
	}), ".")
}

type EnvelopeFieldValue struct {
	Path  Path            // where to assign within the struct/ref
	Value AssignmentValue // what to assign
}

type AssignmentEnvelope struct {
	Type   Type // Should be a ref or a struct only
	Values []EnvelopeFieldValue
}

type AssignmentValue struct {
	Argument *Argument           `json:",omitempty"`
	Constant any                 `json:",omitempty"`
	Envelope *AssignmentEnvelope `json:",omitempty"`
}

type AssignmentMethod string

const (
	DirectAssignment AssignmentMethod = "direct" // `foo = bar`
	AppendAssignment AssignmentMethod = "append" // `foo = append(foo, bar)`
)

type Assignment struct {
	// Where
	Path Path

	// What
	Value AssignmentValue

	// How
	Method AssignmentMethod

	Constraints []TypeConstraint `json:",omitempty"`
}

type AssignmentOpt func(assignment *Assignment)

func Constraints(constraints []TypeConstraint) AssignmentOpt {
	return func(assignment *Assignment) {
		assignment.Constraints = constraints
	}
}

func Method(method AssignmentMethod) AssignmentOpt {
	return func(assignment *Assignment) {
		assignment.Method = method
	}
}

func ConstantAssignment(path Path, value any, opts ...AssignmentOpt) Assignment {
	assignment := Assignment{
		Path: path,
		Value: AssignmentValue{
			Constant: value,
		},
		Method: DirectAssignment,
	}

	for _, opt := range opts {
		opt(&assignment)
	}

	return assignment
}

func ArgumentAssignment(path Path, argument Argument, opts ...AssignmentOpt) Assignment {
	assignment := Assignment{
		Path: path,
		Value: AssignmentValue{
			Argument: &argument,
		},
		Method: DirectAssignment,
	}

	for _, opt := range opts {
		opt(&assignment)
	}

	return assignment
}

func FieldAssignment(field StructField, opts ...AssignmentOpt) Assignment {
	var constraints []TypeConstraint
	if field.Type.Kind == KindScalar {
		constraints = field.Type.AsScalar().Constraints
	}

	argument := Argument{Name: field.Name, Type: field.Type}
	allOpts := []AssignmentOpt{Constraints(constraints)}
	allOpts = append(allOpts, opts...)

	return ArgumentAssignment(PathFromStructField(field), argument, allOpts...)
}

type BuilderGenerator struct {
}

func (generator *BuilderGenerator) FromAST(schemas Schemas) []Builder {
	builders := make([]Builder, 0, len(schemas))

	for _, schema := range schemas {
		for _, object := range schema.Objects {
			// we only want builders for structs or references to structs
			if object.Type.Kind == KindRef {
				ref := object.Type.AsRef()
				referredObj, found := schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)
				if !found {
					continue
				}

				if referredObj.Type.Kind != KindStruct {
					continue
				}
			}

			if object.Type.Kind != KindStruct && object.Type.Kind != KindRef {
				continue
			}

			builders = append(builders, generator.structObjectToBuilder(schemas, schema, object))
		}
	}

	return builders
}

func (generator *BuilderGenerator) structObjectToBuilder(schemas Schemas, schema *Schema, object Object) Builder {
	builder := Builder{
		Package: schema.Package,
		Schema:  schema,
		For:     object,
		Name:    object.Name,
	}

	var structType StructType
	if object.Type.Kind == KindStruct {
		structType = object.Type.AsStruct()
	} else {
		ref := object.Type.AsRef()
		referredObj, _ := schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)
		structType = referredObj.Type.AsStruct()
	}

	for _, field := range structType.Fields {
		if generator.fieldHasStaticValue(field) {
			constantAssignment := ConstantAssignment(PathFromStructField(field), field.Type.AsScalar().Value)
			builder.Initializations = append(builder.Initializations, constantAssignment)
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

	return field.Type.AsScalar().IsConcrete()
}

func (generator *BuilderGenerator) structFieldToOption(field StructField) Option {
	opt := Option{
		Name:     field.Name,
		Comments: field.Comments,
		Args: []Argument{
			{Name: field.Name, Type: field.Type},
		},
		Assignments: []Assignment{
			FieldAssignment(field),
		},
	}

	if field.Type.Default != nil {
		opt.Default = &OptionDefault{
			ArgsValues: []any{field.Type.Default},
		}
	}

	return opt
}
