package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	imports importMap
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

		files = append(
			files,
			*codejen.NewFile(builder.Package+"/"+strings.ToLower(builder.For.Name)+"/builder_gen.ts", output, jenny),
		)
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(builders ast.Builders, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	jenny.imports = newImportMap()

	objectName := tools.UpperCamelCase(builder.For.Name)

	// imports
	jenny.imports.Add("types", fmt.Sprintf("../../types/%s/types_gen", builder.Package))
	buffer.WriteString("import { CogOptionsBuilder } from \"../../options_builder_gen\";\n\n")

	// Builder class declaration
	buffer.WriteString(fmt.Sprintf("export class %[1]sBuilder implements CogOptionsBuilder<types.%[1]s> {\n", objectName))

	// internal property, representing the object being built
	buffer.WriteString(fmt.Sprintf("\tprivate internal: types.%[1]s;\n", objectName))

	// Add a constructor for the builder
	constructorCode := jenny.generateConstructor(builders, builder)
	buffer.WriteString(constructorCode)

	// Allow builders to expose the resource they're building
	buffer.WriteString(fmt.Sprintf(`
	build(): types.%s {
		return this.internal;
	}

`, objectName))

	// Define options
	for _, option := range builder.Options {
		buffer.WriteString(jenny.generateOption(builders, builder, option))
	}

	// End builder class declaration
	buffer.WriteString("}\n")

	importStatements := jenny.imports.Format()
	if importStatements != "" {
		importStatements += "\n"
	}

	return []byte(importStatements + buffer.String()), nil
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

	args = strings.Join(argsList, ", ")
	fieldsInit = strings.Join(fieldsInitList, "\n")
	typeInit := jenny.emptyValueForType(builders, builder.Package, builder.For.Type)

	buffer.WriteString(fmt.Sprintf(`
	constructor(%[2]s) {
		this.internal = %[3]s;
%[4]s
	}
`, typeName, args, typeInit, fieldsInit))

	return buffer.String()
}

func (jenny *Builder) emptyValueForType(builders ast.Builders, pkg string, typeDef ast.Type) string {
	switch typeDef.Kind {
	case ast.KindDisjunction:
		return jenny.emptyValueForType(builders, pkg, typeDef.AsDisjunction().Branches[0])
	case ast.KindRef:
		ref := typeDef.AsRef()
		referredTypeBuilder, _ := builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)

		return jenny.emptyValueForType(builders, referredTypeBuilder.Package, referredTypeBuilder.For.Type)
	case ast.KindStruct:
		return jenny.emptyValueForStruct(builders, pkg, typeDef.AsStruct())
	case ast.KindEnum:
		return formatScalar(typeDef.AsEnum().Values[0].Value)
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

func (jenny *Builder) emptyValueForStruct(builders ast.Builders, pkg string, structType ast.StructType) string {
	var buffer strings.Builder

	var fieldsInit []string
	for _, field := range structType.Fields {
		if field.Type.Default != nil {
			fieldsInit = append(fieldsInit, fmt.Sprintf("%s: %s, // default value", field.Name, formatScalar(field.Type.Default)))
			continue
		}

		if !field.Required {
			continue
		}

		fieldsInit = append(fieldsInit, fmt.Sprintf("%s: %s, // zero value", field.Name, jenny.emptyValueForType(builders, pkg, field.Type)))
	}

	buffer.WriteString(fmt.Sprintf(`{
%[1]s
}`, strings.Join(fieldsInit, "\n")))

	return buffer.String()
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
		return "unknown"
	}
}

func (jenny *Builder) typeHasBuilder(builders ast.Builders, builder ast.Builder, t ast.Type) (string, bool) {
	if t.Kind != ast.KindRef {
		return "", false
	}

	ref := t.AsRef()
	_, builderFound := builders.LocateByObject(builder.Package, ref.ReferredType)

	return ref.ReferredPkg, builderFound
}

func (jenny *Builder) generateInitAssignment(builders ast.Builders, builder ast.Builder, assignment ast.Assignment) string {
	fieldPath := assignment.Path

	if _, valueHasBuilder := jenny.typeHasBuilder(builders, builder, assignment.ValueType); valueHasBuilder {
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

	return generatedConstraints + fmt.Sprintf("\t\tthis.internal.%[1]s = %[2]s;", fieldPath, argName)
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
	typeName := formatType(arg.Type, func(pkg string) string {
		if pkg == builder.Package {
			return "types"
		}

		jenny.imports.Add(pkg, fmt.Sprintf("../../types/%s/types_gen", pkg))

		return pkg
	})

	if referredTypeName, found := jenny.typeHasBuilder(builders, builder, arg.Type); found {
		return fmt.Sprintf(`%[1]s: OptionsBuilder<types.%[2]s>`, arg.Name, referredTypeName)
	}

	name := tools.LowerCamelCase(arg.Name)

	return fmt.Sprintf("%s: %s", name, typeName)
}

func (jenny *Builder) generateAssignment(builders ast.Builders, builder ast.Builder, assignment ast.Assignment) string {
	fieldPath := assignment.Path

	if _, found := jenny.typeHasBuilder(builders, builder, assignment.ValueType); found {
		return fmt.Sprintf("\t\tthis.internal.%[1]s = %[2]s.build();", fieldPath, assignment.ArgumentName)
	}

	if assignment.ArgumentName == "" {
		return fmt.Sprintf("\t\tthis.internal.%[1]s = %[2]s;", fieldPath, formatScalar(assignment.Value))
	}

	argName := tools.LowerCamelCase(assignment.ArgumentName)

	generatedConstraints := strings.Join(jenny.constraints(argName, assignment.Constraints), "\n")
	if generatedConstraints != "" {
		generatedConstraints += "\n\n"
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
