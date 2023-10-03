package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	imports importMap
}

func (jenny *Builder) JennyName() string {
	return "GoBuilder"
}

func (jenny *Builder) Generate(builders []ast.Builder) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range builders {
		output := jenny.generateBuilder(builders, builder)
		filename := filepath.Join(
			strings.ToLower(builder.RootPackage),
			strings.ToLower(builder.Package),
			"builder_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(builders ast.Builders, builder ast.Builder) []byte {
	var buffer strings.Builder

	jenny.imports = newImportMap()

	builderSource := jenny.generateBuilderSource(builders, builder)

	// package declaration
	buffer.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(builder.Package)))

	// write import statements
	buffer.WriteString(jenny.imports.Format())
	buffer.WriteString("\n\n")

	// write the builder source code
	buffer.WriteString(builderSource)

	return []byte(buffer.String())
}

func (jenny *Builder) generateBuilderSource(builders ast.Builders, builder ast.Builder) string {
	var buffer strings.Builder

	objectName := tools.UpperCamelCase(builder.For.Name)

	// import generated types
	importAlias := jenny.importType(builder.For)

	// Option type declaration
	buffer.WriteString("type Option func(builder *Builder) error\n\n")

	// Builder type declaration
	buffer.WriteString(fmt.Sprintf(`type Builder struct {
	internal *%s.%s
}
`, importAlias, objectName))

	// Add a constructor for the builder
	constructorCode := jenny.generateConstructor(builders, builder)
	buffer.WriteString(constructorCode)

	// Allow builders to expose the resource they're building
	// TODO: better name, with less conflict chance
	buffer.WriteString(fmt.Sprintf(`
func (builder *Builder) Internal() *%s.%s {
	return builder.internal
}
`, importAlias, objectName))

	// Define options
	for _, option := range builder.Options {
		buffer.WriteString(jenny.generateOption(builders, builder, option) + "\n")
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

func (jenny *Builder) generateConstructor(builders ast.Builders, builder ast.Builder) string {
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
		argsList = append(argsList, jenny.generateArgument(builders, builder, opt.Args[0]))
		fieldsInitList = append(
			fieldsInitList,
			jenny.generateInitAssignment(builders, builder, opt.Assignments[0]),
		)
	}

	for _, init := range builder.Initializations {
		fieldsInitList = append(
			fieldsInitList,
			jenny.generateInitAssignment(builders, builder, init),
		)
	}

	if len(argsList) != 0 {
		args = strings.Join(argsList, ", ") + ", "
	}
	if len(fieldsInitList) != 0 {
		fieldsInit = strings.Join(fieldsInitList, ",\n") + ",\n"
	}

	buffer.WriteString(fmt.Sprintf(`func New(%[2]soptions ...Option) (Builder, error) {
	resource := &%[4]s.%[1]s{
		%[3]s
	}
	builder := &Builder{internal: resource}

	for _, opt := range append(defaults(), options...) {
		if err := opt(builder); err != nil {
			return *builder, err
		}
	}

	return *builder, nil
}
`, typeName, args, fieldsInit, jenny.typeImportAlias(builder.For)))

	return buffer.String()
}

func (jenny *Builder) formatFieldPath(fieldPath string) string {
	parts := strings.Split(fieldPath, ".")
	formatted := make([]string, 0, len(parts))

	for _, part := range parts {
		formatted = append(formatted, tools.UpperCamelCase(part))
	}

	return strings.Join(formatted, ".")
}

func (jenny *Builder) generateInitAssignment(builders ast.Builders, builder ast.Builder, assignment ast.Assignment) string {
	fieldPath := jenny.formatFieldPath(assignment.Path)
	valueType := assignment.ValueType

	if _, valueHasBuilder := jenny.builderForType(builders, builder, assignment.ValueType); valueHasBuilder {
		return "constructor init assignment with type that has a builder is not supported yet"
	}

	if assignment.ArgumentName == "" {
		return fmt.Sprintf("%[1]s: %[2]s", fieldPath, formatScalar(assignment.Value))
	}

	argName := jenny.escapeVarName(tools.LowerCamelCase(assignment.ArgumentName))

	asPointer := ""
	// FIXME: this condition is probably wrong
	if valueType.Kind != ast.KindArray && valueType.Kind != ast.KindStruct && assignment.IntoNullableField {
		asPointer = "&"
	}

	generatedConstraints := strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
	if generatedConstraints != "" {
		generatedConstraints += "\n\n"
	}

	return generatedConstraints + fmt.Sprintf("%[1]s: %[3]s%[2]s", fieldPath, argName, asPointer)
}

func (jenny *Builder) generateOption(builders ast.Builders, builder ast.Builder, def ast.Option) string {
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
			argsList = append(argsList, jenny.generateArgument(builders, builder, arg))
		}

		arguments = strings.Join(argsList, ", ")
	}

	// Assignments
	assignmentsList := make([]string, 0, len(def.Assignments))
	for _, assignment := range def.Assignments {
		assignmentsList = append(assignmentsList, jenny.generateAssignment(builders, builder, assignment))
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

func (jenny *Builder) builderForType(builders ast.Builders, builder ast.Builder, t ast.Type) (ast.Builder, bool) {
	if t.Kind != ast.KindRef {
		return ast.Builder{}, false
	}

	// TODO: shouldn't we using the package from the ref?!
	ref := t.AsRef()
	return builders.LocateByObject(builder.Package, ref.ReferredType)
}

func (jenny *Builder) generateArgument(builders ast.Builders, builder ast.Builder, arg ast.Argument) string {
	if referredBuilder, found := jenny.builderForType(builders, builder, arg.Type); found {
		importAlias := jenny.importBuilder(referredBuilder)

		return fmt.Sprintf(`opts ...%s.Option`, importAlias)
	}

	builderImportAlias := jenny.typeImportAlias(builder.For)
	typeName := strings.Trim(formatType(arg.Type, func(pkg string) string {
		if pkg == builder.For.SelfRef.ReferredPkg {
			return builderImportAlias
		}

		jenny.imports.Add(pkg, fmt.Sprintf("github.com/grafana/cog/generated/types/%s", pkg))

		return pkg
	}), "*")

	name := jenny.escapeVarName(tools.LowerCamelCase(arg.Name))

	return fmt.Sprintf("%s %s", name, typeName)
}

func (jenny *Builder) generateAssignment(builders ast.Builders, builder ast.Builder, assignment ast.Assignment) string {
	fieldPath := jenny.formatFieldPath(assignment.Path)
	valueType := assignment.ValueType

	if referredBuilder, found := jenny.builderForType(builders, builder, assignment.ValueType); found {
		referredBuilderAlias := jenny.importBuilder(referredBuilder)
		intoPointer := "*"
		if assignment.IntoNullableField {
			intoPointer = ""
		}

		return fmt.Sprintf(`resource, err := %[2]s.New(opts...)
		if err != nil {
			return err
		}

		builder.internal.%[1]s = %[3]sresource.Internal()
`, fieldPath, referredBuilderAlias, intoPointer)
	}

	if assignment.ArgumentName == "" {
		return fmt.Sprintf("builder.internal.%[1]s = %[2]s", fieldPath, formatScalar(assignment.Value))
	}

	argName := jenny.escapeVarName(tools.LowerCamelCase(assignment.ArgumentName))

	asPointer := ""
	// FIXME: this condition is probably wrong
	if valueType.Kind != ast.KindArray && valueType.Kind != ast.KindMap && assignment.IntoNullableField {
		asPointer = "&"
	}

	generatedConstraints := strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
	if generatedConstraints != "" {
		generatedConstraints += "\n\n"
	}

	return generatedConstraints + fmt.Sprintf("builder.internal.%[1]s = %[3]s%[2]s", fieldPath, argName, asPointer)
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
func (jenny *Builder) typeImportAlias(object ast.Object) string {
	// all types within a schema are generated under the same package
	return object.SelfRef.ReferredPkg
}

// importType declares an import statement for the type definition of
// the given object and returns an alias to it.
func (jenny *Builder) importType(object ast.Object) string {
	pkg := jenny.typeImportAlias(object)

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
