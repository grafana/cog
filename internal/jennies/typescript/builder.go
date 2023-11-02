package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	imports importMap
}

type Tmpl struct {
	BuilderName string
	ObjectName  string
	Imports     importMap
	ImportAlias string
	Options     []options
	Constructor constructor
}

type constructor struct {
	Args        []argument
	Assignments []assignment
}

type options struct {
	Name        string
	Comments    []string
	Args        []argument
	Assignments []assignment
}

type argument struct {
	Name          string
	Type          string
	ReferredAlias string
	ReferredName  string
}

type assignment struct {
	Path           string
	InitSafeguards []string
	Value          string
	IsBuilder      bool
	Method         ast.AssignmentMethod
	Constraints    []constraint
}

type constraint struct {
	Name     string
	Op       ast.Op
	Arg      any
	IsString bool
}

func (jenny *Builder) JennyName() string {
	return "TypescriptBuilder"
}

func (jenny *Builder) Generate(context context.Builders) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range context.Builders {
		output, err := jenny.generateBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			strings.ToLower(builder.Package),
			fmt.Sprintf("%s_builder_gen.ts", strings.ToLower(builder.For.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(context context.Builders, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	jenny.imports = newImportMap()
	importAlias := jenny.importType(builder.For.SelfRef)

	constructorCode := jenny.generateConstructor(context, builder)

	// Define options
	opts := make([]options, len(builder.Options))
	for i, o := range builder.Options {
		opts[i] = jenny.generateOption(context, builder, o)
	}

	err := templates.Lookup("builder.tmpl").Execute(&buffer, Tmpl{
		BuilderName: builder.Name,
		ObjectName:  builder.For.Name,
		Imports:     jenny.imports,
		ImportAlias: importAlias,
		Options:     opts,
		Constructor: constructorCode,
	})

	return []byte(buffer.String()), err
}

func (jenny *Builder) generateConstructor(context context.Builders, builder ast.Builder) constructor {
	var argsList []argument
	var assignments []assignment
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		// FIXME: this is assuming that there's only one argument for that option
		argsList = append(argsList, jenny.generateArgument(context, builder, opt.Args[0]))
		assignments = append(assignments, jenny.generateAssignment(context, builder, opt.Assignments[0]))
	}

	for _, init := range builder.Initializations {
		assignments = append(assignments, jenny.generateAssignment(context, builder, init))
	}

	return constructor{
		Args:        argsList,
		Assignments: assignments,
	}
}

func (jenny *Builder) generateOption(context context.Builders, builder ast.Builder, def ast.Option) options {
	// Arguments list
	argsList := make([]argument, 0, len(def.Args))
	if len(def.Args) != 0 {
		for _, arg := range def.Args {
			argsList = append(argsList, jenny.generateArgument(context, builder, arg))
		}
	}

	// Assignments
	assignmentsList := make([]assignment, 0, len(def.Assignments))
	for _, assign := range def.Assignments {
		assignmentsList = append(assignmentsList, jenny.generateAssignment(context, builder, assign))
	}

	return options{
		Name:        def.Name,
		Comments:    def.Comments,
		Args:        argsList,
		Assignments: assignmentsList,
	}
}

func (jenny *Builder) generateArgument(context context.Builders, builder ast.Builder, arg ast.Argument) argument {
	if referredBuilder, found := context.BuilderForType(arg.Type); found {
		referredTypeAlias := jenny.typeImportAlias(referredBuilder.For.SelfRef)

		return argument{
			Name:          arg.Name,
			ReferredAlias: referredTypeAlias,
			ReferredName:  referredBuilder.For.Name,
		}
	}

	builderImportAlias := jenny.typeImportAlias(builder.For.SelfRef)
	typeName := formatType(arg.Type, func(pkg string) string {
		if pkg == builder.For.SelfRef.ReferredPkg {
			return builderImportAlias
		}

		jenny.imports.Add(pkg, fmt.Sprintf("../%s/types_gen", pkg))

		return pkg
	})

	return argument{Name: tools.LowerCamelCase(arg.Name), Type: typeName}
}

func (jenny *Builder) generatePathInitializationSafeGuard(currentBuilder ast.Builder, path ast.Path) string {
	fieldPath := jenny.formatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := formatValue(defaultValueForType(currentBuilder.Schema, valueType, func(pkg string) string {
		jenny.imports.Add(pkg, fmt.Sprintf("../%s/types_gen", pkg))

		return pkg
	}))

	return fmt.Sprintf(`		if (!this.internal.%[1]s) {
			this.internal.%[1]s = %[2]s;
		}`, fieldPath, emptyValue)
}

func (jenny *Builder) generateAssignment(context context.Builders, builder ast.Builder, assign ast.Assignment) assignment {
	fieldPath := jenny.formatFieldPath(assign.Path)

	var pathInitSafeGuards []string
	for i, chunk := range assign.Path {
		if i == len(assign.Path)-1 && assign.Method != ast.AppendAssignment {
			continue
		}

		chunkType := chunk.Type
		if chunk.TypeHint != nil {
			chunkType = *chunk.TypeHint
		}

		maybeUndefined := chunkType.Nullable ||
			chunkType.Kind == ast.KindMap ||
			chunkType.Kind == ast.KindArray ||
			chunkType.Kind == ast.KindRef ||
			chunkType.Kind == ast.KindStruct

		if !maybeUndefined {
			continue
		}

		subPath := assign.Path[:i+1]
		pathInitSafeGuards = append(pathInitSafeGuards, jenny.generatePathInitializationSafeGuard(builder, subPath))
	}

	var constraints []constraint
	isBuilder := false

	if assign.Value.Argument != nil {
		argName := tools.LowerCamelCase(assign.Value.Argument.Name)
		_, isBuilder = context.BuilderForType(assign.Value.Argument.Type)
		constraints = jenny.constraints(argName, assign.Constraints)
	}

	return assignment{
		Path:           fieldPath,
		Method:         assign.Method,
		InitSafeguards: pathInitSafeGuards,
		Value:          jenny.formatAssignmentValue(assign.Value),
		IsBuilder:      isBuilder,
		Constraints:    constraints,
	}
}

func (jenny *Builder) formatAssignmentValue(value ast.AssignmentValue) string {
	// constant value, not into a pointer type
	if value.Constant != nil {
		return formatScalar(value.Constant)
	}

	// envelope
	if value.Envelope != nil {
		return jenny.formatEnvelopeAssignmentValue(value)
	}

	// argument
	return tools.LowerCamelCase(value.Argument.Name)
}

func (jenny *Builder) formatEnvelopeAssignmentValue(value ast.AssignmentValue) string {
	var allValues string
	for _, item := range value.Envelope.Values {
		val := jenny.formatAssignmentValue(item.Value)
		allValues += fmt.Sprintf("%s: %s,\n", item.Path[0].Identifier, val)
	}

	envelopeValue := fmt.Sprintf(`{
	%s
}`, allValues)

	return envelopeValue
}

func (jenny *Builder) constraints(argumentName string, constraints []ast.TypeConstraint) []constraint {
	output := make([]constraint, 0, len(constraints))

	for _, c := range constraints {
		op, isString := jenny.constraintComparison(c)
		output = append(output, constraint{
			Name:     argumentName,
			Op:       op,
			Arg:      c.Args[0],
			IsString: isString,
		})
	}

	return output
}

func (jenny *Builder) constraintComparison(constraint ast.TypeConstraint) (ast.Op, bool) {
	if constraint.Op == ast.MinLengthOp {
		return ast.GreaterThanEqualOp, true
	}
	if constraint.Op == ast.MaxLengthOp {
		return ast.LessThanEqualOp, true
	}

	return constraint.Op, false
}

// typeImportAlias returns the alias to use when importing the given object's type definition.
func (jenny *Builder) typeImportAlias(typeRef ast.RefType) string {
	// all types within a schema are generated under the same package
	return typeRef.ReferredPkg
}

// importType declares an import statement for the type definition of
// the given object and returns an alias to it.
func (jenny *Builder) importType(typeRef ast.RefType) string {
	pkg := jenny.typeImportAlias(typeRef)

	jenny.imports.Add(pkg, fmt.Sprintf("../%s/types_gen", pkg))

	return pkg
}

func (jenny *Builder) formatFieldPath(fieldPath ast.Path) string {
	return strings.Join(tools.Map(fieldPath, func(chunk ast.PathItem) string {
		return chunk.Identifier
	}), ".")
}
