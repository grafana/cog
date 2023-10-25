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

type Type string

const (
	TypeType         = "type"
	TypeEnum         = "enum"
	TypeConst        = "const"
	TypeInterface    = "interface"
	TypeIntersection = "intersection"
	TypeEmpty        = ""
)

type RawTmpl struct {
	Imports importMap
	Objects []Object
}

type Object struct {
	Name         string
	Type         Type
	Comments     []string
	Value        any
	Fields       []Field
	HasDefault   bool
	DefaultValue string
}

type Field struct {
	Name     string
	Value    any
	Comments []string
	Fields   []Field
}

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

	objects := make([]Object, len(schema.Objects))
	for i, typeDef := range schema.Objects {
		typeDefGen, err := jenny.formatObject(schema, typeDef, packageMapper)
		if err != nil {
			return nil, err
		}
		objects[i] = typeDefGen
	}

	err := templates.Lookup("types.tmpl").Execute(&buffer, RawTmpl{
		Imports: imports,
		Objects: objects,
	})

	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}

func (jenny RawTypes) formatObject(schema *ast.Schema, def ast.Object, packageMapper pkgMapper) (Object, error) {
	objectType := TypeEmpty
	var typeValue any
	fields := make([]Field, 0)

	switch def.Type.Kind {
	case ast.KindStruct:
		objectType = TypeInterface
		fields = formatStruct(def.Type.AsStruct().Fields, packageMapper)
	case ast.KindEnum:
		objectType = TypeEnum
		fields = formatEnum(def.Type.AsEnum().Values)
	case ast.KindDisjunction, ast.KindMap, ast.KindArray, ast.KindRef:
		objectType = TypeType
		typeValue = formatType(def.Type, packageMapper)
	case ast.KindScalar:
		scalarType := def.Type.AsScalar()
		switch {
		case def.Type.Hints["kind"] == "type":
			objectType = TypeType
			typeValue = scalarType.Value
		case !scalarType.IsConcrete():
			objectType = TypeType
			typeValue = formatScalarKind(scalarType.ScalarKind)
		default:
			objectType = TypeConst
			typeValue = scalarType.Value
		}
	case ast.KindIntersection:
		objectType = TypeIntersection
		fields = formatIntersection(def.Type.AsIntersection(), packageMapper)
	default:
		return Object{}, fmt.Errorf("unhandled type def kind: %s", def.Type.Kind)
	}

	hasDefault := false
	defaultValue := ""
	// generate a "default value factory" for every object, except for constants
	if def.Type.Kind != ast.KindScalar || (def.Type.Kind == ast.KindScalar && !def.Type.AsScalar().IsConcrete()) {
		hasDefault = true
		// TODO: Improve this
		defaultValue = formatValue(defaultValueForObject(schema, def, packageMapper))
	}

	return Object{
		Name:         def.Name,
		Type:         Type(objectType),
		Comments:     def.Comments,
		Value:        typeValue,
		Fields:       fields,
		HasDefault:   hasDefault,
		DefaultValue: defaultValue,
	}, nil
}

func formatStruct(structFields []ast.StructField, packageMapper pkgMapper) []Field {
	fields := make([]Field, len(structFields))
	for i, field := range structFields {
		nestedFields := make([]Field, 0)
		if field.Type.Kind == ast.KindStruct {
			nestedFields = formatStruct(field.Type.AsStruct().Fields, packageMapper)
		}
		fields[i] = Field{
			Name:     field.Name,
			Value:    formatType(field.Type, packageMapper),
			Comments: field.Comments,
			Fields:   nestedFields,
		}
	}

	return fields
}

func formatEnum(values []ast.EnumValue) []Field {
	fields := make([]Field, len(values))
	for i, v := range values {
		fields[i] = Field{
			Name:  tools.CleanupNames(tools.UpperCamelCase(v.Name)),
			Value: v.Value,
		}
	}

	return fields
}

func formatReference(ref ast.RefType, packageMapper pkgMapper) string {
	referredPkg := packageMapper(ref.ReferredPkg)
	if referredPkg != "" {
		return referredPkg + "." + ref.ReferredType
	}

	return ref.ReferredType
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

func formatIntersection(def ast.IntersectionType, packageMapper pkgMapper) []Field {
	fields := make([]Field, 0, len(def.Branches))

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
		fields = append(fields, Field{
			Name: formatReference(ref.AsRef(), packageMapper),
		})
	}

	for _, r := range rest {
		if r.Struct != nil {
			fields = append(fields, formatStruct(r.AsStruct().Fields, packageMapper)...)
		}

		// Map different types ??
	}

	return fields
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
	case ast.KindIntersection:
		return defaultValuesForIntersection(schema, typeDef.AsIntersection(), packageMapper)
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

func defaultValuesForIntersection(schema *ast.Schema, intersectDef ast.IntersectionType, packageMapper pkgMapper) *orderedmap.Map[string, any] {
	defaults := orderedmap.New[string, any]()

	for _, branch := range intersectDef.Branches {
		if branch.Ref != nil {
			continue
		}

		if branch.Struct != nil {
			strctDef := defaultValuesForStructType(schema, branch.AsStruct(), packageMapper)
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
