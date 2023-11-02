package golang

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
	imports          importMap
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

	jenny.imports = newImportMap()

	builderSource := jenny.generateBuilderSource(context, builder)

	// package declaration
	buffer.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(builder.Package)))

	// write import statements
	buffer.WriteString(jenny.imports.Format())
	buffer.WriteString("\n\n")

	// write the builder source code
	buffer.WriteString(builderSource)

	return []byte(buffer.String())
}

func (jenny *Builder) generateBuilderSource(context context.Builders, builder ast.Builder) string {
	var buffer strings.Builder

	builderName := tools.UpperCamelCase(builder.Name)

	// import generated types
	cogAlias := jenny.importCog()
	qualifiedObjectName := jenny.importType(builder.For.SelfRef)

	buildObjectSignature := qualifiedObjectName
	if builder.For.Type.ImplementsVariant() {
		buildObjectSignature = variantInterface(builder.For.Type.ImplementedVariant())
	}

	// just to make explicit that this builder implements the generic Cog builder interface
	buffer.WriteString(fmt.Sprintf("var _ %[1]s.Builder[%[2]s] = (*%[3]sBuilder)(nil)\n\n", cogAlias, buildObjectSignature, builderName))

	// Builder type declaration
	buffer.WriteString(fmt.Sprintf(`type %[2]sBuilder struct {
	internal *%[1]s
	errors map[string]%[3]s.BuildErrors
}
`, qualifiedObjectName, builderName, cogAlias))

	// Add a constructor for the builder
	constructorCode := jenny.generateConstructor(context, builder)
	buffer.WriteString(constructorCode)

	// Allow builders to expose the resource they're building
	buffer.WriteString(fmt.Sprintf(`
func (builder *%[2]sBuilder) Build() (%[4]s, error) {
	var errs %[3]s.BuildErrors

	for _, err := range builder.errors {
		errs = append(errs, %[3]s.MakeBuildErrors("%[2]s", err)...)
	}

	if len(errs) != 0 {
		return %[1]s{}, errs
	}

	return *builder.internal, nil
}
`, qualifiedObjectName, builderName, cogAlias, buildObjectSignature))

	// Define options
	for _, option := range builder.Options {
		buffer.WriteString(jenny.generateOption(context, builder, option) + "\n")
	}

	// add calls to set default values
	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("func (builder *%[1]sBuilder) applyDefaults() {\n", builderName))
	for _, opt := range builder.Options {
		if opt.Default != nil {
			buffer.WriteString(jenny.generateDefaultCall(opt) + "\n")
		}
	}
	buffer.WriteString("}\n")

	return buffer.String()
}

func (jenny *Builder) generateConstructor(context context.Builders, builder ast.Builder) string {
	var buffer strings.Builder

	cogAlias := jenny.importCog()
	args := ""
	fieldsInit := ""
	var argsList []string
	var fieldsInitList []string
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		// FIXME: this is assuming that there's only one argument for that option
		argsList = append(argsList, jenny.generateArgument(context, opt.Args[0]))
		fieldsInitList = append(
			fieldsInitList,
			jenny.generateAssignment(context, opt.Assignments[0]),
		)
	}

	for _, init := range builder.Initializations {
		fieldsInitList = append(
			fieldsInitList,
			jenny.generateAssignment(context, init),
		)
	}

	if len(argsList) != 0 {
		args = strings.Join(argsList, ", ")
	}
	if len(fieldsInitList) != 0 {
		fieldsInit = strings.Join(fieldsInitList, "\n") + "\n"
	}

	qualifiedObjectName := jenny.importType(builder.For.SelfRef)
	builderName := tools.UpperCamelCase(builder.Name)

	buffer.WriteString(fmt.Sprintf(`func New%[1]sBuilder(%[2]s) *%[1]sBuilder {
	resource := &%[4]s{}
	builder := &%[1]sBuilder{
		internal: resource,
		errors: make(map[string]%[5]s.BuildErrors),
	}

	%[3]s
	builder.applyDefaults()

	return builder
}
`, builderName, args, fieldsInit, qualifiedObjectName, cogAlias))

	return buffer.String()
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

