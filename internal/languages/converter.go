package languages

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type MappingGuard struct {
	Path  ast.Path
	Op    ast.Op
	Value any
}

func (guard MappingGuard) String() string {
	return fmt.Sprintf("%s %s %v", guard.Path, guard.Op, guard.Value)
}

type DirectArgMapping struct {
	ValuePath ast.Path // direct mapping between a JSON value and an argument
	ValueType ast.Type
}

// mapping between a JSON value and an argument delegated to a builder
type BuilderArgMapping struct {
	ValuePath   ast.Path
	ValueType   ast.Type
	BuilderPkg  string
	BuilderName string
}

type ArrayArgMapping struct {
	For       ast.Path
	ForType   ast.Type
	ForArg    *ArgumentMapping
	ValueAs   ast.Path
	ValueType ast.Type
}

type RuntimeArgMapping struct {
	FuncName string
	Args     []*DirectArgMapping
}

type ArgumentMapping struct {
	Direct  *DirectArgMapping
	Runtime *RuntimeArgMapping
	Builder *BuilderArgMapping
	Array   *ArrayArgMapping

	Guards []MappingGuard
}

type OptionMapping struct {
	Option ast.Option // option in the builder

	// for _, panel := range input.Panels { WithPanel(panel) }
	RepeatFor ast.Path // `input.Panels`
	RepeatAs  string   // `panel`

	Guards []MappingGuard
	Args   []ArgumentMapping
}

type ConverterInput struct {
	ArgName string
	TypeRef ast.RefType
}

type Converter struct {
	Package string

	BuilderName string
	Input       ConverterInput

	// FIXME: assuming we only have direct mappings here is... *optimistic*.
	ConstructorArgs []DirectArgMapping

	Mappings []OptionMapping
}

func (converter Converter) inputRootPath() ast.Path {
	return ast.Path{
		{
			Identifier: converter.Input.ArgName,
			Type:       converter.Input.TypeRef.AsType(),
			Root:       true,
		},
	}
}

type ConverterGenerator struct {
}

func (generator *ConverterGenerator) FromBuilder(context Context, builder ast.Builder) Converter {
	converter := Converter{
		Package:     builder.Package,
		BuilderName: builder.Name,

		Input: ConverterInput{
			ArgName: "input",
			TypeRef: builder.For.SelfRef,
		},
	}

	converter.ConstructorArgs = generator.constructorArgs(converter, builder)

	converter.Mappings = tools.Map(builder.Options, func(option ast.Option) OptionMapping {
		return generator.convertOption(context, converter, option)
	})

	return converter
}

func (generator *ConverterGenerator) constructorArgs(converter Converter, builder ast.Builder) []DirectArgMapping {
	// we're only interested in assignments made from a constructor argument
	// (as opposed to constant initializations for example)
	argAssignments := tools.Filter(builder.Constructor.Assignments, func(assignment ast.Assignment) bool {
		return assignment.Value.Argument != nil
	})

	return tools.Map(argAssignments, func(assignment ast.Assignment) DirectArgMapping {
		return DirectArgMapping{
			ValuePath: converter.inputRootPath().Append(assignment.Path),
			ValueType: assignment.Path.Last().Type,
		}
	})
}

func (generator *ConverterGenerator) convertOption(context Context, converter Converter, option ast.Option) OptionMapping {
	mapping := OptionMapping{
		Option: option,
		Guards: generator.optionGuards(converter, option),
	}

	nonConstantAssignments := tools.Filter(option.Assignments, func(assignment ast.Assignment) bool {
		return assignment.Value.Constant == nil
	})

	if len(nonConstantAssignments) == 1 && nonConstantAssignments[0].Method == ast.AppendAssignment {
		mapping.RepeatFor = converter.inputRootPath().Append(nonConstantAssignments[0].Path)
		mapping.RepeatAs = "item"
	}

	i := 0
	for _, assignment := range nonConstantAssignments {
		i++

		argName := fmt.Sprintf("arg%d", i)
		valueType := assignment.Path.Last().Type
		valuePath := converter.inputRootPath().Append(assignment.Path)
		if mapping.RepeatFor != nil {
			valueType = valueType.AsArray().ValueType
			valuePath = ast.Path{
				{Identifier: mapping.RepeatAs, Type: valueType, Root: true},
			}
		}

		if argument, ok := generator.argumentFromDisjunctionStruct(context, converter, argName, valuePath, assignment); ok {
			mapping.Args = append(mapping.Args, argument)
			continue
		}

		if assignment.Value.Envelope != nil {
			arguments := generator.argumentsForEnvelope(context, converter, argName, valuePath, assignment)
			mapping.Args = append(mapping.Args, arguments...)
			continue
		}

		argument := generator.argumentForType(context, converter, argName, valuePath, valueType)
		mapping.Args = append(mapping.Args, argument)
	}

	return mapping
}

