package typescript

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

type Tmpl struct {
	Package     string
	Name        string
	Imports     importMap
	ImportAlias string
	Options     []string
	Constructor constructor
}

type constructor struct {
	Args         []string
	Items        map[string]any
	Initializers []string
}

func (jenny *Builder) JennyName() string {
	return "TypescriptBuilder"
}

func (jenny *Builder) Generate(builders []ast.Builder) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range builders {
		output, err := jenny.generateBuilder(builders, builder)
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

func (jenny *Builder) generateBuilder(builders ast.Builders, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	jenny.imports = newImportMap()
	importAlias := jenny.importType(builder.For)

	objectName := builder.For.Name
	constructorCode := jenny.generateConstructor(builders, builder)

	// Define options
	options := make([]string, len(builder.Options))
	for i, option := range builder.Options {
		options[i] = jenny.generateOption(builders, builder, option)
	}

	tmpl := templates.Lookup("builder.tmpl")
	err := tmpl.Execute(&buffer, Tmpl{
		Package:     builder.Package,
		Name:        objectName,
		Imports:     jenny.imports,
		ImportAlias: importAlias,
		Options:     options,
		Constructor: constructorCode,
	})

	return []byte(buffer.String()), err
}

func (jenny *Builder) generateConstructor(builders ast.Builders, builder ast.Builder) constructor {
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

	return constructor{
		Args:         argsList,
		Items:        jenny.getDefaultValues(builders, builder.Package, builder.For.Type),
		Initializers: fieldsInitList,
	}
}

func (jenny *Builder) getDefaultValues(builders ast.Builders, pkg string, typeDef ast.Type) map[string]any {
	switch typeDef.Kind {
	case ast.KindDisjunction:
		return jenny.getDefaultValues(builders, pkg, typeDef.AsDisjunction().Branches[0])
	case ast.KindRef:
		ref := typeDef.AsRef()
		referredTypeBuilder, _ := builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)
		return jenny.getDefaultValues(builders, referredTypeBuilder.Package, referredTypeBuilder.For.Type)
	case ast.KindStruct:
		return jenny.emptyValueForStruct(builders, pkg, typeDef.AsStruct())
	default:
		return map[string]any{"": "unknown"}
	}
}

func (jenny *Builder) emptyValueForType(builders ast.Builders, pkg string, typeDef ast.Type) string {
	switch typeDef.Kind {
	case ast.KindDisjunction:
		return jenny.emptyValueForType(builders, pkg, typeDef.AsDisjunction().Branches[0])
	case ast.KindRef:
		ref := typeDef.AsRef()
		referredTypeBuilder, _ := builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)
		return jenny.emptyValueForType(builders, referredTypeBuilder.Package, referredTypeBuilder.For.Type)
	case ast.KindEnum:
		return jenny.formatEnumDefault(typeDef.AsEnum().Values)
	case ast.KindMap:
		return "{}"
	case ast.KindArray:
		return "[]"
	case ast.KindScalar:
		return jenny.emptyValueForScalar(typeDef.AsScalar())

	default:
		return "unknown"
	}
}

func (jenny *Builder) emptyValueForStruct(builders ast.Builders, pkg string, structType ast.StructType) map[string]any {
	fieldsInit := make(map[string]any, len(structType.Fields))
	for _, field := range structType.Fields {
		if field.Type.Default != nil {
			fieldsInit[field.Name] = formatScalar(field.Type.Default)
			continue
		}

		if !field.Required {
			continue
		}

		if field.Type.Kind == ast.KindStruct {
			return jenny.emptyValueForStruct(builders, pkg, field.Type.AsStruct())
		}

		fieldsInit[field.Name] = jenny.emptyValueForType(builders, pkg, field.Type)
	}

	return fieldsInit
}

func (jenny *Builder) formatEnumDefault(values []ast.EnumValue) string {
	for _, v := range values {
		if v.Type.Default != nil {
			return formatScalar(v.Type.Default)
		}
	}
	return formatScalar(values[0].Value)
}