func (jenny *Builder) generateOption(context context.Builders, builder ast.Builder, def ast.Option) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	builderName := tools.UpperCamelCase(builder.Name)
	optionName := tools.UpperCamelCase(def.Name)

	// Arguments list
	arguments := ""
	if len(def.Args) != 0 {
		argsList := make([]string, 0, len(def.Args))
		for _, arg := range def.Args {
			argsList = append(argsList, jenny.generateArgument(context, arg))
		}

		arguments = strings.Join(argsList, ", ")
	}

	// Assignments
	assignmentsList := make([]string, 0, len(def.Assignments))
	for _, assignment := range def.Assignments {
		assignmentsList = append(assignmentsList, jenny.generateAssignment(context, assignment))
	}
	assignments := strings.Join(assignmentsList, "\n")

	buffer.WriteString(fmt.Sprintf(`func (builder *%[1]sBuilder) %[2]s(%[3]s) *%[1]sBuilder {
	%[4]s

	return builder
}
`, builderName, optionName, arguments, assignments))

	return buffer.String()
}

func (jenny *Builder) generateArgument(context context.Builders, arg ast.Argument) string {
	argName := jenny.escapeVarName(tools.LowerCamelCase(arg.Name))

	if composableSlot, isRefToComposable := context.RefToComposableSlot(arg.Type); isRefToComposable {
		cogAlias := jenny.importCog()
		qualifiedType := formatType(composableSlot, jenny.typeImportMapper)

		return fmt.Sprintf(`%[1]s %[2]s.Builder[%[3]s]`, argName, cogAlias, qualifiedType)
	}

	if referredBuilder, found := context.BuilderForType(arg.Type); found {
		cogAlias := jenny.importCog()
		qualifiedType := jenny.importType(referredBuilder.For.SelfRef)

		return fmt.Sprintf(`%[1]s %[2]s.Builder[%[3]s]`, argName, cogAlias, qualifiedType)
	}

	typeName := strings.Trim(formatType(arg.Type, jenny.typeImportMapper), "*")

	return fmt.Sprintf("%s %s", argName, typeName)
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

func (jenny *Builder) generateAssignment(context context.Builders, assignment ast.Assignment) string {
	fieldPath := jenny.formatFieldPath(assignment.Path)
	pathEnd := assignment.Path.Last()
	valueType := pathEnd.Type

	var pathInitSafeGuards []string
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
			pathInitSafeGuards = append(pathInitSafeGuards, jenny.generatePathInitializationSafeGuard(subPath))
		}
	}

	assignmentSafeGuards := strings.Join(pathInitSafeGuards, "\n")
	if assignmentSafeGuards != "" {
		assignmentSafeGuards += "\n\n"
	}

	constraintsChecks := ""
	if assignment.Value.Argument != nil {
		argName := jenny.escapeVarName(tools.LowerCamelCase(assignment.Value.Argument.Name))
		constraintsChecks = strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
		if constraintsChecks != "" {
			constraintsChecks += "\n\n"
		}
	}

	assignmentSetup, assignmentSource := jenny.formatAssignmentValue(context, assignment.Value, valueType)

	return constraintsChecks + assignmentSafeGuards + assignmentSetup + jenny.formatAssignment(assignment, fieldPath, assignmentSource)
}

func (jenny *Builder) formatAssignmentValue(context context.Builders, value ast.AssignmentValue, valueType ast.Type) (string, string) {
	// constant value, not into a pointer type
	if value.Constant != nil {
		return jenny.formatConstantAssignmentValue(value, valueType)
	}

	// envelope
	if value.Envelope != nil {
		return jenny.formatEnvelopeAssignmentValue(context, value)
	}

	// argument
	return jenny.formatArgumentAssignmentValue(context, value, valueType)
}

func (jenny *Builder) formatArgumentAssignmentValue(context context.Builders, value ast.AssignmentValue, valueType ast.Type) (string, string) {
	argName := jenny.escapeVarName(tools.LowerCamelCase(value.Argument.Name))

	switch valueType.Kind {
	case ast.KindArray:
		valueType = valueType.AsArray().ValueType
	case ast.KindMap:
		valueType = valueType.AsMap().ValueType
	}

	maybeAsPointer := ""
	if valueType.Nullable {
		maybeAsPointer = "&"
	}

	_, hasBuilder := context.BuilderForType(value.Argument.Type)
	_, isRefToComposable := context.RefToComposableSlot(value.Argument.Type)
	if hasBuilder || isRefToComposable {
		cogAlias := jenny.importCog()

		resourceBuildSource := fmt.Sprintf(`resource, err := %[1]s.Build()
if err != nil {
	builder.errors["%[1]s"] = err.(%[2]s.BuildErrors)
	return builder
}

`, argName, cogAlias)

		return resourceBuildSource, fmt.Sprintf("%[1]sresource", maybeAsPointer)
	}

	return "", maybeAsPointer + argName
}

