package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
)

type typeFormatter struct {
	config        Config
	packageMapper func(pkg string, class string) string
	context       common.Context
}

func createFormatter(ctx common.Context, config Config) *typeFormatter {
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
		ref := def.ValueType.AsRef()
		tf.packageMapper(ref.ReferredPkg, ref.ReferredType)
		mapType = ref.ReferredType
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
			return fmt.Sprintf("new %s.Builder().build()", tf.formatPackage(refDef))
		}
		return fmt.Sprintf("new %s()", tf.formatPackage(refDef))
	case ast.KindStruct:
		return "new Object()"
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
	Class string
	Value string
	Path  string
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
	for _, p := range fieldPath[1 : len(fieldPath)-1] {
		castedPath = fmt.Sprintf("%s.%s", castedPath, tools.LowerCamelCase(p.Identifier))
	}

	return CastPath{
		Class: fmt.Sprintf("%s.%s", tf.formatPackage(refPkg), refType),
		Value: refType,
		Path:  castedPath,
	}
}

// formatAssignmentPath generates the pad to assign the value. When the value is a generic one (Object) like Custom or FieldConfig
// we should return until this pad to set the object to it.
func (tf *typeFormatter) formatAssignmentPath(fieldPath ast.Path) string {
	path := escapeVarName(tools.LowerCamelCase(fieldPath[0].Identifier))

	if len(fieldPath[1:]) == 1 && fieldPath[0].TypeHint != nil && fieldPath[0].TypeHint.Kind == ast.KindRef {
		return path
	}

	for _, p := range fieldPath[1:] {
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
