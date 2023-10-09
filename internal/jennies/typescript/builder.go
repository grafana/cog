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
	Package     string
	Name        string
	Imports     importMap
	ImportAlias string
	Options     []options
	Constructor constructor
}

type constructor struct {
	Args         []argument
	Items        map[string]any
	Initializers []string
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
	Constraints    []string
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
			strings.ToLower(builder.RootPackage),
			strings.ToLower(builder.Package),
			"builder_gen.ts",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(context context.Builders, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	jenny.imports = newImportMap()
	importAlias := jenny.importType(builder.For)

	constructorCode := jenny.generateConstructor(context, builder)

	// Define options
	opts := make([]options, len(builder.Options))
	for i, o := range builder.Options {
		opts[i] = jenny.generateOption(context, builder, o)
	}

	err := templates.Lookup("builder.tmpl").Execute(&buffer, Tmpl{
		Package:     builder.Package,
		Name:        builder.For.Name,
		Imports:     jenny.imports,
		ImportAlias: importAlias,
		Options:     opts,
		Constructor: constructorCode,
	})

	return []byte(buffer.String()), err
}

func (jenny *Builder) generateConstructor(context context.Builders, builder ast.Builder) constructor {
	var argsList []argument
	var fieldsInitList []string
	for _, opt := range builder.Options {
		if !opt.IsConstructorArg {
			continue
		}

		// FIXME: this is assuming that there's only one argument for that option
		argsList = append(argsList, jenny.generateArgument(context, builder, opt.Args[0]))
		fieldsInitList = append(
			fieldsInitList,
			jenny.generateInitAssignment(context, opt.Assignments[0]),
		)
	}

	for _, init := range builder.Initializations {
		fieldsInitList = append(
			fieldsInitList,
			jenny.generateInitAssignment(context, init),
		)
	}

	return constructor{
		Args:         argsList,
		Initializers: fieldsInitList,
	}
}

func (jenny *Builder) generateInitAssignment(context context.Builders, assignment ast.Assignment) string {
	fieldPath := assignment.Path
	valueType := assignment.Path.Last().Type

	if _, valueHasBuilder := context.BuilderForType(valueType); valueHasBuilder {
		return "constructor init assignment with type that has a builder is not supported yet"
	}

	if assignment.ArgumentName == "" {
		return fmt.Sprintf("this.internal.%[1]s = %[2]s;", fieldPath, formatScalar(assignment.Value))
	}

	argName := tools.LowerCamelCase(assignment.ArgumentName)

	generatedConstraints := strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
	if generatedConstraints != "" {
		generatedConstraints += "\n\n"
	}

	return generatedConstraints + fmt.Sprintf("this.internal.%[1]s = %[2]s;", fieldPath, argName)
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
		referredTypeAlias := jenny.typeImportAlias(referredBuilder.For)

		return argument{
			Name:          arg.Name,
			ReferredAlias: referredTypeAlias,
			ReferredName:  referredBuilder.For.Name,
		}
	}

	builderImportAlias := jenny.typeImportAlias(builder.For)
	typeName := formatType(arg.Type, func(pkg string) string {
		if pkg == builder.For.SelfRef.ReferredPkg {
			return builderImportAlias
		}

		jenny.imports.Add(pkg, fmt.Sprintf("../../types/%s/types_gen", pkg))

		return pkg
	})

	return argument{Name: tools.LowerCamelCase(arg.Name), Type: typeName}
}

func (jenny *Builder) generatePathInitializationSafeGuard(currentBuilder ast.Builder, path ast.Path) string {
	fieldPath := jenny.formatFieldPath(path, currentBuilder)
	valueType := path.Last().Type
	if path.Last().TypeHint != nil {
		valueType = *path.Last().TypeHint
	}

	emptyValue := defaultValueForType(currentBuilder.Schema, valueType, func(pkg string) string {
		jenny.imports.Add(pkg, fmt.Sprintf("../../types/%s/types_gen", pkg))

		return pkg
	})

	return fmt.Sprintf(`		if (!this.internal.%[1]s) {
			this.internal.%[1]s = %[2]s;
		}`, fieldPath, emptyValue)
}

func (jenny *Builder) generateAssignment(context context.Builders, builder ast.Builder, assign ast.Assignment) assignment {
	fieldPath := jenny.formatFieldPath(assign.Path, builder)
	pathEnd := assign.Path.Last()
	valueType := pathEnd.Type

	var pathInitSafeGuards []string
	for i, chunk := range assign.Path {
		if i == len(assign.Path)-1 {
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

	assignmentSafeGuards := strings.Join(pathInitSafeGuards, "\n")
	if assignmentSafeGuards != "" {
		assignmentSafeGuards = assignmentSafeGuards + "\n\n"
	}

	if _, found := context.BuilderForType(valueType); found {
		return assignment{
			Path:           fieldPath,
			InitSafeguards: pathInitSafeGuards,
			Value:          assign.ArgumentName,
			IsBuilder:      true,
		}
	}

	if assign.ArgumentName == "" {
		return assignment{
			Path:           fieldPath,
			InitSafeguards: pathInitSafeGuards,
			Value:          formatScalar(assign.Value),
		}
	}

	argName := tools.LowerCamelCase(assign.ArgumentName)

	return assignment{
		Path:           fieldPath,
		InitSafeguards: pathInitSafeGuards,
		Value:          argName,
		Constraints:    jenny.constraints(argName, assign.Constraints),
	}
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

	buffer.WriteString(fmt.Sprintf("\t\tif (!(%s)) {\n", jenny.constraintComparison(argumentName, constraint)))
	buffer.WriteString(fmt.Sprintf("\t\t\tthrow new Error(\"%[1]s must be %[2]s %[3]v\");\n", argumentName, constraint.Op, constraint.Args[0]))
	buffer.WriteString("\t\t}\n")

	return buffer.String()
}

func (jenny *Builder) constraintComparison(argumentName string, constraint ast.TypeConstraint) string {
	if constraint.Op == ast.MinLengthOp {
		return fmt.Sprintf("%[1]s.length >= %[2]v", argumentName, constraint.Args[0])
	}
	if constraint.Op == ast.MaxLengthOp {
		return fmt.Sprintf("%[1]s.length <= %[2]v", argumentName, constraint.Args[0])
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

	jenny.imports.Add(pkg, fmt.Sprintf("../../types/%s/types_gen", pkg))

	return pkg
}

func (jenny *Builder) formatFieldPath(fieldPath ast.Path, currentBuilder ast.Builder) string {
	if len(fieldPath) != 0 {
		return fieldPath.String()
	}
	formattedPath := ""

	builderImportAlias := jenny.typeImportAlias(currentBuilder.For)
	for i := range fieldPath {
		output := fieldPath[i].Identifier

		// don't generate type hints if:
		// * there isn't one defined
		// * the type isn't "any"
		// * as a trailing element in the path
		if !fieldPath[i].Type.IsAny() || fieldPath[i].TypeHint == nil || i == len(fieldPath)-1 {
			formattedPath += "." + output
			continue
		}

		formattedTypeHint := formatType(*fieldPath[i].TypeHint, func(pkg string) string {
			if pkg == currentBuilder.For.SelfRef.ReferredPkg {
				return builderImportAlias
			}

			jenny.imports.Add(pkg, fmt.Sprintf("../../types/%s/types_gen", pkg))

			return pkg
		})

		formattedPath += "(" + formattedPath + " as " + formattedTypeHint + ")." + output
	}

	return formattedPath
}
