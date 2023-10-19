package typescript

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type raw string

type pkgMapper func(string) string

type RawTypes struct {
}

func (jenny RawTypes) JennyName() string {
	return "TypescriptRawTypes"
}

func (jenny RawTypes) Generate(schema *ast.Schema) (codejen.Files, error) {
	output, err := jenny.generateSchema(schema)
	if err != nil {
		return nil, err
	}

	filename := filepath.Join(
		"types",
		strings.ToLower(schema.Package),
		"types_gen.ts",
	)

	return codejen.Files{
		*codejen.NewFile(filename, output, jenny),
	}, nil
}

func (jenny RawTypes) generateSchema(schema *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	imports := newImportMap()

	packageMapper := func(pkg string) string {
		if pkg == schema.Package {
			return ""
		}

		imports.Add(pkg, fmt.Sprintf("../%s/types_gen", pkg))

		return pkg
	}

	for _, typeDef := range schema.Objects {
		typeDefGen, err := jenny.formatObject(schema, typeDef, packageMapper)
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

func (jenny RawTypes) formatObject(schema *ast.Schema, def ast.Object, packageMapper pkgMapper) ([]byte, error) {
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
	case ast.KindEnum:
		buffer.WriteString(fmt.Sprintf("enum %s {\n", def.Name))
		for _, val := range def.Type.AsEnum().Values {
			name := tools.CleanupNames(tools.UpperCamelCase(val.Name))
			buffer.WriteString(fmt.Sprintf("\t%s = %s,\n", name, formatScalar(val.Value)))
		}
		buffer.WriteString("}\n")
	case ast.KindDisjunction, ast.KindMap, ast.KindArray, ast.KindRef:
		buffer.WriteString(fmt.Sprintf("type %s = %s;\n", def.Name, formatType(def.Type, packageMapper)))
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()
		typeValue := formatScalar(scalarType.Value)

		if !scalarType.IsConcrete() || def.Type.Hints["kind"] == "type" {
			if !scalarType.IsConcrete() {
				typeValue = formatScalarKind(scalarType.ScalarKind)
			}

			buffer.WriteString(fmt.Sprintf("type %s = %s;\n", def.Name, typeValue))
		} else {
			buffer.WriteString(fmt.Sprintf("const %s = %s;\n", def.Name, typeValue))
		}
	case ast.KindIntersection:
		buffer.WriteString(fmt.Sprintf("interface %s ", def.Name))
		buffer.WriteString(formatIntersection(def.Type.AsIntersection(), packageMapper))
		buffer.WriteString("\n")
	default:
		return nil, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	// generate a "default value factory" for every object, except for constants
	if def.Type.Kind != ast.KindScalar || (def.Type.Kind == ast.KindScalar && !def.Type.AsScalar().IsConcrete()) {
		buffer.WriteString("\n")

		buffer.WriteString(fmt.Sprintf("export const default%[1]s = (): %[2]s => (", tools.UpperCamelCase(def.Name), def.Name))

		formattedDefaults := formatValue(defaultValueForObject(schema, def, packageMapper))
		buffer.WriteString(formattedDefaults)

		buffer.WriteString(");\n")
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
				prefixLinesWith(fieldDefGen, "\t"),
				"\t",
			),
		)
	}

	buffer.WriteString("}")

	return buffer.String()
}

func formatField(def ast.StructField, packageMapper pkgMapper) string {
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

	return buffer.String()
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
	case ast.KindIntersection:
		return formatIntersection(def.AsIntersection(), packageMapper)
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

	if def.ValueType.Kind == ast.KindDisjunction {
		return fmt.Sprintf("(%s)[]", subTypeString)
	}

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
		return defaultValuesForStructType(schema, typeDef.AsStruct(), packageMapper)
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

	default:
		return "unknown"
	}
}

func defaultValuesForStructType(schema *ast.Schema, structDef ast.StructType, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, field := range structDef.Fields {
		if field.Type.Default != nil {
			defaults.Set(field.Name, field.Type.Default)
			continue
		}

		if !field.Required {
			continue
		}

		defaults.Set(field.Name, defaultValueForType(schema, field.Type, packageMapper))
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

func formatValue(val any) string {
	if rawVal, ok := val.(raw); ok {
		return string(rawVal)
	}

	var buffer strings.Builder

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

func formatIntersection(def ast.IntersectionType, packageMapper pkgMapper) string {
	var buffer strings.Builder

	refs := make([]ast.Type, 0)
	rest := make([]ast.Type, 0)
	for _, b := range def.Branches {
		if b.Ref != nil {
			refs = append(refs, b)
			continue
		}
		rest = append(rest, b)
	}

	if len(refs) > 0 {
		buffer.WriteString("extends ")
	}

	for i, ref := range refs {
		if i != 0 && i < len(refs) {
			buffer.WriteString(", ")
		}

		buffer.WriteString(formatType(ref, packageMapper))
	}

	buffer.WriteString(" {\n")

	for _, r := range rest {
		if r.Struct != nil {
			for _, fieldDef := range r.AsStruct().Fields {
				buffer.WriteString("\t" + formatField(fieldDef, packageMapper))
			}
			continue
		}
		buffer.WriteString("\t" + formatType(r, packageMapper))
	}

	buffer.WriteString("}")

	return buffer.String()
}
