package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	imports          template.ImportMap
	typeImportMapper func(pkg string) string
}

func (jenny *Builder) JennyName() string {
	return "GoBuilder"
}

func (jenny *Builder) Generate(context context.Builders) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range context.Builders {
		jenny.typeImportMapper = func(pkg string) string {
			if pkg == builder.Package {
				return ""
			}

			jenny.imports.Add(pkg, "github.com/grafana/cog/generated/"+pkg)

			return pkg
		}

		output := jenny.generateBuilder(context, builder)
		filename := filepath.Join(
			strings.ToLower(builder.Package),
			fmt.Sprintf("%s_builder_gen.go", strings.ToLower(builder.For.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(context context.Builders, builder ast.Builder) []byte {
	var buffer strings.Builder

	jenny.imports = template.NewImportMap()
	jenny.imports.Add("cog", "github.com/grafana/cog/generated")

	fullObjectName := builder.For.Name
	if builder.For.SelfRef.ReferredPkg != builder.Package {
		fullObjectName = builder.For.SelfRef.ReferredPkg + "." + fullObjectName
	}

	err := templates.ExecuteTemplate(&buffer, "builder.tmpl", template.Tmpl{
		Package:     builder.Package,
		Imports:     jenny.imports,
		BuilderName: tools.UpperCamelCase(builder.Name),
		ObjectName:  fullObjectName,
		Constructor: jenny.generateConstructor(context, builder),
		// TODO: Add arguments and assignments to constructor
		Options:        jenny.generateOptions(context, builder),
		DefaultBuilder: jenny.genDefaultBuilder(builder),
	})

	if err != nil {
		return nil
	}

	return []byte(buffer.String())
}

func (jenny *Builder) genDefaultBuilder(builder ast.Builder) template.DefaultBuilder {
	arguments := make([]template.Argument, 0)
	for _, opt := range builder.Options {
		if opt.Default != nil {
			arguments = append(arguments, jenny.generateDefaultCall(opt)...)
		}
	}

	return template.DefaultBuilder{
		Name: tools.UpperCamelCase(builder.Name),
		Args: arguments,
	}
}

func (jenny *Builder) generateConstructor(context context.Builders, builder ast.Builder) template.Constructor {
	var argsList []template.Argument
	var assignmentList []template.Assignment
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		// FIXME: this is assuming that there's only one argument for that option
		argsList = append(argsList, jenny.generateArgument(context, opt.Args[0]))
		assignmentList = append(assignmentList, jenny.generateAssignment(context, opt.Assignments[0]))
	}

	for _, init := range builder.Initializations {
		assignmentList = append(assignmentList, jenny.generateAssignment(context, init))
	}

	return template.Constructor{
		Args:        argsList,
		Assignments: assignmentList,
	}
}

func (jenny *Builder) formatFieldPath(fieldPath ast.Path) string {
	parts := make([]string, len(fieldPath))

	for i := range fieldPath {
		output := tools.UpperCamelCase(fieldPath[i].Identifier)

		// don't generate type hints if:
		// * there isn't one defined
		// * the type isn't "any"
		// * as a trailing element in the path
		if !fieldPath[i].Type.IsAny() || fieldPath[i].TypeHint == nil || i == len(fieldPath)-1 {
			parts[i] = output
			continue
		}

		formattedTypeHint := formatType(*fieldPath[i].TypeHint, jenny.typeImportMapper)

		parts[i] = output + fmt.Sprintf(".(*%s)", formattedTypeHint)
	}

	return strings.Join(parts, ".")
}

func (jenny *Builder) generateOptions(context context.Builders, builder ast.Builder) []template.Option {
	options := make([]template.Option, len(builder.Options))
	for i, opt := range builder.Options {
		options[i] = jenny.generateOption(context, opt)
	}

	return options
}

func (jenny *Builder) generateOption(context context.Builders, def ast.Option) template.Option {
	// Arguments list
	argsList := make([]template.Argument, 0, len(def.Args))
	if len(def.Args) != 0 {
		for _, arg := range def.Args {
			argsList = append(argsList, jenny.generateArgument(context, arg))
		}
	}

	// Assignments
	assignmentsList := make([]template.Assignment, 0, len(def.Assignments))
	for _, assignment := range def.Assignments {
		assignmentsList = append(assignmentsList, jenny.generateAssignment(context, assignment))
	}

	return template.Option{
		Name:        tools.UpperCamelCase(def.Name),
		Comments:    def.Comments,
		Args:        argsList,
		Assignments: assignmentsList,
	}
}

func (jenny *Builder) generateArgument(context context.Builders, arg ast.Argument) template.Argument {
	argName := jenny.escapeVarName(tools.LowerCamelCase(arg.Name))

	if referredBuilder, found := context.BuilderForType(arg.Type); found {
		qualifiedType := jenny.importType(referredBuilder.For.SelfRef)

		return template.Argument{
			Name:          argName,
			Type:          qualifiedType,
			ReferredAlias: "cog",
		}
	}

	return template.Argument{
		Name: argName,
		Type: strings.Trim(formatType(arg.Type, jenny.typeImportMapper), "*"),
	}
}

func (jenny *Builder) generatePathInitializationSafeGuard(path ast.Path) string {
	fieldPath := jenny.formatFieldPath(path)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := jenny.emptyValueForType(valueType)
	// This should be alright since there shouldn't be any scalar in the middle of a path
	if emptyValue[0] == '*' {
		emptyValue = "&" + emptyValue[1:]
	}

	if path.Last().Type.IsAny() && emptyValue[0] != '&' {
		emptyValue = "&" + emptyValue
	}

	return fmt.Sprintf(`if builder.internal.%[1]s == nil {
	builder.internal.%[1]s = %[2]s
}`, fieldPath, emptyValue)
}

func (jenny *Builder) generateAssignment(context context.Builders, assignment ast.Assignment) template.Assignment {
	fieldPath := jenny.formatFieldPath(assignment.Path)
	pathEnd := assignment.Path.Last()
	valueType := pathEnd.Type

	var initSafeGuards []string
	for i, chunk := range assignment.Path {
		if i == len(assignment.Path)-1 {
			continue
		}

		nullable := chunk.Type.Nullable ||
			chunk.Type.Kind == ast.KindMap ||
			chunk.Type.Kind == ast.KindArray ||
			chunk.Type.IsAny()
		if nullable {
			subPath := assignment.Path[:i+1]
			initSafeGuards = append(initSafeGuards, jenny.generatePathInitializationSafeGuard(subPath))
		}
	}

	constraints := make([]template.Constraint, 0)
	if assignment.Value.Argument != nil {
		argName := jenny.escapeVarName(tools.LowerCamelCase(assignment.Value.Argument.Name))
		constraints = append(jenny.constraints(argName, assignment.Constraints))
	}

	assigmentValueType, value := jenny.formatAssignmentValue(context, assignment.Value)
	isBuilder := false
	if assigmentValueType == template.ValueTypeAssigment {
		_, isBuilder = context.BuilderForType(assignment.Value.Argument.Type)
	}

	return template.Assignment{
		Path:           fieldPath,
		InitSafeguards: initSafeGuards,
		Constraints:    constraints,
		Method:         assignment.Method,
		Value:          value,
		ValueType:      assigmentValueType,
		IsBuilder:      isBuilder,
		IntoNullable:   valueType.Nullable && valueType.Kind != ast.KindArray,
	}
}

func (jenny *Builder) formatAssignmentValue(context context.Builders, value ast.AssignmentValue) (template.ValueType, string) {
	// constant value, not into a pointer type
	if value.Constant != nil {
		return template.ValueTypeConstant, formatScalar(value.Constant)
	}

	// envelope
	if value.Envelope != nil {
		return template.ValueTypeEnvelope, jenny.formatEnvelopeAssignmentValue(context, value)
	}

	// argument
	return template.ValueTypeAssigment, jenny.escapeVarName(tools.LowerCamelCase(value.Argument.Name))
}

func (jenny *Builder) formatEnvelopeAssignmentValue(context context.Builders, value ast.AssignmentValue) string {
	envelope := value.Envelope
	formattedType := formatType(envelope.Type, jenny.typeImportMapper)

	var allValues string
	for _, item := range envelope.Values {
		_, val := jenny.formatAssignmentValue(context, item.Value)
		allValues += fmt.Sprintf("%s: %s,\n", tools.UpperCamelCase(item.Path[0].Identifier), val)
	}

	envelopeValue := fmt.Sprintf(`%[1]s{
	%[2]s
}`, formattedType, allValues)

	return envelopeValue
}

func (jenny *Builder) emptyValueForType(typeDef ast.Type) string {
	switch typeDef.Kind {
	case ast.KindRef:
		return formatType(typeDef, jenny.typeImportMapper) + "{}"
	case ast.KindStruct:
		return formatType(typeDef, jenny.typeImportMapper) + "{}"
	case ast.KindEnum:
		return formatScalar(typeDef.AsEnum().Values[0].Value)
	case ast.KindArray, ast.KindMap:
		return formatType(typeDef, jenny.typeImportMapper) + "{}"
	case ast.KindScalar:
		return "" // no need to do anything here

	default:
		return "unknown"
	}
}

func (jenny *Builder) escapeVarName(varName string) string {
	if isReservedGoKeyword(varName) {
		return varName + "Arg"
	}

	return varName
}

func (jenny *Builder) generateDefaultCall(option ast.Option) []template.Argument {
	args := make([]template.Argument, 0, len(option.Default.ArgsValues))
	for _, arg := range option.Default.ArgsValues {
		args = append(args, template.Argument{
			Name: tools.UpperCamelCase(option.Name),
			Type: formatScalar(arg),
		})
	}

	return args
}

func (jenny *Builder) constraints(argumentName string, constraints []ast.TypeConstraint) []template.Constraint {
	output := make([]template.Constraint, 0, len(constraints))

	for _, constraint := range constraints {
		op, isString := jenny.constraintComparison(constraint)
		output = append(output, template.Constraint{
			Name:     argumentName,
			Op:       op,
			Arg:      constraint.Args[0],
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

// importType declares an import statement for the type definition of
// the given object and returns a fully qualified type name for it.
func (jenny *Builder) importType(typeRef ast.RefType) string {
	pkg := jenny.typeImportMapper(typeRef.ReferredPkg)
	typeName := tools.UpperCamelCase(typeRef.ReferredType)
	if pkg == "" {
		return typeName
	}

	return fmt.Sprintf("%s.%s", pkg, typeName)
}

func isReservedGoKeyword(input string) bool {
	// see: https://go.dev/ref/spec#Keywords
	return input == "break" ||
		input == "case" ||
		input == "chan" ||
		input == "continue" ||
		input == "const" ||
		input == "default" ||
		input == "defer" ||
		input == "else" ||
		input == "error" ||
		input == "fallthrough" ||
		input == "for" ||
		input == "func" ||
		input == "go" ||
		input == "goto" ||
		input == "if" ||
		input == "import" ||
		input == "interface" ||
		input == "map" ||
		input == "package" ||
		input == "range" ||
		input == "return" ||
		input == "select" ||
		input == "struct" ||
		input == "switch" ||
		input == "type" ||
		input == "var"
}
