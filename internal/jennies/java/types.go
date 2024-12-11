package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

const fasterXMLPackageName = "com.fasterxml.jackson"
const javaNullableField = "@JsonInclude(JsonInclude.Include.NON_NULL)"
const javaDefaultEmptyField = "@JsonSetter(nulls = Nulls.AS_EMPTY)"
const javaEmptyField = "@JsonInclude(JsonInclude.Include.NON_EMPTY)"

type typeFormatter struct {
	config        Config
	packageMapper func(pkg string, class string) string
	context       languages.Context
}

func createFormatter(ctx languages.Context, config Config) *typeFormatter {
	return &typeFormatter{context: ctx, config: config}
}

func (tf *typeFormatter) withPackageMapper(packageMapper func(pkg string, class string) string) *typeFormatter {
	tf.packageMapper = packageMapper
	return tf
}

func (tf *typeFormatter) formatFieldType(def ast.Type) string {
	switch def.Kind {
	case ast.KindScalar:
		return formatScalarType(def.AsScalar())
	case ast.KindRef:
		return tf.formatReference(def.AsRef())
	case ast.KindArray:
		return tf.formatArray(def.AsArray())
	case ast.KindComposableSlot:
		return tf.formatComposable(def.AsComposableSlot())
	case ast.KindMap:
		return tf.formatMap(def.AsMap())
	case ast.KindStruct:
		// TODO: Manage anonymous structs
		return "Object"
	}

	return "unknown"
}

func (tf *typeFormatter) typeHasBuilder(def ast.Type) bool {
	return tf.context.ResolveToBuilder(def)
}

func (tf *typeFormatter) resolvesToComposableSlot(typeDef ast.Type) bool {
	_, found := tf.context.ResolveToComposableSlot(typeDef)
	return found
}

func (tf *typeFormatter) formatBuilderFieldType(def ast.Type) string {
	value := tf.formatFieldType(def)
	if tf.resolvesToComposableSlot(def) || tf.typeHasBuilder(def) {
		value = fmt.Sprintf("%s.Builder<%s>", tf.formatPackage("cog"), value)
	}

	return value
}

func (tf *typeFormatter) formatReference(def ast.RefType) string {
	object, _ := tf.context.LocateObject(def.ReferredPkg, def.ReferredType)
	switch object.Type.Kind {
	case ast.KindScalar:
		return formatScalarType(object.Type.AsScalar())
	case ast.KindMap:
		return tf.formatMap(object.Type.AsMap())
	case ast.KindArray:
		return tf.formatArray(object.Type.AsArray())
	default:
		tf.packageMapper(def.ReferredPkg, def.ReferredType)
		return def.ReferredType
	}
}

func (tf *typeFormatter) formatArray(def ast.ArrayType) string {
	tf.packageMapper("java.util", "List")
	return fmt.Sprintf("List<%s>", tf.formatFieldType(def.ValueType))
}

func (tf *typeFormatter) formatMap(def ast.MapType) string {
	tf.packageMapper("java.util", "Map")
	mapType := "unknown"
	switch def.ValueType.Kind {
	case ast.KindRef:
		mapType = tf.formatReference(def.ValueType.AsRef())
	case ast.KindScalar:
		mapType = formatScalarType(def.ValueType.AsScalar())
	case ast.KindMap:
		mapType = tf.formatMap(def.ValueType.AsMap())
	case ast.KindArray:
		mapType = tf.formatArray(def.ValueType.AsArray())
	}

	return fmt.Sprintf("Map<String, %s>", mapType)
}

func (tf *typeFormatter) formatComposable(def ast.ComposableSlotType) string {
	variant := tools.UpperCamelCase(string(def.Variant))
	tf.packageMapper("cog.variants", variant)
	return variant
}

func formatScalarType(def ast.ScalarType) string {
	scalarType := "unknown"

	switch def.ScalarKind {
	case ast.KindString:
		scalarType = "String"
	case ast.KindBytes:
		scalarType = "Byte"
	case ast.KindInt16, ast.KindUint16:
		scalarType = "Short"
	case ast.KindInt8, ast.KindUint8, ast.KindInt32, ast.KindUint32:
		scalarType = "Integer"
	case ast.KindInt64, ast.KindUint64:
		scalarType = "Long"
	case ast.KindFloat32:
		scalarType = "Float"
	case ast.KindFloat64:
		scalarType = "Double"
	case ast.KindBool:
		scalarType = "Boolean"
	case ast.KindAny:
		scalarType = "Object"
	}

	return scalarType
}

