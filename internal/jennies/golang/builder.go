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
		builderImportAlias := jenny.typeImportAlias(builder.For.SelfRef)
		jenny.typeImportMapper = func(pkg string) string {
			if pkg == builder.For.SelfRef.ReferredPkg {
				return builderImportAlias
			}

			jenny.imports.Add(pkg, "github.com/grafana/cog/generated/types/"+pkg)

			return pkg
		}

		output := jenny.generateBuilder(context, builder)
		filename := filepath.Join(
			strings.ToLower(builder.RootPackage),
			strings.ToLower(builder.Package),
			"builder_gen.go",
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

	objectName := tools.UpperCamelCase(builder.For.Name)

	// import generated types
	importAlias := jenny.importType(builder.For.SelfRef)

	// Option type declaration
	buffer.WriteString("type Option func(builder *Builder) error\n\n")

	// Builder type declaration
	buffer.WriteString(fmt.Sprintf(`type Builder struct {
	internal *%s.%s
}
`, importAlias, objectName))

	// Add a constructor for the builder
	constructorCode := jenny.generateConstructor(context, builder)
	buffer.WriteString(constructorCode)

	// Allow builders to expose the resource they're building
	buffer.WriteString(fmt.Sprintf(`
func (builder *Builder) Build() *%s.%s {
	return builder.internal
}
`, importAlias, objectName))

	// Define options
	for _, option := range builder.Options {
		buffer.WriteString(jenny.generateOption(context, option) + "\n")
	}

	// add calls to set default values
	buffer.WriteString("\n")
	buffer.WriteString("func defaults() []Option {\n")
	buffer.WriteString("return []Option{\n")
	for _, opt := range builder.Options {
		if opt.Default != nil {
			buffer.WriteString(jenny.generateDefaultCall(opt) + ",\n")
		}
	}
	buffer.WriteString("}\n")
	buffer.WriteString("}\n")

	return buffer.String()
}

func (jenny *Builder) generateConstructor(context context.Builders, builder ast.Builder) string {
	var buffer strings.Builder

	typeName := tools.UpperCamelCase(builder.For.Name)
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
		args = strings.Join(argsList, ", ") + ", "
	}
	if len(fieldsInitList) != 0 {
		fieldsInit = strings.Join(fieldsInitList, "\n") + "\n"
	}

	buffer.WriteString(fmt.Sprintf(`func New(%[2]soptions ...Option) (Builder, error) {
	resource := &%[4]s.%[1]s{}
	builder := &Builder{internal: resource}

	%[3]s
	for _, opt := range append(defaults(), options...) {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
}
`, typeName, args, fieldsInit, jenny.typeImportAlias(builder.For.SelfRef)))

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

func (jenny *Builder) generateOption(context context.Builders, def ast.Option) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	// Option name
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

	buffer.WriteString(fmt.Sprintf(`func %[1]s(%[2]s) Option {
	return func(builder *Builder) error {
		%[3]s

		return nil
	}
}
`, optionName, arguments, assignments))

	return buffer.String()
}

func (jenny *Builder) generateArgument(context context.Builders, arg ast.Argument) string {
	if referredBuilder, found := context.BuilderForType(arg.Type); found {
		importAlias := jenny.importBuilder(referredBuilder)

		return fmt.Sprintf(`opts ...%s.Option`, importAlias)
	}

	typeName := strings.Trim(formatType(arg.Type, jenny.typeImportMapper), "*")
	name := jenny.escapeVarName(tools.LowerCamelCase(arg.Name))

	return fmt.Sprintf("%s %s", name, typeName)
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

		nullable := chunk.Type.Nullable || chunk.Type.Kind == ast.KindMap || chunk.Type.Kind == ast.KindArray || chunk.Type.IsAny()
		if nullable {
			subPath := assignment.Path[:i+1]
			pathInitSafeGuards = append(pathInitSafeGuards, jenny.generatePathInitializationSafeGuard(subPath))
		}
	}

	assignmentSafeGuards := strings.Join(pathInitSafeGuards, "\n")
	if assignmentSafeGuards != "" {
		assignmentSafeGuards += "\n\n"
	}

	if referredBuilder, found := context.BuilderForType(valueType); found {
		referredBuilderAlias := jenny.importBuilder(referredBuilder)
		intoPointer := "*"
		if valueType.Nullable {
			intoPointer = ""
		}

		return fmt.Sprintf(`resource, err := %[2]s.New(opts...)
		if err != nil {
			return err
		}

		%[4]sbuilder.internal.%[1]s = %[3]sresource.Build()
`, fieldPath, referredBuilderAlias, intoPointer, assignmentSafeGuards)
	}

	// constant value, not into a pointer type
	if assignment.ArgumentName == "" && !valueType.Nullable {
		return fmt.Sprintf("%[3]sbuilder.internal.%[1]s = %[2]s", fieldPath, formatScalar(assignment.Value), assignmentSafeGuards)
	}
	// constant value, into a pointer type
	if assignment.ArgumentName == "" && valueType.Nullable {
		tmpVarName := "val" + tools.UpperCamelCase(assignment.Path.Last().Identifier)

		return fmt.Sprintf(`%[3]s%[4]s := %[2]s
builder.internal.%[1]s = &%[4]s`, fieldPath, formatScalar(assignment.Value), assignmentSafeGuards, tmpVarName)
	}

	argName := jenny.escapeVarName(tools.LowerCamelCase(assignment.ArgumentName))

	asPointer := ""
	// FIXME: this condition is probably wrong
	if valueType.Kind != ast.KindArray && valueType.Kind != ast.KindMap && valueType.Nullable {
		asPointer = "&"
	}

	generatedConstraints := strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
	if generatedConstraints != "" {
		generatedConstraints += "\n\n"
	}

	return generatedConstraints + fmt.Sprintf("%[4]sbuilder.internal.%[1]s = %[3]s%[2]s", fieldPath, argName, asPointer, assignmentSafeGuards)
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
		return fmt.Sprintf(
			"make(%s)",
			formatType(typeDef, jenny.typeImportMapper),
		)
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

	return fmt.Sprintf("%s(%s)", tools.UpperCamelCase(option.Name), strings.Join(args, ", "))
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

	buffer.WriteString(fmt.Sprintf("if !(%s) {\n", jenny.constraintComparison(argumentName, constraint)))
	buffer.WriteString(fmt.Sprintf("return errors.New(\"%[1]s must be %[2]s %[3]v\")\n", argumentName, constraint.Op, constraint.Args[0]))
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

// typeImportAlias returns the alias to use when importing the given object's type definition.
func (jenny *Builder) typeImportAlias(ref ast.RefType) string {
	// all types within a schema are generated under the same package
	return ref.ReferredPkg
}

// importType declares an import statement for the type definition of
// the given object and returns an alias to it.
func (jenny *Builder) importType(typeRef ast.RefType) string {
	pkg := jenny.typeImportAlias(typeRef)

	jenny.imports.Add(pkg, "github.com/grafana/cog/generated/types/"+pkg)

	return pkg
}

// importBuilderForObject declares an import statement for the builder definition of
// the given object and returns an alias to it.
func (jenny *Builder) importBuilder(builder ast.Builder) string {
	pkg := strings.ToLower(builder.For.Name) // conflict with type import?

	jenny.imports.Add(pkg, fmt.Sprintf("github.com/grafana/cog/generated/%s/%s", strings.ToLower(builder.RootPackage), strings.ToLower(builder.Package)))

	return pkg
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
