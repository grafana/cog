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

type BuilderChoiceMapping struct {
	Guards  []MappingGuard
	Builder BuilderArgMapping
}

type ArrayArgMapping struct {
	For       ast.Path
	ForType   ast.Type
	ForArg    *ArgumentMapping
	ValueAs   ast.Path
	ValueType ast.Type
}

type MapArgMapping struct {
	For       ast.Path
	ForType   ast.Type
	ForArg    *ArgumentMapping
	ValueAs   ast.Path
	IndexType ast.Type
	ValueType ast.Type
}

type RuntimeArgMapping struct {
	FuncName string
	Args     []*DirectArgMapping
}

type DisjunctionBranchArgMapping struct {
	Type ast.Type
	Of   *DirectArgMapping
	Arg  *ArgumentMapping
}

type DisjunctionArgMapping struct {
	Branches []DisjunctionBranchArgMapping
}

type ArgumentMapping struct {
	Direct             *DirectArgMapping
	Runtime            *RuntimeArgMapping
	Builder            *BuilderArgMapping
	BuilderDisjunction []BuilderChoiceMapping
	Array              *ArrayArgMapping
	Map                *MapArgMapping
	Disjunction        *DisjunctionArgMapping

	// TODO: used? necessary?
	Guards []MappingGuard
}

type ConversionMapping struct {
	// for _, panel := range input.Panels { WithPanel(panel) }
	RepeatFor ast.Path // `input.Panels`
	RepeatAs  string   // `panel`

	Options []OptionMapping
}

type OptionMapping struct {
	Option ast.Option // option in the builder

	Guards []MappingGuard
	Args   []ArgumentMapping
}

func (optMapping OptionMapping) ArgumentGuards() []MappingGuard {
	var guards []MappingGuard

	for _, arg := range optMapping.Args {
		guards = append(guards, arg.Guards...)
	}

	return guards
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

	Mappings []ConversionMapping
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
	nullableTypes NullableConfig

	// generatedPaths lets us keep track of the paths in the input that we generated option mappings for.
	// Since several options can represent with a single path, it allows us to not have "duplicates".
	generatedPaths map[string]struct{}

	listOfDisjunctionOptions map[string][]ast.Option
}