func (jenny *Builder) formatEnvelopeAssignmentValue(context context.Builders, value ast.AssignmentValue) (string, string) {
	envelope := value.Envelope
	formattedType := formatType(envelope.Type, jenny.typeImportMapper)

	var allSetup, allValues string

	for _, item := range envelope.Values {
		setup, val := jenny.formatAssignmentValue(context, item.Value, item.Path.Last().Type)
		allSetup += setup
		allValues += fmt.Sprintf("%s: %s,\n", tools.UpperCamelCase(item.Path[0].Identifier), val)
	}

	envelopeValue := fmt.Sprintf(`%[1]s{
	%[2]s
}`, formattedType, allValues)

	return allSetup, envelopeValue
}

func (jenny *Builder) formatConstantAssignmentValue(value ast.AssignmentValue, valueType ast.Type) (string, string) {
	// constant value, not into a pointer type
	if !valueType.Nullable {
		return "", formatScalar(value.Constant)
	}

	// constant value, into a pointer type
	// todo: find a better name for that temporary variable
	tmpVarName := "val" + tools.UpperCamelCase(ast.TypeName(valueType))
	tmpVarSource := fmt.Sprintf("%[1]s := %[2]s\n", tmpVarName, formatScalar(value.Constant))

	return tmpVarSource, "&" + tmpVarName
}

func (jenny *Builder) formatAssignment(assignment ast.Assignment, fieldPath string, value string) string {
	if assignment.Method == ast.DirectAssignment {
		return fmt.Sprintf("builder.internal.%[1]s = %[2]s", fieldPath, value)
	}

	return fmt.Sprintf("builder.internal.%[1]s = append(builder.internal.%[1]s, %[2]s)", fieldPath, value)
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

func (jenny *Builder) generateDefaultCall(option ast.Option) string {
	args := make([]string, 0, len(option.Default.ArgsValues))
	for _, arg := range option.Default.ArgsValues {
		args = append(args, formatScalar(arg))
	}

	return fmt.Sprintf("builder.%s(%s)", tools.UpperCamelCase(option.Name), strings.Join(args, ", "))
}

func (jenny *Builder) constraints(argumentName string, constraints []ast.TypeConstraint) []string {
	output := make([]string, 0, len(constraints))

	for _, constraint := range constraints {
		output = append(output, jenny.constraint(argumentName, constraint))
	}

	return output
}

func (jenny *Builder) constraint(argumentName string, constraint ast.TypeConstraint) string {
	var buffer strings.Builder

	cogAlias := jenny.importCog()

	buffer.WriteString(fmt.Sprintf("if !(%s) {\n", jenny.constraintComparison(argumentName, constraint)))
	buffer.WriteString(fmt.Sprintf(`builder.errors["%[1]s"] = %[1]s.MakeBuildErrors("%[2]s", errors.New("value must be %[3]s %[4]v"))
`, cogAlias, argumentName, constraint.Op, constraint.Args[0]))
	buffer.WriteString("return builder\n")
	buffer.WriteString("}\n")

	return buffer.String()
}

func (jenny *Builder) constraintComparison(argumentName string, constraint ast.TypeConstraint) string {
	if constraint.Op == ast.MinLengthOp {
		return fmt.Sprintf("len([]rune(%[1]s)) >= %[2]v", argumentName, constraint.Args[0])
	}
	if constraint.Op == ast.MaxLengthOp {
		return fmt.Sprintf("len([]rune(%[1]s)) <= %[2]v", argumentName, constraint.Args[0])
	}

	return fmt.Sprintf("%[1]s %[2]s %#[3]v", argumentName, constraint.Op, constraint.Args[0])
}

func (jenny *Builder) importCog() string {
	jenny.imports.Add("cog", "github.com/grafana/cog/generated")

	return "cog"
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
