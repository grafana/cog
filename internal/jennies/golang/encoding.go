package golang

import (
	"fmt"
	"path/filepath"
	"strings"

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

		buffer.WriteString(fmt.Sprintf(`func (resource *%s) EncodeToGo() string {
	if resource == nil {
		return "nil"
	}
`, tools.UpperCamelCase(object.Name)))

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

func typeHasEncoder(context languages.Context, typeDef ast.Type) bool {
	if typeNeedsEncoding(context, typeDef) {
		return true
	}

	if !typeDef.IsRef() {
		return false
	}

	obj, found := context.LocateObject(typeDef.Ref.ReferredPkg, typeDef.Ref.ReferredType)
	if !found {
		return false
	}

	if context.ResolveToBuilder(obj.AsRef()) {
		return false
	}

	return typeHasEncoder(context, obj.Type)
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

	return "", fmt.Errorf("could not determine how to render encoding function for object of type %s", obj.Type.Kind)
}

func typeNameToGo(typeDef ast.Type) string {
	switch typeDef.Kind {
	case ast.KindArray:
		return fmt.Sprintf("[]%s", typeNameToGo(typeDef.Array.ValueType))
	case ast.KindMap:
		return fmt.Sprintf("map[%s, %s]", typeNameToGo(typeDef.Map.IndexType), typeNameToGo(typeDef.Map.ValueType))
	case ast.KindRef:
		return fmt.Sprintf("%s.%s", formatPackageName(typeDef.Ref.ReferredPkg), tools.UpperCamelCase(typeDef.Ref.ReferredType))
	case ast.KindScalar:
		return string(typeDef.Scalar.ScalarKind)
	}

	return "unknown"
}

func (jenny Encoding) renderStructEncode(context languages.Context, obj ast.Object) (string, error) {
	objectName := tools.UpperCamelCase(obj.Name)

	fields := make([]string, 0, len(obj.Type.Struct.Fields))
	for _, field := range obj.Type.Struct.Fields {
		fieldName := tools.UpperCamelCase(field.Name)
		valuePath := fmt.Sprintf("resource.%s", fieldName)

		// TODO: what if value is a list/map?
		if field.Type.IsArray() && !field.Type.Array.IsArrayOfScalars() {
			continue
		}

		if field.Type.IsMap() {
			continue
		}

		if typeHasEncoder(context, field.Type) {
			fields = append(fields, fmt.Sprintf(`buffer.WriteString(fmt.Sprintf("%s: %%s,\n", cog.Dump(%s)))`, fieldName, valuePath))
			continue
		}

		val := valuePath
		encodedVal := "%#v"
		if field.Type.Nullable && !field.Type.IsAnyOf(ast.KindArray, ast.KindMap) && !field.Type.IsAny() {
			val = "*" + val
			encodedVal = fmt.Sprintf("cog.ToPtr[%s](%%#v)", typeNameToGo(field.Type))
		}

		assignment := fmt.Sprintf(`buffer.WriteString(fmt.Sprintf("%[1]s: %[2]s,\n", %[3]s))`, fieldName, encodedVal, val)
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

	return buffer.String()`, formatPackageName(obj.SelfRef.ReferredPkg), objectName, strings.Join(fields, "\n")), nil
}

func (jenny Encoding) renderMapEncode(context languages.Context, obj ast.Object) (string, error) {
	objectName := tools.UpperCamelCase(obj.Name)

	valueFormatter := `fmt.Sprintf("%#v", value)` // TODO: what if value is a list?
	if typeNeedsEncoding(context, obj.Type.Map.ValueType) {
		valueFormatter = `value.EncodeToGo()`
	}

	return fmt.Sprintf(`var buffer strings.Builder
	
	buffer.WriteString("\"")
	buffer.WriteString("%[1]s{\n")

	for key, value := range *resource {
		buffer.WriteString(fmt.Sprintf("%%#v: %%s,\n", key, %[2]s))
	}

	buffer.WriteString("\n}\n")
	buffer.WriteString("\"")

	return buffer.String()`, objectName, valueFormatter), nil
}

func (jenny Encoding) renderEnumEncode(obj ast.Object) (string, error) {
	cases := make([]string, 0, len(obj.Type.Enum.Values))

	for _, member := range obj.Type.Enum.Values {
		cases = append(cases, fmt.Sprintf(`if *resource == %[2]s {
	return "%[1]s.%[2]s"
}
`, formatPackageName(obj.SelfRef.ReferredPkg), member.Name))
	}

	return fmt.Sprintf(`%[1]s

	return "unknown"`, strings.Join(cases, "\n")), nil
}
