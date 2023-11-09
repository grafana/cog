package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

type pkgMapper func(string) string

type RawTypes struct {
}

func (jenny RawTypes) JennyName() string {
	return "GoRawTypes"
}

func (jenny RawTypes) Generate(context context.Builders) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			strings.ToLower(schema.Package),
			"types_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder
	imports := template.NewImportMap()

	packageMapper := func(pkg string) string {
		if pkg == schema.Package {
			return ""
		}

		return imports.Add(pkg, "github.com/grafana/cog/generated/"+pkg)
	}

	for _, object := range schema.Objects {
		objectOutput, err := jenny.formatObject(object, packageMapper)
		if err != nil {
			return nil, err
		}

		buffer.Write(objectOutput)
		buffer.WriteString("\n")
	}

	importStatements := formatImports(imports)
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(fmt.Sprintf(`package %[1]s

%[2]s%[3]s`, strings.ToLower(schema.Package), importStatements, buffer.String())), nil
}

func (jenny RawTypes) formatObject(def ast.Object, packageMapper pkgMapper) ([]byte, error) {
	var buffer strings.Builder

	defName := tools.UpperCamelCase(def.Name)

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	switch def.Type.Kind {
	case ast.KindEnum:
		buffer.WriteString(jenny.formatEnumDef(def, packageMapper))
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()

		//nolint: gocritic
		if scalarType.Value != nil {
			buffer.WriteString(fmt.Sprintf("const %s = %s", defName, formatScalar(scalarType.Value)))
		} else if scalarType.ScalarKind == ast.KindBytes {
			buffer.WriteString(fmt.Sprintf("type %s %s", defName, "[]byte"))
		} else {
			buffer.WriteString(fmt.Sprintf("type %s %s", defName, formatType(def.Type, packageMapper)))
		}
	case ast.KindMap, ast.KindRef, ast.KindArray, ast.KindStruct, ast.KindIntersection:
		buffer.WriteString(fmt.Sprintf("type %s %s", defName, formatType(def.Type, packageMapper)))
	default:
		return nil, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	buffer.WriteString("\n")

	if def.Type.Kind == ast.KindStruct && def.Type.ImplementsVariant() {
		variant := tools.UpperCamelCase(def.Type.ImplementedVariant())

		buffer.WriteString(fmt.Sprintf("func (resource %s) Implements%sVariant() {}\n", defName, variant))
		buffer.WriteString("\n")
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatEnumDef(def ast.Object, packageMapper pkgMapper) string {
	var buffer strings.Builder

	enumName := tools.UpperCamelCase(def.Name)
	enumType := def.Type.AsEnum()

	buffer.WriteString(fmt.Sprintf("type %s %s\n", enumName, formatType(enumType.Values[0].Type, packageMapper)))

	buffer.WriteString("const (\n")
	for _, val := range enumType.Values {
		name := tools.CleanupNames(tools.UpperCamelCase(val.Name))
		buffer.WriteString(fmt.Sprintf("\t%s %s = %#v\n", name, enumName, val.Value))
	}
	buffer.WriteString(")\n")

	return buffer.String()
}

func formatStructBody(def ast.StructType, packageMapper pkgMapper) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	for _, fieldDef := range def.Fields {
		buffer.WriteString("\t" + formatField(fieldDef, packageMapper))
	}

	buffer.WriteString("}")

	return buffer.String()
}

func formatField(def ast.StructField, packageMapper pkgMapper) string {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	jsonOmitEmpty := ""
	if !def.Required {
		jsonOmitEmpty = ",omitempty"
	}

	buffer.WriteString(fmt.Sprintf(
		"%s %s `json:\"%s%s\"`\n",
		tools.UpperCamelCase(def.Name),
		formatType(def.Type, packageMapper),
		def.Name,
		jsonOmitEmpty,
	))

	return buffer.String()
}

func formatType(def ast.Type, packageMapper pkgMapper) string {
	if def.IsAny() {
		return "any"
	}

	if def.Kind == ast.KindComposableSlot {
		return variantInterface(string(def.AsComposableSlot().Variant), packageMapper)
	}

	if def.Kind == ast.KindArray {
		return formatArray(def.AsArray(), packageMapper)
	}

	if def.Kind == ast.KindMap {
		return formatMap(def.AsMap(), packageMapper)
	}

	if def.Kind == ast.KindScalar {
		typeName := def.AsScalar().ScalarKind
		if def.Nullable {
			typeName = "*" + typeName
		}

		return string(typeName)
	}

	if def.Kind == ast.KindRef {
		return formatRef(def, packageMapper)
	}

	// anonymous struct or struct body
	if def.Kind == ast.KindStruct {
		output := formatStructBody(def.AsStruct(), packageMapper)
		if def.Nullable {
			output = "*" + output
		}

		return output
	}

	if def.Kind == ast.KindIntersection {
		return formatIntersection(def.AsIntersection(), packageMapper)
	}

	// FIXME: we should never be here
	return "unknown"
}

func variantInterface(variant string, packageMapper pkgMapper) string {
	referredPkg := packageMapper("cog/variants")

	return fmt.Sprintf("%s.%s", referredPkg, tools.UpperCamelCase(variant))
}

func formatArray(def ast.ArrayType, packageMapper pkgMapper) string {
	subTypeString := formatType(def.ValueType, packageMapper)

	return fmt.Sprintf("[]%s", subTypeString)
}

func formatMap(def ast.MapType, packageMapper pkgMapper) string {
	keyTypeString := formatType(def.IndexType, packageMapper)
	valueTypeString := formatType(def.ValueType, packageMapper)

	return fmt.Sprintf("map[%s]%s", keyTypeString, valueTypeString)
}

func formatScalar(val any) string {
	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		// FIXME: this is wrong, we can't just assume a list of strings.
		return fmt.Sprintf("[]string{%s}", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}

func formatRef(def ast.Type, packageMapper pkgMapper) string {
	referredPkg := packageMapper(def.AsRef().ReferredPkg)
	typeName := tools.UpperCamelCase(def.AsRef().ReferredType)

	if referredPkg != "" {
		typeName = referredPkg + "." + typeName
	}

	if def.Nullable {
		typeName = "*" + typeName
	}

	return typeName
}

func formatIntersection(def ast.IntersectionType, packageMapper pkgMapper) string {
	var buffer strings.Builder

	buffer.WriteString("struct {\n")

	refs := make([]ast.Type, 0)
	rest := make([]ast.Type, 0)
	for _, b := range def.Branches {
		if b.Ref != nil {
			refs = append(refs, b)
			continue
		}

		rest = append(rest, b)
	}

	for _, ref := range refs {
		buffer.WriteString("\t" + formatRef(ref, packageMapper) + "\n")
	}

	if len(refs) > 0 {
		buffer.WriteString("\n")
	}

	for _, r := range rest {
		if r.Struct != nil {
			for _, fieldDef := range r.AsStruct().Fields {
				buffer.WriteString("\t" + formatField(fieldDef, packageMapper))
			}
			continue
		}
		buffer.WriteString("\t" + formatType(r, packageMapper) + "\n")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func formatImports(importMap map[string]string) string {
	if len(importMap) == 0 {
		return ""
	}

	statements := make([]string, 0, len(importMap))

	for alias, importPath := range importMap {
		statements = append(statements, fmt.Sprintf(`	%s "%s"`, alias, importPath))
	}

	return fmt.Sprintf(`import (
%[1]s
)`, strings.Join(statements, "\n"))
}
