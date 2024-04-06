package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Encoding struct {
}

func (jenny Encoding) JennyName() string {
	return "GoEncoding"
}

func (jenny Encoding) Generate(context languages.Context) (codejen.Files, error) {
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

func (jenny Encoding) generateSchema(context languages.Context, schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder
	var err error

	schema.Objects.Iterate(func(_ string, object ast.Object) {
		if context.ResolveToBuilder(object.AsRef()) || !typeNeedsEncoding(context, object.Type) {
			return
		}

		encode, innerErr := jenny.renderEncode(context, object)
		if innerErr != nil {
			err = innerErr
			return
		}

		valueName := tools.LowerCamelCase(object.Name)
		objectName := tools.UpperCamelCase(object.Name)

		buffer.WriteString(fmt.Sprintf(`func (%[1]s *%[2]s) EncodeToGo() string {
	if %[1]s == nil {
		return "nil"
	}
`, valueName, objectName))

		buffer.WriteString(encode)

		buffer.WriteString("\n}\n\n")
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

func typeNeedsEncoding(context languages.Context, typeDef ast.Type) bool {
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

	if typeDef.IsRef() {
		return false
	}

	return true
}

func (jenny Encoding) renderEncode(context languages.Context, obj ast.Object) (string, error) {
	if obj.Type.IsEnum() {
		return jenny.renderEnumEncode(obj)
	}

	if obj.Type.IsMap() {
		return jenny.renderMapEncode(context, obj)
	}

	if obj.Type.IsStruct() {
		return jenny.renderStructEncode(context, obj)
	}

	spew.Dump(obj)
	return "", fmt.Errorf("could not determine how to render encoding function for object of type %s", obj.Type.Kind)

	//return "", fmt.Errorf("could not determine how to render encoding function for object of type %s", obj.Type.Kind)
}

func (jenny Encoding) renderStructEncode(context languages.Context, obj ast.Object) (string, error) {
	valueName := tools.LowerCamelCase(obj.Name)
	objectName := tools.UpperCamelCase(obj.Name)

	fields := make([]string, 0, len(obj.Type.Struct.Fields))
	for _, field := range obj.Type.Struct.Fields {
		fieldName := tools.UpperCamelCase(field.Name)
		valuePath := fmt.Sprintf("%s.%s", valueName, fieldName)

		if typeNeedsEncoding(context, field.Type) {
			fields = append(fields, fmt.Sprintf(`buffer.WriteString("%s: %s.EncodeToGo(),\n")`, fieldName, valuePath))
			continue
		}

		if field.Type.IsAny() {
			// TODO: what to do here?
			continue
		}

		val := valuePath
		encodedVal := "%#v"
		if field.Type.Nullable && !field.Type.IsAnyOf(ast.KindArray, ast.KindMap) {
			val = "*" + val
			encodedVal = "cog.ToPtr(%#v)"
		}

		assignment := fmt.Sprintf(`buffer.WriteString(fmt.Sprintf("%[1]s: %[2]s,\n", %[3]s))`, fieldName, encodedVal, val) // TODO: what if value is a list?
		if field.Type.Nullable {
			assignment = fmt.Sprintf(`if %s != nil {
	%s
}`, valuePath, assignment)
		}

		fields = append(fields, assignment)
	}

	return fmt.Sprintf(`var buffer strings.Builder
	
	buffer.WriteString("%[1]s.%[2]s{\n")
	%[3]s
	buffer.WriteString("\n}")

	return buffer.String()`, obj.SelfRef.ReferredPkg, objectName, strings.Join(fields, "\n")), nil
}

func (jenny Encoding) renderMapEncode(context languages.Context, obj ast.Object) (string, error) {
	valueName := tools.LowerCamelCase(obj.Name)
	objectName := tools.UpperCamelCase(obj.Name)

	valueFormatter := `fmt.Sprintf("%#v", value)` // TODO: what if value is a list?
	if typeNeedsEncoding(context, obj.Type.Map.ValueType) {
		valueFormatter = `value.EncodeToGo()`
	}

	return fmt.Sprintf(`var buffer strings.Builder
	
	buffer.WriteString("\"")
	buffer.WriteString("%[2]s{\n")

	for key, value := range *%[1]s {
		buffer.WriteString(fmt.Sprintf("%%#v: %%s,\n", key, %[3]s))
	}

	buffer.WriteString("\n}\n")
	buffer.WriteString("\"")

	return buffer.String()`, valueName, objectName, valueFormatter), nil
}

func (jenny Encoding) renderEnumEncode(obj ast.Object) (string, error) {
	valueName := tools.LowerCamelCase(obj.Name)

	cases := make([]string, 0, len(obj.Type.Enum.Values))

	for _, member := range obj.Type.Enum.Values {
		cases = append(cases, fmt.Sprintf(`if *%[1]s == %[3]s {
	return "%[2]s.%[3]s"
}
`, valueName, obj.SelfRef.ReferredPkg, member.Name))
	}

	return fmt.Sprintf(`%[1]s

	return "unknown"`, strings.Join(cases, "\n")), nil
}
