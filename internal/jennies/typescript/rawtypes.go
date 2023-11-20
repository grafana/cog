package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type raw string

type pkgMapper func(string) string

type RawTypes struct {
	typeFormatter *typeFormatter
	schemas       ast.Schemas
}

func (jenny RawTypes) JennyName() string {
	return "TypescriptRawTypes"
}

func (jenny RawTypes) Generate(context common.Context) (codejen.Files, error) {
	jenny.schemas = context.Schemas
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
		typeDefGen, err := jenny.formatObject(typeDef, packageMapper)
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

func (jenny RawTypes) formatObject(def ast.Object, packageMapper pkgMapper) ([]byte, error) {
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

		formattedDefaults := formatValue(jenny.defaultValueForObject(def, packageMapper))
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

func (jenny RawTypes) defaultValueForObject(object ast.Object, packageMapper pkgMapper) any {
	switch object.Type.Kind {
	case ast.KindEnum:
		return defaultValueForEnumType(object.Name, object.Type, nil)
	default:
		return jenny.defaultValueForType(object.Type, packageMapper)
	}
}

func (jenny RawTypes) defaultValueForType(typeDef ast.Type, packageMapper pkgMapper) any {
	if typeDef.Default != nil {
		return typeDef.Default
	}

	switch typeDef.Kind {
	case ast.KindDisjunction:
		return jenny.defaultValueForType(typeDef.AsDisjunction().Branches[0], packageMapper)
	case ast.KindStruct:
		return jenny.defaultValuesForStructType(typeDef, packageMapper)
	case ast.KindEnum: // anonymous enum
		defaultValue := typeDef.AsEnum().Values[0].Value
		if typeDef.Default != nil {
			defaultValue = typeDef.Default
		}

		return defaultValue
	case ast.KindRef:
		return jenny.defaultValuesForReference(typeDef, packageMapper)
	case ast.KindMap:
		return raw("{}")
	case ast.KindArray:
		return raw("[]")
	case ast.KindScalar:
		return defaultValueForScalar(typeDef.AsScalar())
	case ast.KindIntersection:
		return jenny.defaultValuesForIntersection(typeDef.AsIntersection(), packageMapper)
	default:
		return "unknown"
	}
}

func (jenny RawTypes) defaultValuesForStructType(structType ast.Type, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, field := range structType.AsStruct().Fields {
		if field.Type.Default != nil {
			if field.Type.Kind == ast.KindRef {
				defaults.Set(field.Name, jenny.defaultValuesForReference(field.Type, packageMapper))
				continue
			}
			defaults.Set(field.Name, field.Type.Default)
			continue
		}

		if !field.Required {
			continue
		}

		defaults.Set(field.Name, jenny.defaultValueForType(field.Type, packageMapper))
	}

	if structType.ImplementsVariant() {
		variant := tools.UpperCamelCase(structType.ImplementedVariant())
		defaults.Set("_implements"+variant+"Variant", raw("() => {}"))
	}

	return defaults
}

func defaultValueForEnumType(name string, typeDef ast.Type, overrideDefault any) any {
	enum := typeDef.AsEnum()
	defaultValue := enum.Values[0].Value
	if typeDef.Default != nil {
		defaultValue = typeDef.Default
	}
	if overrideDefault != nil {
		defaultValue = overrideDefault
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

func (jenny RawTypes) defaultValuesForIntersection(intersectDef ast.IntersectionType, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, branch := range intersectDef.Branches {
		if branch.Ref != nil {
			continue
		}

		if branch.Struct != nil {
			strctDef := jenny.defaultValuesForStructType(branch, packageMapper)
			strctDef.Iterate(func(key string, value any) {
				defaults.Set(key, value)
			})
		}

		// TODO: Add them for other types?
	}

	return defaults
}

func (jenny RawTypes) defaultValuesForReference(typeDef ast.Type, packageMapper pkgMapper) any {
	ref := typeDef.AsRef()

	pkg := packageMapper(ref.ReferredPkg)
	referredType, _ := jenny.schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)

	// is the reference to a constant?
	if referredType.Type.Kind == ast.KindScalar && referredType.Type.AsScalar().IsConcrete() {
		return raw(fmt.Sprintf("%s.%s", ref.ReferredPkg, ref.ReferredType))
	}

	if referredType.Type.Kind == ast.KindEnum {
		if pkg != "" {
			return raw(fmt.Sprintf("%s.%s", pkg, defaultValueForEnumType(ref.ReferredType, referredType.Type, typeDef.Default)))
		}

		return defaultValueForEnumType(ref.ReferredType, referredType.Type, typeDef.Default)
	}

	if pkg != "" {
		return raw(fmt.Sprintf("%s.default%s()", ref.ReferredPkg, ref.ReferredType))
	}

	return raw(fmt.Sprintf("default%s()", ref.ReferredType))
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
