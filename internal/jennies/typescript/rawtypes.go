package typescript

import (
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

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

	for _, typeDef := range file.Definitions {
		typeDefGen, err := jenny.formatObject(typeDef, "")
		if err != nil {
			return nil, err
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatObject(def ast.Object, typesPkg string) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString("export ")

	switch def.Type.Kind {
	case ast.KindStruct:
		buffer.WriteString(fmt.Sprintf("interface %s ", def.Name))
		buffer.WriteString(formatStructFields(def.Type.AsStruct().Fields, typesPkg))
		buffer.WriteString("\n")
	case ast.KindEnum:
		buffer.WriteString(fmt.Sprintf("enum %s {\n", def.Name))
		for _, val := range def.Type.AsEnum().Values {
			name := tools.CleanupNames(tools.UpperCamelCase(val.Name))
			buffer.WriteString(fmt.Sprintf("\t%s%s = %s,\n", def.Name, name, formatScalar(val.Value)))
		}
		buffer.WriteString("}\n")
	case ast.KindRef:
		refType := def.Type.AsRef()

		buffer.WriteString(fmt.Sprintf("type %s = %s;", def.Name, refType.ReferredType))
	case ast.KindDisjunction, ast.KindMap, ast.KindArray:
		buffer.WriteString(fmt.Sprintf("type %s = %s;\n", def.Name, formatType(def.Type, "")))
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

func formatStructFields(fields []ast.StructField, typesPkg string) string {
	var buffer strings.Builder

	buffer.WriteString("{\n")

	for i, fieldDef := range fields {
		fieldDefGen := formatField(fieldDef, typesPkg)

		buffer.WriteString(
			strings.TrimSuffix(
				prefixLinesWith(string(fieldDefGen), "\t"),
				"\n\t",
			),
		)

		if i != len(fields)-1 {
			buffer.WriteString("\n")
		}
	}

	buffer.WriteString("\n}")

	return buffer.String()
}

func formatField(def ast.StructField, typesPkg string) []byte {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	required := ""
	if !def.Required {
		required = "?"
	}

	formattedType := formatType(def.Type, typesPkg)

	buffer.WriteString(fmt.Sprintf(
		"%s%s: %s;\n",
		tools.LowerCamelCase(def.Name),
		required,
		formattedType,
	))

	return []byte(buffer.String())
}

func formatType(def ast.Type, typesPkg string) string {
	switch def.Kind {
	case ast.KindDisjunction:
		return formatDisjunction(def.AsDisjunction(), typesPkg)
	case ast.KindRef:
		if typesPkg != "" {
			return typesPkg + "." + def.AsRef().ReferredType
		}

		return def.AsRef().ReferredType
	case ast.KindArray:
		return formatArray(def.AsArray(), typesPkg)
	case ast.KindStruct:
		return formatStructFields(def.AsStruct().Fields, typesPkg)
	case ast.KindMap:
		return formatMap(def.AsMap(), typesPkg)
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

func formatArray(def ast.ArrayType, typesPkg string) string {
	subTypeString := formatType(def.ValueType, typesPkg)

	return fmt.Sprintf("%s[]", subTypeString)
}

func formatDisjunction(def ast.DisjunctionType, typesPkg string) string {
	subTypes := make([]string, 0, len(def.Branches))
	for _, subType := range def.Branches {
		subTypes = append(subTypes, formatType(subType, typesPkg))
	}

	return strings.Join(subTypes, " | ")
}

func formatMap(def ast.MapType, typesPkg string) string {
	keyTypeString := formatType(def.IndexType, typesPkg)
	valueTypeString := formatType(def.ValueType, typesPkg)

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
