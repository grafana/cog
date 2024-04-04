package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type Encoding struct {
}

func (jenny Encoding) JennyName() string {
	return "GoEncoding"
}

func (jenny Encoding) Generate(context ast.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}
		if output == nil {
			continue
		}

		filename := filepath.Join(
			formatPackageName(schema.Package),
			"types_go_encoding_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny Encoding) generateSchema(context ast.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder
	var err error

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		if context.ResolveToBuilder(object.AsRef()) || !jenny.typeNeedsEncoding(context, object.Type) {
			return
		}

		encode, innerErr := jenny.renderEncode(context, object)
		if innerErr != nil {
			err = innerErr
			return
		}
		buffer.WriteString(encode)
		buffer.WriteString("\n")
	})
	if err != nil {
		return nil, err
	}

	if buffer.Len() == 0 {
		return nil, nil
	}

	return []byte(fmt.Sprintf(`package %[1]s
%[2]s`, formatPackageName(schema.Package), buffer.String())), nil
}

func (jenny Encoding) typeNeedsEncoding(context ast.Context, typeDef ast.Type) bool {
	if context.ResolveToBuilder(typeDef) {
		return false
	}

	if typeDef.IsAny() || typeDef.IsScalar() {
		return false
	}

	if typeDef.IsArray() && typeDef.Array.IsArrayOfScalars() {
		return false
	}

	if typeDef.IsMap() && typeDef.Map.ValueType.IsScalar() {
		return false
	}

	return true
}

func (jenny Encoding) renderEncode(context ast.Context, obj ast.Object) (string, error) {
	if obj.Type.IsEnum() {
		return jenny.renderEnumEncode(obj)
	}

	if obj.Type.IsMap() {
		return jenny.renderMapEncode(context, obj)
	}

	if obj.Type.IsStruct() {
		return jenny.renderStructEncode(context, obj)
	}

	return "", nil

	//return "", fmt.Errorf("could not determine how to render encoding function for object of type %s", obj.Type.Kind)
}

func (jenny Encoding) renderStructEncode(context ast.Context, obj ast.Object) (string, error) {
	valueName := tools.LowerCamelCase(obj.Name)
	objectName := tools.UpperCamelCase(obj.Name)

	fields := make([]string, 0, len(obj.Type.Struct.Fields))
	for _, field := range obj.Type.Struct.Fields {
		fieldName := tools.UpperCamelCase(field.Name)

		valueFormatter := fmt.Sprintf(`fmt.Sprintf(\"%%#v\", %s.%s)`, valueName, fieldName) // TODO: what if value is a list?
		if jenny.typeNeedsEncoding(context, field.Type) {
			valueFormatter = fmt.Sprintf(`%s.%s.encodeToGo()`, valueName, fieldName)
		}

		fields = append(fields, fmt.Sprintf(`buffer.WriteString("%s: %s,\n")`, fieldName, valueFormatter))
	}

	return fmt.Sprintf(`func (%[1]s *%[2]s) encodeToGo() string {
	if %[1]s == nil {
		return "nil"
	}

	var buffer strings.Builder
	
	buffer.WriteString("\"")
	buffer.WriteString("%[2]s{\n")
	%[3]s
	buffer.WriteString("\n}\n")
	buffer.WriteString("\"")

	return buffer.String()
}`, valueName, objectName, strings.Join(fields, "\n")), nil
}

func (jenny Encoding) renderMapEncode(context ast.Context, obj ast.Object) (string, error) {
	valueName := tools.LowerCamelCase(obj.Name)
	objectName := tools.UpperCamelCase(obj.Name)

	valueFormatter := `fmt.Sprintf("%#v", value)` // TODO: what if value is a list?
	if jenny.typeNeedsEncoding(context, obj.Type.Map.ValueType) {
		valueFormatter = `value.encodeToGo()`
	}

	return fmt.Sprintf(`func (%[1]s *%[2]s) encodeToGo() string {
	if %[1]s == nil {
		return "nil"
	}

	var buffer strings.Builder
	
	buffer.WriteString("\"")
	buffer.WriteString("%[2]s{\n")

	for key, value := range *%[1]s {
		buffer.WriteString(fmt.Sprintf("%%#v: %%s,\n", key, %[3]s))
	}

	buffer.WriteString("\n}\n")
	buffer.WriteString("\"")

	return buffer.String()
}`, valueName, objectName, valueFormatter), nil
}

func (jenny Encoding) renderEnumEncode(obj ast.Object) (string, error) {
	valueName := tools.LowerCamelCase(obj.Name)
	objectName := tools.UpperCamelCase(obj.Name)

	cases := make([]string, 0, len(obj.Type.Enum.Values))

	for _, member := range obj.Type.Enum.Values {
		cases = append(cases, fmt.Sprintf(`if *%[1]s == %[2]s {
	return "%[2]s"
}
`, valueName, member.Name))
	}

	return fmt.Sprintf(`func (%[1]s *%[2]s) encodeToGo() string {
	if %[1]s == nil {
		return "nil"
	}

	%[3]s

	return "unknown"
}`, valueName, objectName, strings.Join(cases, "\n")), nil
}