func (tf *typeFormatter) defaultValueFor(def ast.Type) string {
	switch def.Kind {
	case ast.KindArray:
		tf.packageMapper("java.util", "LinkedList")
		return "new LinkedList<>()"
	case ast.KindMap:
		tf.packageMapper("java.util", "HashMap")
		return "new Hashmap<>()"
	case ast.KindRef:
		refDef := fmt.Sprintf("%s.%s", def.AsRef().ReferredPkg, def.AsRef().ReferredType)
		if tf.typeHasBuilder(def) {
			return fmt.Sprintf("new %sBuilder().build()", tf.formatPackage(refDef))
		}
		return fmt.Sprintf("new %s()", tf.formatPackage(refDef))
	case ast.KindStruct:
		return "new Object()"
	case ast.KindScalar:
		switch def.AsScalar().ScalarKind {
		case ast.KindBool:
			return "false"
		case ast.KindFloat32:
			return "0.0f"
		case ast.KindFloat64:
			return "0.0"
		case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16, ast.KindInt32, ast.KindUint32:
			return "0"
		case ast.KindInt64, ast.KindUint64:
			return "0L"
		case ast.KindString:
			return `""`
		case ast.KindBytes:
			return "(byte) 0"
		default:
			return "unknown"
		}
	default:
		return "unknown"
	}
}

func (tf *typeFormatter) formatScalar(v any) string {
	if list, ok := v.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, tf.formatScalar(item))
		}

		// FIXME: this is wrong, we can't just assume a list of strings.
		return strings.Join(items, ", ")
	}

	return fmt.Sprintf("%#v", v)
}

type CastPath struct {
	Class        string
	Value        string
	Path         string
	IsNilChecked bool
}

// formatCastValue identifies if the object to set is a generic one, so it needs
// to do a cast to the desired object to be able to set their values.
func (tf *typeFormatter) formatCastValue(fieldPath ast.Path) CastPath {
	refPkg := ""
	refType := ""
	for _, path := range fieldPath {
		if path.TypeHint != nil && path.TypeHint.Kind == ast.KindRef {
			refPkg = path.TypeHint.AsRef().ReferredPkg
			refType = path.TypeHint.AsRef().ReferredType
		}
	}

	if refType == "" {
		return CastPath{}
	}

	castedPath := fieldPath[0].Identifier
	isNilChecked := false
	genericFound := false

	for i, p := range fieldPath {
		if i > 0 && fieldPath[i-1].Type.IsAny() && i != len(fieldPath)-1 {
			isNilChecked = true
		}

		if !genericFound {
			if i > 0 {
				castedPath = fmt.Sprintf("%s.%s", castedPath, tools.LowerCamelCase(p.Identifier))
			}
			genericFound = p.Type.IsAny()
		}
	}

	return CastPath{
		Class:        fmt.Sprintf("%s.%s", tf.formatPackage(refPkg), refType),
		Value:        refType,
		Path:         castedPath,
		IsNilChecked: isNilChecked,
	}
}

func (tf *typeFormatter) shouldCastNilCheck(fieldPath ast.Path) CastPath {
	refPkg := ""
	refType := ""
	for _, path := range fieldPath {
		if path.TypeHint == nil && path.Type.Kind == ast.KindRef {
			refPkg = path.Type.AsRef().ReferredPkg
			refType = path.Type.AsRef().ReferredType
		}
	}

	if refType == "" {
		return CastPath{}
	}

	castedPath := fieldPath[0].Identifier
	isNilChecked := false
	genericFound := false

	for i, p := range fieldPath {
		if i > 0 && p.Type.IsRef() && !fieldPath[i-1].Type.IsRef() {
			refType = fieldPath[i-1].Identifier
			isNilChecked = true
		}

		if !genericFound {
			if i > 0 {
				castedPath = fmt.Sprintf("%s.%s", castedPath, tools.LowerCamelCase(p.Identifier))
			}
			genericFound = p.Type.IsAny()
		}
	}

	return CastPath{
		Class:        fmt.Sprintf("%s.%s", tf.formatPackage(refPkg), tools.UpperCamelCase(refType)),
		Value:        refType,
		Path:         castedPath,
		IsNilChecked: isNilChecked,
	}
}

func (tf *typeFormatter) formatFieldPath(fieldPath ast.Path) string {
	parts := make([]string, 0)
	for i, part := range fieldPath {
		output := tools.LowerCamelCase(part.Identifier)

		if i > 0 && fieldPath[i-1].Type.IsAny() {
			return output
		}

		parts = append(parts, output)
	}

	return strings.Join(parts, ".")
}