func (generator *ConverterGenerator) argumentsForEnvelope(context Context, converter Converter, argName string, valuePath ast.Path, assignment ast.Assignment) []ArgumentMapping {
	var mappings []ArgumentMapping
	for _, envelopeField := range assignment.Value.Envelope.Values {
		fieldValuePath := valuePath.Append(envelopeField.Path)
		mappings = append(mappings, generator.argumentForType(context, converter, argName, fieldValuePath, fieldValuePath.Last().Type))
	}
	return mappings
}

func (generator *ConverterGenerator) argumentForType(context Context, converter Converter, argName string, valuePath ast.Path, typeDef ast.Type) ArgumentMapping {
	if typeDef.IsComposableSlot() {
		var slotTypeHintsArgs []*DirectArgMapping
		var guards []MappingGuard
		converterRootType := context.ResolveRefs(converter.inputRootPath().Last().Type)

		if converterRootType.IsStruct() {
			for _, field := range converterRootType.Struct.Fields {
				if !field.Type.IsRef() || field.Type.AsRef().ReferredType != "DataSourceRef" {
					continue
				}

				refPath := ast.PathFromStructField(field)
				datasourceRefType := context.ResolveRefs(field.Type)

				typeField, found := datasourceRefType.Struct.FieldByName("type")
				if !found {
					continue
				}

				typePath := refPath.AppendStructField(typeField)
				guards = append(guards, generator.pathNotNullGuards(converter, typePath)...)

				slotTypeHintsArgs = append(slotTypeHintsArgs, &DirectArgMapping{
					ValuePath: converter.inputRootPath().Append(typePath),
					ValueType: typeField.Type,
				})
			}
		}

		return ArgumentMapping{
			Runtime: &RuntimeArgMapping{
				FuncName: fmt.Sprintf("Convert%sToGo", tools.UpperCamelCase(string(typeDef.AsComposableSlot().Variant))),
				Args: append([]*DirectArgMapping{
					{ValuePath: valuePath, ValueType: typeDef},
				}, slotTypeHintsArgs...),
			},
			Guards: guards,
		}
	}

	if typeDef.IsArray() {
		valueAs := ast.Path{
			{Identifier: argName, Type: typeDef.Array.ValueType, Root: true},
		}

		forArg := generator.argumentForType(context, converter, argName+"Value", valueAs, typeDef.Array.ValueType)

		return ArgumentMapping{
			Array: &ArrayArgMapping{
				For:       valuePath,
				ForType:   typeDef,
				ForArg:    &forArg,
				ValueAs:   valueAs,
				ValueType: typeDef.Array.ValueType,
			},
		}
	}

	builder, found := context.ResolveAsBuilder(typeDef)
	if found && strings.EqualFold(builder.Package, "dashboard") && strings.EqualFold("panel", builder.For.Name) {
		typeField, _ := builder.For.Type.Struct.FieldByName("type")

		return ArgumentMapping{
			Runtime: &RuntimeArgMapping{
				FuncName: "ConvertPanelToGo",
				Args: []*DirectArgMapping{
					{ValuePath: valuePath, ValueType: typeDef},
					{ValuePath: valuePath.AppendStructField(typeField), ValueType: typeField.Type},
				},
			},
		}
	}
	if found {
		return ArgumentMapping{
			Builder: &BuilderArgMapping{
				ValuePath:   valuePath,
				ValueType:   typeDef,
				BuilderPkg:  builder.Package,
				BuilderName: builder.Name,
			},
		}
	}

	return ArgumentMapping{
		Direct: &DirectArgMapping{
			ValuePath: valuePath,
			ValueType: typeDef,
		},
	}
}

