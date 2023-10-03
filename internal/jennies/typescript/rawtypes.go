package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type pkgMapper func(string) string

type RawTypes struct {
}

func (jenny RawTypes) JennyName() string {
	return "TypescriptRawTypes"
}

func (jenny RawTypes) Generate(file *ast.File) (codejen.Files, error) {
	output, err := jenny.generateFile(file)
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("types/"+file.Package+"/types_gen.ts", output, jenny),
	}, nil
}

func (jenny RawTypes) generateFile(file *ast.File) ([]byte, error) {
	var buffer strings.Builder

	imports := newImportMap()

	packageMapper := func(pkg string) string {
		if pkg == file.Package {
			return ""
		}

		imports.Add(pkg, fmt.Sprintf("../%s/types_gen", pkg))

		return pkg
	}

	for _, typeDef := range file.Definitions {
		typeDefGen, err := jenny.formatObject(typeDef, packageMapper)
		if err != nil {
			return nil, err
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	}

	importStatements := imports.Format()
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(importStatements + buffer.String()), nil
}

func (jenny RawTypes) formatObject(def ast.Object, packageMapper pkgMapper) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString("export ")

	switch def.Type.Kind {
	case ast.KindStruct:
		buffer.WriteString(fmt.Sprintf("interface %s ", def.Name))
		buffer.WriteString(formatStructFields(def.Type.AsStruct().Fields, packageMapper))
		buffer.WriteString("\n")

		buffer.WriteString("\n")
		buffer.WriteString(formatStructDefaults(def))
		buffer.WriteString("\n")
	case ast.KindEnum:
		buffer.WriteString(fmt.Sprintf("enum %s {\n", def.Name))
		for _, val := range def.Type.AsEnum().Values {
			buffer.WriteString(fmt.Sprintf("\t%s = %s,\n", tools.UpperCamelCase(val.Name), formatScalar(val.Value)))
		}
		buffer.WriteString("}\n")
	case ast.KindDisjunction, ast.KindMap, ast.KindArray, ast.KindRef:
		buffer.WriteString(fmt.Sprintf("type %s = %s;\n", def.Name, formatType(def.Type, packageMapper)))
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()
		if scalarType.Value != nil {
			buffer.WriteString(fmt.Sprintf("const %s = %s;\n", def.Name, formatScalar(scalarType.Value)))
		} else {
			buffer.WriteString(fmt.Sprintf("type %s = %s;\n", def.Name, formatScalarKind(scalarType.ScalarKind)))
		}
	default:
		return nil, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	return []byte(buffer.String()), nil
}

func formatStructFields(fields []ast.StructField, packageMapper pkgMapper) string {
	var buffer strings.Builder

	buffer.WriteString("{\n")

	for _, fieldDef := range fields {
		fieldDefGen := formatField(fieldDef, packageMapper)

		buffer.WriteString(
			strings.TrimSuffix(
				prefixLinesWith(string(fieldDefGen), "\t"),
				"\t",
			),
		)
	}

	buffer.WriteString("}")

	return buffer.String()
}

func formatStructDefaults(def ast.Object) string {
	var buffer strings.Builder

	buffer.WriteString(fmt.Sprintf("export const default%[1]s: Partial<%[2]s> = {\n", tools.UpperCamelCase(def.Name), def.Name))

	fields := def.Type.AsStruct().Fields
	for _, field := range fields {
		if field.Default == nil {
			continue
		}

		buffer.WriteString(fmt.Sprintf("\t%s: %s,\n", field.Name, formatScalar(field.Default)))
	}

	buffer.WriteString("};")

	return buffer.String()
}

func formatField(def ast.StructField, packageMapper pkgMapper) []byte {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	required := ""
	if !def.Required {
		required = "?"
	}

	formattedType := formatType(def.Type, packageMapper)

	buffer.WriteString(fmt.Sprintf(
		"%s%s: %s;\n",
		def.Name,
		required,
		formattedType,
	))

	return []byte(buffer.String())
}

func formatType(def ast.Type, packageMapper pkgMapper) string {
	switch def.Kind {
	case ast.KindDisjunction:
		return formatDisjunction(def.AsDisjunction(), packageMapper)
	case ast.KindRef:
		referredPkg := packageMapper(def.AsRef().ReferredPkg)
		if referredPkg != "" {
			return referredPkg + "." + def.AsRef().ReferredType
		}

		return def.AsRef().ReferredType
	case ast.KindArray:
		return formatArray(def.AsArray(), packageMapper)
	case ast.KindStruct:
		return formatStructFields(def.AsStruct().Fields, packageMapper)
	case ast.KindMap:
		return formatMap(def.AsMap(), packageMapper)
	case ast.KindEnum:
		return formatAnonymousEnum(def.AsEnum())
	case ast.KindScalar:
		// This scalar actually refers to a constant
		if def.AsScalar().Value != nil {
			return formatScalar(def.AsScalar().Value)
		}

		return formatScalarKind(def.AsScalar().ScalarKind)

	default:
		return string(def.Kind)
	}
}

func formatScalarKind(kind ast.ScalarKind) string {
	switch kind {
	case ast.KindNull:
		return "null"
	case ast.KindAny:
		return "any"

	case ast.KindBytes, ast.KindString:
		return "string"

	case ast.KindFloat32, ast.KindFloat64:
		return "number"
	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return "number"
	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return "number"

	case ast.KindBool:
		return "boolean"
	default:
		return string(kind)
	}
}

func formatArray(def ast.ArrayType, packageMapper pkgMapper) string {
	subTypeString := formatType(def.ValueType, packageMapper)

	return fmt.Sprintf("%s[]", subTypeString)
}

func formatDisjunction(def ast.DisjunctionType, packageMapper pkgMapper) string {
	subTypes := make([]string, 0, len(def.Branches))
	for _, subType := range def.Branches {
		subTypes = append(subTypes, formatType(subType, packageMapper))
	}

	return strings.Join(subTypes, " | ")
}

func formatMap(def ast.MapType, packageMapper pkgMapper) string {
	keyTypeString := formatType(def.IndexType, packageMapper)
	valueTypeString := formatType(def.ValueType, packageMapper)

	return fmt.Sprintf("Record<%s, %s>", keyTypeString, valueTypeString)
}

func formatAnonymousEnum(def ast.EnumType) string {
	values := make([]string, 0, len(def.Values))
	for _, value := range def.Values {
		values = append(values, fmt.Sprintf("%#v", value.Value))
	}

	enumeration := strings.Join(values, " | ")

	return enumeration
}

func prefixLinesWith(input string, prefix string) string {
	lines := strings.Split(input, "\n")
	prefixed := make([]string, 0, len(lines))

	for _, line := range lines {
		prefixed = append(prefixed, prefix+line)
	}

	return strings.Join(prefixed, "\n")
}

func formatScalar(val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}