func (jenny *Builder) emptyValueForScalar(scalar ast.ScalarType) string {
	switch scalar.ScalarKind {
	case ast.KindNull:
		return "null"
	case ast.KindAny:
		return "{}"

	case ast.KindBytes, ast.KindString:
		return "''"

	case ast.KindFloat32, ast.KindFloat64:
		return "0"

	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return "0"

	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return "0"

	case ast.KindBool:
		return "false"

	default:
		return "undefined"
	}
}

func (jenny *Builder) builderForType(builders ast.Builders, builder ast.Builder, t ast.Type) (ast.Builder, bool) {
	if t.Kind != ast.KindRef {
		return ast.Builder{}, false
	}

	// TODO: shouldn't we using the package from the ref?!
	ref := t.AsRef()
	return builders.LocateByObject(builder.Package, ref.ReferredType)
}

func (jenny *Builder) generateInitAssignment(builders ast.Builders, builder ast.Builder, assignment ast.Assignment) string {
	fieldPath := assignment.Path

	if _, valueHasBuilder := jenny.builderForType(builders, builder, assignment.ValueType); valueHasBuilder {
		return "constructor init assignment with type that has a builder is not supported yet"
	}

	if assignment.ArgumentName == "" {
		return fmt.Sprintf("\t\tthis.internal.%[1]s = %[2]s;", fieldPath, formatScalar(assignment.Value))
	}

	argName := tools.LowerCamelCase(assignment.ArgumentName)

	generatedConstraints := strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
	if generatedConstraints != "" {
		generatedConstraints += "\n\n"
	}

	return generatedConstraints + fmt.Sprintf("this.internal.%[1]s = %[2]s;", fieldPath, argName)
}

func (jenny *Builder) generateOption(builders ast.Builders, builder ast.Builder, def ast.Option) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("\t// %s\n", commentLine))
	}

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

	// Option body
	buffer.WriteString(fmt.Sprintf(`	%[1]s(%[2]s): this {
%[3]s

		return this;
	}

`, def.Name, arguments, assignments))

	return buffer.String()
}

func (jenny *Builder) generateArgument(builders ast.Builders, builder ast.Builder, arg ast.Argument) string {
	if referredBuilder, found := jenny.builderForType(builders, builder, arg.Type); found {
		referredTypeAlias := jenny.typeImportAlias(referredBuilder.For)

		return fmt.Sprintf(`%[1]s: OptionsBuilder<%[2]s.%[3]s>`, arg.Name, referredTypeAlias, referredBuilder.For.Name)
	}

	builderImportAlias := jenny.typeImportAlias(builder.For)
	typeName := formatType(arg.Type, func(pkg string) string {
		if pkg == builder.For.SelfRef.ReferredPkg {
			return builderImportAlias
		}

		jenny.imports.Add(pkg, fmt.Sprintf("../../types/%s/types_gen", pkg))

		return pkg
	})

	name := tools.LowerCamelCase(arg.Name)

	return fmt.Sprintf("%s: %s", name, typeName)
}

func (jenny *Builder) generateAssignment(builders ast.Builders, builder ast.Builder, assignment ast.Assignment) string {
	fieldPath := assignment.Path

	if _, found := jenny.builderForType(builders, builder, assignment.ValueType); found {
		return fmt.Sprintf("\t\tthis.internal.%[1]s = %[2]s.build();", fieldPath, assignment.ArgumentName)
	}

	if assignment.ArgumentName == "" {
		return fmt.Sprintf("\t\tthis.internal.%[1]s = %[2]s;", fieldPath, formatScalar(assignment.Value))
	}

	argName := tools.LowerCamelCase(assignment.ArgumentName)

	generatedConstraints := strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
	if generatedConstraints != "" {
		generatedConstraints += "\n"
	}

	return generatedConstraints + fmt.Sprintf("\t\tthis.internal.%[1]s = %[2]s;", fieldPath, argName)
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
