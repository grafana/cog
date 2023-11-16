package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type raw string

type pkgMapper func(string) string

type RawTypes struct {
	typeFormatter *typeFormatter
}

func (jenny RawTypes) JennyName() string {
	return "TypescriptRawTypes"
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
			"types_gen.ts",
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

		return imports.Add(pkg, fmt.Sprintf("../%s", pkg))
	}

	jenny.typeFormatter = defaultTypeFormatter(packageMapper)

	for _, typeDef := range schema.Objects {
		typeDefGen, err := jenny.formatObject(schema, typeDef, packageMapper)
		if err != nil {
			return nil, err
		}

		buffer.Write(typeDefGen)
		buffer.WriteString("\n")
	}

	importStatements := formatImports(imports)
	if importStatements != "" {
		importStatements += "\n\n"
	}

	return []byte(importStatements + buffer.String()), nil
}

func (jenny RawTypes) formatObject(schema *ast.Schema, def ast.Object, packageMapper pkgMapper) ([]byte, error) {
	var buffer strings.Builder

	for _, commentLine := range def.Comments {
		buffer.WriteString(fmt.Sprintf("// %s\n", commentLine))
	}

	buffer.WriteString("export ")

	switch def.Type.Kind {
	case ast.KindStruct:
		buffer.WriteString(fmt.Sprintf("interface %s ", def.Name))
		buffer.WriteString(jenny.typeFormatter.formatStructFields(def.Type))
		buffer.WriteString("\n")
	case ast.KindEnum:
		buffer.WriteString(fmt.Sprintf("enum %s {\n", def.Name))
		for _, val := range def.Type.AsEnum().Values {
			name := tools.CleanupNames(tools.UpperCamelCase(val.Name))
			buffer.WriteString(fmt.Sprintf("\t%s = %s,\n", name, formatScalar(val.Value)))
		}
		buffer.WriteString("}\n")
	case ast.KindDisjunction, ast.KindMap, ast.KindArray, ast.KindRef:
		buffer.WriteString(fmt.Sprintf("type %s = %s;\n", def.Name, jenny.typeFormatter.formatType(def.Type)))
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()
		typeValue := formatScalar(scalarType.Value)

		if !scalarType.IsConcrete() || def.Type.Hints["kind"] == "type" {
			if !scalarType.IsConcrete() {
				typeValue = jenny.typeFormatter.formatScalarKind(scalarType.ScalarKind)
			}

			buffer.WriteString(fmt.Sprintf("type %s = %s;\n", def.Name, typeValue))
		} else {
			buffer.WriteString(fmt.Sprintf("const %s = %s;\n", def.Name, typeValue))
		}
	case ast.KindIntersection:
		buffer.WriteString(fmt.Sprintf("interface %s ", def.Name))
		buffer.WriteString(jenny.typeFormatter.formatType(def.Type))
		buffer.WriteString("\n")
	case ast.KindComposableSlot:
		buffer.WriteString(fmt.Sprintf("interface %s %s\n", def.Name, jenny.typeFormatter.variantInterface(string(def.Type.AsComposableSlot().Variant))))
	default:
		return nil, fmt.Errorf("unhandled object of type: %s", def.Type.Kind)
	}
	// generate a "default value factory" for every object, except for constants or composability slots
	if (def.Type.Kind != ast.KindScalar && def.Type.Kind != ast.KindComposableSlot) || (def.Type.Kind == ast.KindScalar && !def.Type.AsScalar().IsConcrete()) {
		buffer.WriteString("\n")

		buffer.WriteString(fmt.Sprintf("export const default%[1]s = (): %[2]s => (", tools.UpperCamelCase(def.Name), def.Name))

		formattedDefaults := formatValue(defaultValueForObject(schema, def, packageMapper))
		buffer.WriteString(formattedDefaults)

		buffer.WriteString(");\n")
	}

	return []byte(buffer.String()), nil
}

func prefixLinesWith(input string, prefix string) string {
	lines := strings.Split(input, "\n")
	prefixed := make([]string, 0, len(lines))

	for _, line := range lines {
		prefixed = append(prefixed, prefix+line)
	}

	return strings.Join(prefixed, "\n")
}

/******************************************************************************
* 					 Default and "empty" values management 					  *
******************************************************************************/

func defaultValueForObject(schema *ast.Schema, object ast.Object, packageMapper pkgMapper) any {
	switch object.Type.Kind {
	case ast.KindEnum:
		return defaultValueForEnumType(object.Name, object.Type)
	default:
		return defaultValueForType(schema, object.Type, packageMapper)
	}
}

