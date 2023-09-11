package golang

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type RawTypes struct {
}

func (jenny RawTypes) JennyName() string {
	return "GoRawTypes"
}

func (jenny RawTypes) Generate(file *ast.File) (codejen.Files, error) {
	output, err := jenny.generateFile(file)
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("types/"+file.Package+"/types_gen.go", output, jenny),
	}, nil
}

func (jenny RawTypes) generateFile(file *ast.File) ([]byte, error) {
	var buffer strings.Builder

	buffer.WriteString("package types\n\n")

	for _, object := range file.Definitions {
		objectOutput, err := jenny.formatObject(object)
		if err != nil {
			return nil, err
		}

		buffer.Write(objectOutput)
		buffer.WriteString("\n")

		// Add JSON (un)marshaling shortcuts
		if !object.Type.IsAny() {
			jsonMarshal, err := jenny.veneer("json_marshal", object)
			if err != nil {
				return nil, err
			}
			buffer.WriteString(jsonMarshal)
		}
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatObject(def ast.Object) ([]byte, error) {
	var buffer strings.Builder

	defName := tools.UpperCamelCase(def.Name)

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	switch def.Type.Kind {
	case ast.KindStruct:
		buffer.WriteString(fmt.Sprintf("type %s ", defName))
		buffer.WriteString(formatStructBody(def.Type.AsStruct(), ""))
		buffer.WriteString("\n")
	case ast.KindEnum:
		buffer.WriteString(jenny.formatEnumDef(def))
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()

		if scalarType.Value != nil {
			buffer.WriteString(fmt.Sprintf("const %s = %s", defName, formatScalar(scalarType.Value)))
		} else {
			buffer.WriteString(fmt.Sprintf("type %s %s", defName, scalarType.ScalarKind))
		}
	case ast.KindMap:
		buffer.WriteString(fmt.Sprintf("type %s %s", defName, formatType(def.Type, true, "")))
	case ast.KindRef:
		buffer.WriteString(fmt.Sprintf("type %s %s", defName, def.Type.AsRef().ReferredType))
	default:
		return nil, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatEnumDef(def ast.Object) string {
	var buffer strings.Builder

	enumName := tools.UpperCamelCase(def.Name)
	enumType := def.Type.AsEnum()

	buffer.WriteString(fmt.Sprintf("type %s %s\n", enumName, formatType(enumType.Values[0].Type, true, "")))

	buffer.WriteString("const (\n")
	for _, val := range enumType.Values {
		buffer.WriteString(fmt.Sprintf("\t%s %s = %#v\n", tools.UpperCamelCase(val.Name), enumName, val.Value))
	}
	buffer.WriteString(")\n")

	return buffer.String()
}

func (jenny RawTypes) veneer(veneerType string, def ast.Object) (string, error) {
	// First, see if there is a definition-specific veneer
	templateFile := fmt.Sprintf("%s.types.%s.go.tmpl", strings.ToLower(def.Name), veneerType)
	tmpl := templates.Lookup(templateFile)

	// If not, get the generic one
	if tmpl == nil {
		tmpl = templates.Lookup(fmt.Sprintf("types.%s.go.tmpl", veneerType))
	}
	// If not, something went wrong.
	if tmpl == nil {
		return "", fmt.Errorf("veneer '%s' not found", veneerType)
	}

	buf := bytes.Buffer{}
	if err := tmpl.Execute(&buf, map[string]any{
		"def": def,
	}); err != nil {
		return "", fmt.Errorf("failed executing veneer template: %w", err)
	}

	return buf.String(), nil
}

func formatStructBody(def ast.StructType, typesPkg string) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	for _, fieldDef := range def.Fields {
		buffer.WriteString("\t" + formatField(fieldDef, typesPkg))
	}

	buffer.WriteString("}")

	return buffer.String()
}

func formatField(def ast.StructField, typesPkg string) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	// ToDo: this doesn't follow references to other types like the builder jenny does
	/*
		if def.Type.Default != nil {
			buffer.WriteString(fmt.Sprintf("// Default: %#v\n", def.Type.Default))
		}
	*/

	jsonOmitEmpty := ""
	if !def.Required {
		jsonOmitEmpty = ",omitempty"
	}

	buffer.WriteString(fmt.Sprintf(
		"%s %s `json:\"%s%s\"`\n",
		tools.UpperCamelCase(def.Name),
		formatType(def.Type, def.Required, typesPkg),
		def.Name,
		jsonOmitEmpty,
	))

	return buffer.String()
}

func formatType(def ast.Type, fieldIsRequired bool, typesPkg string) string {
	if def.IsAny() {
		return "any"
	}

	if def.Kind == ast.KindDisjunction {
		return formatDisjunction(def.AsDisjunction(), typesPkg)
	}

	if def.Kind == ast.KindArray {
		return formatArray(def.AsArray(), typesPkg)
	}

	if def.Kind == ast.KindMap {
		return formatMap(def.AsMap(), typesPkg)
	}

	if def.Kind == ast.KindScalar {
		typeName := def.AsScalar().ScalarKind
		if !fieldIsRequired {
			typeName = "*" + typeName
		}

		return string(typeName)
	}

	if def.Kind == ast.KindRef {
		typeName := def.AsRef().ReferredType

		if typesPkg != "" {
			typeName = typesPkg + "." + typeName
		}

		if !fieldIsRequired {
			typeName = "*" + typeName
		}

		return typeName
	}

	if def.Kind == ast.KindEnum {
		return "enum here"
	}

	// anonymous struct
	if def.Kind == ast.KindStruct {
		return formatStructBody(def.AsStruct(), typesPkg)
	}

	// FIXME: we shouldn't be here
	return "unknown"
}

func formatArray(def ast.ArrayType, typesPkg string) string {
	subTypeString := formatType(def.ValueType, true, typesPkg)

	return fmt.Sprintf("[]%s", subTypeString)
}

func formatMap(def ast.MapType, typesPkg string) string {
	keyTypeString := formatType(def.IndexType, true, typesPkg)
	valueTypeString := formatType(def.ValueType, true, typesPkg)

	return fmt.Sprintf("map[%s]%s", keyTypeString, valueTypeString)
}

func formatDisjunction(def ast.DisjunctionType, typesPkg string) string {
	subTypes := make([]string, 0, len(def.Branches))
	for _, subType := range def.Branches {
		subTypes = append(subTypes, formatType(subType, true, typesPkg))
	}

	return fmt.Sprintf("disjunction<%s>", strings.Join(subTypes, " | "))
}

func formatScalar(val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		// TODO: we can't assume a list of strings
		return fmt.Sprintf("[]string{%s}", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}