func NewConverterGenerator(nullableTypes NullableConfig) *ConverterGenerator {
	return &ConverterGenerator{
		nullableTypes:            nullableTypes,
		generatedPaths:           make(map[string]struct{}),
		listOfDisjunctionOptions: make(map[string][]ast.Option),
	}
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

	converter.Mappings = tools.Map(builder.Options, func(option ast.Option) ConversionMapping {
		return generator.convertOption(context, converter, option)
	})

	for _, opts := range generator.listOfDisjunctionOptions {
		converter.Mappings = append(converter.Mappings, generator.convertListOfDisjunctionOptions(context, converter, opts))
	}

	converter.Mappings = tools.Filter(converter.Mappings, func(mapping ConversionMapping) bool {
		return len(mapping.Options) != 0
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

func (generator *ConverterGenerator) convertListOfDisjunctionOptions(context Context, converter Converter, options []ast.Option) ConversionMapping {
	mapping := ConversionMapping{
		RepeatFor: converter.inputRootPath().Append(options[0].Assignments[0].Path),
		RepeatAs:  "item",
	}

	mapping.Options = tools.Map(options, func(option ast.Option) OptionMapping {
		return generator.mappingForOption(context, converter, mapping, option)
	})
	mapping.Options = tools.Filter(mapping.Options, func(optMapping OptionMapping) bool {
		return optMapping.Option.Name != ""
	})

	return mapping
}

func (generator *ConverterGenerator) convertOption(context Context, converter Converter, option ast.Option) ConversionMapping {
	assignments := tools.Filter(option.Assignments, func(assignment ast.Assignment) bool {
		_, pathAlreadyGenerated := generator.generatedPaths[generator.assignmentKey(assignment)]
		return !pathAlreadyGenerated
	})
	if len(assignments) == 0 {
		return ConversionMapping{}
	}

	mapping := ConversionMapping{}

	if len(assignments) == 1 && assignments[0].Method == ast.AppendAssignment {
		mapping.RepeatFor = converter.inputRootPath().Append(assignments[0].Path)
		mapping.RepeatAs = "item"
	}

	// if the option appends one possible branch of a disjunction to a list,
	// we need to treat it differently
	if mapping.RepeatFor != nil && generator.isAssignmentFromDisjunctionStruct(context, assignments[0]) {
		path := assignments[0].Path.String()
		generator.listOfDisjunctionOptions[path] = append(generator.listOfDisjunctionOptions[path], option)
		return ConversionMapping{}
	}

	optMapping := generator.mappingForOption(context, converter, mapping, option)
	if optMapping.Option.Name == "" {
		return ConversionMapping{}
	}

	mapping.Options = []OptionMapping{optMapping}

	return mapping
}

func (generator *ConverterGenerator) mappingForOption(context Context, converter Converter, mapping ConversionMapping, option ast.Option) OptionMapping {
	optMapping := OptionMapping{
		Option: option,
		Guards: generator.guardForAssignments(converter.inputRootPath(), option.Assignments),
	}

	assignments := tools.Filter(option.Assignments, func(assignment ast.Assignment) bool {
		_, pathAlreadyGenerated := generator.generatedPaths[generator.assignmentKey(assignment)]
		return !pathAlreadyGenerated
	})
	if len(assignments) == 0 {
		return OptionMapping{}
	}

	i := 0
	for _, assignment := range assignments {
		i++

		generator.generatedPaths[generator.assignmentKey(assignment)] = struct{}{}

		// no need for an argument if the assignment uses a constant value
		if assignment.Value.Constant != nil {
			continue
		}

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
			optMapping.Args = append(optMapping.Args, argument)
			continue
		}

		if assignment.Value.Envelope != nil {
			arguments := generator.argumentsForEnvelope(context, converter, argName, valuePath, assignment)
			optMapping.Args = append(optMapping.Args, arguments...)
			continue
		}

		argument := generator.argumentForType(context, converter, argName, valuePath, valueType)
		optMapping.Args = append(optMapping.Args, argument)
	}

	return optMapping
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
		return ArgumentMapping{
			Runtime: &RuntimeArgMapping{
				FuncName: fmt.Sprintf("Convert%sToCode", tools.UpperCamelCase(string(typeDef.AsComposableSlot().Variant))),
				Args: []*DirectArgMapping{
					{ValuePath: valuePath, ValueType: typeDef},
				},
			},
		}
	}

	if typeDef.IsDisjunction() {
		return ArgumentMapping{
			Disjunction: &DisjunctionArgMapping{
				Branches: tools.Map(typeDef.Disjunction.Branches, func(branch ast.Type) DisjunctionBranchArgMapping {
					arg := generator.argumentForType(context, converter, argName, valuePath, branch)
					return DisjunctionBranchArgMapping{
						Type: branch,
						Of:   &DirectArgMapping{ValuePath: valuePath, ValueType: typeDef},
						Arg:  &arg,
					}
				}),
			},
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

	if typeDef.IsMap() {
		valueAs := ast.Path{
			{Identifier: argName, Type: typeDef.Map.ValueType, Root: true},
		}

		forArg := generator.argumentForType(context, converter, argName+"Value", valueAs, typeDef.Map.ValueType)

		return ArgumentMapping{
			Map: &MapArgMapping{
				For:       valuePath,
				ForType:   typeDef,
				ForArg:    &forArg,
				ValueAs:   valueAs,
				IndexType: typeDef.Map.IndexType,
				ValueType: typeDef.Map.ValueType,
			},
		}
	}

	possibleBuilders := generator.buildersForType(context, typeDef)
	// hack to use the runtime to convert panels
	// TODO: find a better way to handle this case (ie: something more generic than hardcoding it :/)
	if len(possibleBuilders) > 1 && strings.EqualFold(possibleBuilders[0].Package, "dashboard") && strings.EqualFold("panel", possibleBuilders[0].For.Name) {
		typeField, _ := possibleBuilders[0].For.Type.Struct.FieldByName("type")

		return ArgumentMapping{
			Runtime: &RuntimeArgMapping{
				FuncName: "ConvertPanelToCode",
				Args: []*DirectArgMapping{
					{ValuePath: valuePath, ValueType: typeDef},
					{ValuePath: valuePath.AppendStructField(typeField), ValueType: typeField.Type},
				},
			},
		}
	}
	if len(possibleBuilders) > 1 && possibleBuilders.HaveConstantConstructorAssignment() {
		choices := make([]BuilderChoiceMapping, 0, len(possibleBuilders))
		for _, possibleBuilder := range possibleBuilders {
			constantAssignments := tools.Filter(possibleBuilder.Constructor.Assignments, func(assignment ast.Assignment) bool {
				return assignment.HasConstantValue()
			})

			choices = append(choices, BuilderChoiceMapping{
				Guards: generator.guardForAssignments(valuePath, constantAssignments),
				Builder: BuilderArgMapping{
					ValuePath:   valuePath,
					ValueType:   typeDef,
					BuilderPkg:  possibleBuilder.Package,
					BuilderName: possibleBuilder.Name,
				},
			})
		}

		return ArgumentMapping{
			BuilderDisjunction: choices,
		}
	}
	if len(possibleBuilders) > 0 {
		return ArgumentMapping{
			Builder: &BuilderArgMapping{
				ValuePath:   valuePath,
				ValueType:   typeDef,
				BuilderPkg:  possibleBuilders[0].Package,
				BuilderName: possibleBuilders[0].Name,
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

func (generator *ConverterGenerator) isAssignmentFromDisjunctionStruct(context Context, assignment ast.Assignment) bool {
	if assignment.Value.Envelope == nil {
		return false
	}

	envelopedType := assignment.Value.Envelope.Type
	if envelopedType.IsRef() {
		referredObject, _ := context.LocateObject(envelopedType.Ref.ReferredPkg, envelopedType.Ref.ReferredType)
		envelopedType = referredObject.Type
	}

	return envelopedType.IsStructGeneratedFromDisjunction()
}

func (generator *ConverterGenerator) argumentFromDisjunctionStruct(context Context, converter Converter, argName string, valuePath ast.Path, assignment ast.Assignment) (ArgumentMapping, bool) {
	if !generator.isAssignmentFromDisjunctionStruct(context, assignment) {
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

func (generator *ConverterGenerator) guardForAssignments(valuesRootPath ast.Path, assignments []ast.Assignment) []MappingGuard {
	// conditions safeguarding the conversion of the current option
	guards := orderedmap.New[string, MappingGuard]()

	// TODO: define guards other than "not null" checks? (0, "", ...)
	// TODO: builders + array of builders (and array of array of builders, ...)
	// TODO: envelopes?
	for _, assignment := range assignments {
		nullPathChunksGuards := generator.pathNotNullGuards(valuesRootPath, assignment.Path)
		for _, guard := range nullPathChunksGuards {
			guards.Set(guard.String(), guard)
		}

		if assignment.Value.Constant != nil {
			guard := MappingGuard{
				Path:  valuesRootPath.Append(assignment.Path),
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
				Path:  valuesRootPath.Append(assignment.Path),
				Op:    ast.MinLengthOp,
				Value: 1,
			}
			guards.Set(guard.String(), guard)
		}

		// For strings: ensure they're not empty
		// TODO: deal with datetime strings
		if assignmentType.IsScalar() && assignmentType.AsScalar().ScalarKind == ast.KindString && !assignmentType.HasHint(ast.HintStringFormatDateTime) {
			guard := MappingGuard{
				Path:  valuesRootPath.Append(assignment.Path),
				Op:    ast.NotEqualOp,
				Value: "",
			}
			guards.Set(guard.String(), guard)
		}

		// For scalar values, add a guard against assignments equal to the default value for that path
		if assignmentType.IsScalar() && assignmentType.Default != nil {
			guard := MappingGuard{
				Path:  valuesRootPath.Append(assignment.Path),
				Op:    ast.NotEqualOp,
				Value: assignmentType.Default,
			}
			guards.Set(guard.String(), guard)
		}

		// TODO: is that correct/needed?
		if assignment.Method != ast.AppendAssignment && assignment.Value.Envelope != nil {
			for _, envelopePath := range assignment.Value.Envelope.Values {
				guard := MappingGuard{
					Path:  valuesRootPath.Append(assignment.Path.Append(envelopePath.Path)),
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

func (generator *ConverterGenerator) pathNotNullGuards(rootPath ast.Path, path ast.Path) []MappingGuard {
	var guards []MappingGuard

	for i, chunk := range path {
		if !generator.nullableTypes.TypeIsNullable(chunk.Type) {
			continue
		}

		guards = append(guards, MappingGuard{
			Path:  rootPath.Append(path[:i+1]),
			Op:    ast.NotEqualOp,
			Value: nil,
		})
	}

	return guards
}

func (generator *ConverterGenerator) assignmentKey(assignment ast.Assignment) string {
	path := assignment.Path.String()

	if assignment.Value.Constant != nil {
		path += fmt.Sprintf("=%v", assignment.Value.Constant)
	}

	if assignment.Value.Envelope != nil {
		// TODO: envelope of envelope?
		for _, envelopeAssignment := range assignment.Value.Envelope.Values {
			path += "," + envelopeAssignment.Path.String()
		}
	}

	return path
}

func (generator *ConverterGenerator) buildersForType(context Context, typeDef ast.Type) ast.Builders {
	var candidateBuilders ast.Builders

	var search func(def ast.Type)
	search = func(def ast.Type) {
		if def.IsArray() {
			search(def.AsArray().ValueType)
			return
		}
		if def.IsMap() {
			search(def.AsMap().ValueType)
			return
		}

		if def.IsDisjunction() {
			for _, branch := range def.AsDisjunction().Branches {
				search(branch)
			}

			return
		}

		if !def.IsRef() {
			return
		}

		ref := def.AsRef()
		builders := context.Builders.LocateAllByObject(ref.ReferredPkg, ref.ReferredType)
		candidateBuilders = append(candidateBuilders, builders...)
	}

	search(typeDef)

	return candidateBuilders
}