func defaultValueForType(schema *ast.Schema, typeDef ast.Type, packageMapper pkgMapper) any {
	if typeDef.Default != nil {
		return typeDef.Default
	}

	switch typeDef.Kind {
	case ast.KindDisjunction:
		return defaultValueForType(schema, typeDef.AsDisjunction().Branches[0], packageMapper)
	case ast.KindStruct:
		return defaultValuesForStructType(schema, typeDef, packageMapper)
	case ast.KindEnum: // anonymous enum
		return typeDef.AsEnum().Values[0].Value
	case ast.KindRef:
		ref := typeDef.AsRef()

		// TODO: handle references to other packages
		referredType := schema.LocateObject(ref.ReferredType)
		// is the reference to a constant?
		if referredType.Type.Kind == ast.KindScalar && referredType.Type.AsScalar().IsConcrete() {
			return raw(fmt.Sprintf("%s.%s", ref.ReferredPkg, ref.ReferredType))
		}

		pkg := packageMapper(ref.ReferredPkg)
		if pkg != "" {
			return raw(fmt.Sprintf("%s.default%s()", ref.ReferredPkg, ref.ReferredType))
		}

		return raw(fmt.Sprintf("default%s()", ref.ReferredType))
	case ast.KindMap:
		return raw("{}")
	case ast.KindArray:
		return raw("[]")
	case ast.KindScalar:
		return defaultValueForScalar(typeDef.AsScalar())
	case ast.KindIntersection:
		return defaultValuesForIntersection(schema, typeDef.AsIntersection(), packageMapper)
	default:
		return "unknown"
	}
}

func defaultValuesForStructType(schema *ast.Schema, structType ast.Type, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, field := range structType.AsStruct().Fields {
		if field.Type.Default != nil {
			defaults.Set(field.Name, field.Type.Default)
			continue
		}

		if !field.Required {
			continue
		}

		defaults.Set(field.Name, defaultValueForType(schema, field.Type, packageMapper))
	}

	if structType.ImplementsVariant() {
		variant := tools.UpperCamelCase(structType.ImplementedVariant())
		defaults.Set("_implements"+variant+"Variant", raw("() => {}"))
	}

	return defaults
}

func defaultValueForEnumType(name string, typeDef ast.Type) any {
	enum := typeDef.AsEnum()
	defaultValue := enum.Values[0].Value
	if typeDef.Default != nil {
		defaultValue = typeDef.Default
	}

	for _, v := range enum.Values {
		if v.Value == defaultValue {
			return raw(fmt.Sprintf("%s.%s", name, tools.CleanupNames(tools.UpperCamelCase(v.Name))))
		}
	}

	return raw(fmt.Sprintf("%s.%s", name, tools.CleanupNames(tools.UpperCamelCase(enum.Values[0].Name))))
}

func defaultValueForScalar(scalar ast.ScalarType) any {
	// The scalar represents a constant
	if scalar.Value != nil {
		return scalar.Value
	}

	switch scalar.ScalarKind {
	case ast.KindNull:
		return raw("null")
	case ast.KindAny:
		return raw("{}")

	case ast.KindBytes, ast.KindString:
		return ""

	case ast.KindFloat32, ast.KindFloat64:
		return 0.0

	case ast.KindUint8, ast.KindUint16, ast.KindUint32, ast.KindUint64:
		return 0

	case ast.KindInt8, ast.KindInt16, ast.KindInt32, ast.KindInt64:
		return 0

	case ast.KindBool:
		return false

	default:
		return "unknown"
	}
}

func defaultValuesForIntersection(schema *ast.Schema, intersectDef ast.IntersectionType, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, branch := range intersectDef.Branches {
		if branch.Ref != nil {
			continue
		}

		if branch.Struct != nil {
			strctDef := defaultValuesForStructType(schema, branch, packageMapper)
			strctDef.Iterate(func(key string, value any) {
				defaults.Set(key, value)
			})
		}

		// TODO: Add them for other types?
	}

	return defaults
}

func formatValue(val any) string {
	if rawVal, ok := val.(raw); ok {
		return string(rawVal)
	}

	var buffer strings.Builder

	if array, ok := val.([]any); ok {
		buffer.WriteString("[\n")
		for _, v := range array {
			buffer.WriteString(fmt.Sprintf("%s,\n", formatValue(v)))
		}
		buffer.WriteString("]")

		return buffer.String()
	}

	if orderedMap, ok := val.(*orderedmap.Map[string, any]); ok {
		buffer.WriteString("{\n")

		orderedMap.Iterate(func(key string, value any) {
			buffer.WriteString(fmt.Sprintf("\t%s: %s,\n", key, formatValue(value)))
		})

		buffer.WriteString("}")

		return buffer.String()
	}

	return fmt.Sprintf("%#v", val)
}

func formatImports(importMap template.ImportMap) string {
	if len(importMap) == 0 {
		return ""
	}

	statements := make([]string, 0, len(importMap))

	for alias, importPath := range importMap {
		statements = append(statements, fmt.Sprintf(`import * as %s from "%s";`, alias, importPath))
	}

	return strings.Join(statements, "\n")
}