// formatAssignmentPath generates the pad to assign the value. When the value is a generic one (Object) like Custom or FieldConfig
// we should return until this pad to set the object to it.
func (tf *typeFormatter) formatAssignmentPath(fieldPath ast.Path) string {
	path := escapeVarName(tools.LowerCamelCase(fieldPath[0].Identifier))

	if len(fieldPath[1:]) == 1 && fieldPath[0].TypeHint != nil && fieldPath[0].TypeHint.Kind == ast.KindRef {
		return path
	}

	for i, p := range fieldPath[1:] {
		if p.Index != nil {
			path += "["
			if p.Index.Constant != nil {
				path += fmt.Sprintf("%#v", p.Index.Constant)
			} else {
				path += tools.LowerCamelCase(p.Index.Argument.Name)
			}
			path += "]"
			continue
		}

		if fieldPath[i].Type.IsAny() && i != len(fieldPath)-1 {
			return path
		}

		path = fmt.Sprintf("%s.%s", path, tools.LowerCamelCase(p.Identifier))
		if p.TypeHint != nil {
			return path
		}
	}

	return path
}

func (tf *typeFormatter) formatPackage(pkg string) string {
	if tf.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", tf.config.PackagePath, pkg)
	}

	return pkg
}

func (tf *typeFormatter) formatRefType(destinationType ast.Type, value any) string {
	if destinationType.IsRef() {
		referredObj, found := tf.context.LocateObject(destinationType.AsRef().ReferredPkg, destinationType.AsRef().ReferredType)
		if found && referredObj.Type.IsEnum() {
			return tf.formatEnumValue(referredObj, value)
		}
	}

	return fmt.Sprintf("%#v", value)
}

func (tf *typeFormatter) formatEnumValue(obj ast.Object, val any) string {
	member, _ := obj.Type.AsEnum().MemberForValue(val)

	return fmt.Sprintf("%s.%s", obj.Name, tools.UpperSnakeCase(member.Name))
}

func (tf *typeFormatter) objectNeedsCustomSerializer(obj ast.Object) bool {
	if !tf.config.generateBuilders || tf.config.SkipRuntime {
		return false
	}
	if obj.Type.HasHint(ast.HintDisjunctionOfScalars) {
		tf.packageMapper(fasterXMLPackageName, "databind.annotation.JsonSerialize")
		return true
	}

	return false
}

func (tf *typeFormatter) objectNeedsCustomDeserializer(obj ast.Object) bool {
	if !tf.config.generateBuilders || tf.config.SkipRuntime {
		return false
	}
	if objectNeedsCustomDeserialiser(tf.context, obj) {
		tf.packageMapper(fasterXMLPackageName, "databind.annotation.JsonDeserialize")
		return true
	}

	return false
}

func (tf *typeFormatter) fillNullableAnnotationPattern(t ast.Type) string {
	if t.Nullable {
		tf.packageMapper(fasterXMLPackageName, "annotation.JsonInclude")
		return javaNullableField
	}

	if t.IsArray() || t.IsMap() {
		tf.packageMapper(fasterXMLPackageName, "annotation.JsonSetter")
		tf.packageMapper(fasterXMLPackageName, "annotation.Nulls")
		return javaDefaultEmptyField
	}

	if t.IsAny() || t.IsStruct() || t.IsRef() {
		tf.packageMapper(fasterXMLPackageName, "annotation.JsonInclude")
		return javaEmptyField
	}

	return ""
}

func (tf *typeFormatter) formatGuardPath(fieldPath ast.Path) string {
	parts := make([]string, 0)
	var castedPath string

	for i := range fieldPath {
		output := fieldPath[i].Identifier
		if !fieldPath[i].Root {
			output = escapeVarName(tools.LowerCamelCase(output))
		}

		// don't generate type hints if:
		// * there isn't one defined
		// * the type isn't "any"
		// * as a trailing element in the path
		if !fieldPath[i].Type.IsAny() || fieldPath[i].TypeHint == nil || i == len(fieldPath)-1 {
			parts = append(parts, output)
			continue
		}

		castedPath = fmt.Sprintf("((%s) %s.%s).", tf.formatFieldType(*fieldPath[i].TypeHint), strings.Join(parts, "."), output)
		parts = nil
	}

	return castedPath + strings.Join(parts, ".")
}