func (generator *ConverterGenerator) argumentFromDisjunctionStruct(context Context, converter Converter, argName string, valuePath ast.Path, assignment ast.Assignment) (ArgumentMapping, bool) {
	if assignment.Value.Envelope == nil {
		return ArgumentMapping{}, false
	}

	envelopedType := assignment.Value.Envelope.Type
	if envelopedType.IsRef() {
		referredObject, _ := context.LocateObject(envelopedType.Ref.ReferredPkg, envelopedType.Ref.ReferredType)
		envelopedType = referredObject.Type
	}

	if !envelopedType.IsStructGeneratedFromDisjunction() {
		return ArgumentMapping{}, false
	}

	envelopeValues := assignment.Value.Envelope.Values

	arg := generator.argumentForType(context, converter, argName, valuePath.Append(envelopeValues[0].Path), envelopeValues[0].Path.Last().Type)
	arg.Guards = tools.Map(envelopeValues, func(envelopedField ast.EnvelopeFieldValue) MappingGuard {
		return MappingGuard{
			Path:  valuePath.Append(envelopedField.Path),
			Op:    ast.NotEqualOp,
			Value: nil,
		}
	})

	return arg, true
}

func (generator *ConverterGenerator) optionGuards(converter Converter, option ast.Option) []MappingGuard {
	// conditions safeguarding the conversion of the current option
	guards := orderedmap.New[string, MappingGuard]()

	// TODO: define guards other than "not null" checks? (0, "", ...)
	// TODO: builders + array of builders (and array of array of builders, ...)
	// TODO: envelopes?
	for _, assignment := range option.Assignments {
		nullPathChunksGuards := generator.pathNotNullGuards(converter, assignment.Path)
		for _, guard := range nullPathChunksGuards {
			guards.Set(guard.String(), guard)
		}

		if assignment.Value.Constant != nil {
			guard := MappingGuard{
				Path:  converter.inputRootPath().Append(assignment.Path),
				Op:    ast.EqualOp,
				Value: assignment.Value.Constant,
			}
			guards.Set(guard.String(), guard)
			continue
		}

		assignmentType := assignment.Path.Last().Type

		// For arrays: ensure they're not empty
		if assignmentType.IsArray() {
			guard := MappingGuard{
				Path:  converter.inputRootPath().Append(assignment.Path),
				Op:    ast.MinLengthOp,
				Value: 1,
			}
			guards.Set(guard.String(), guard)
		}

		// For strings: ensure they're not empty
		// TODO: deal with datetime strings
		if assignmentType.IsScalar() && assignmentType.AsScalar().ScalarKind == ast.KindString && !assignmentType.HasHint(ast.HintStringFormatDateTime) {
			guard := MappingGuard{
				Path:  converter.inputRootPath().Append(assignment.Path),
				Op:    ast.NotEqualOp,
				Value: "",
			}
			guards.Set(guard.String(), guard)
		}

		// TODO: is that correct/needed?
		if assignment.Method != ast.AppendAssignment && assignment.Value.Envelope != nil {
			for _, envelopePath := range assignment.Value.Envelope.Values {
				guard := MappingGuard{
					Path:  converter.inputRootPath().Append(assignment.Path.Append(envelopePath.Path)),
					Op:    ast.NotEqualOp,
					Value: nil,
				}
				guards.Set(guard.String(), guard)
			}
			continue
		}
	}

	return guards.Values()
}

func (generator *ConverterGenerator) pathNotNullGuards(converter Converter, path ast.Path) []MappingGuard {
	var guards []MappingGuard

	for i, chunk := range path {
		chunkType := chunk.Type
		if chunk.TypeHint != nil {
			chunkType = *chunk.TypeHint
		}

		// TODO: this is language-specific
		maybeNull := chunkType.Nullable || chunkType.IsAnyOf(ast.KindMap, ast.KindArray)
		if !maybeNull {
			continue
		}

		guards = append(guards, MappingGuard{
			Path:  converter.inputRootPath().Append(path[:i+1]),
			Op:    ast.NotEqualOp,
			Value: nil,
		})
	}

	return guards
}
