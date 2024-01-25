package java

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"strings"
)

type typeFormatter struct {
	packageMapper func(pkg string, class string) string
	context       common.Context
}

func createFormatter(ctx common.Context) *typeFormatter {
	return &typeFormatter{context: ctx}
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

func (tf *typeFormatter) formatBuilderArgs(def ast.Type) string {
	value := tf.formatFieldType(def)
	if _, ok := tf.context.ResolveToComposableSlot(def); ok || tf.context.ResolveToBuilder(def) {
		value = fmt.Sprintf("Builder<%s>", value)
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
	mapType := "unknown"
	switch def.ValueType.Kind {
	case ast.KindRef:
		ref := def.ValueType.AsRef()
		tf.packageMapper(ref.ReferredPkg, ref.ReferredType)
		tf.packageMapper("java.util", "Map")
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
		return fmt.Sprintf("new %s()", refDef)
	case ast.KindStruct:
		return "new Object()"
	default:
		return "unknown"
	}
}

func formatScalarType(def ast.ScalarType) string {
	scalarType := "unknown"

	switch def.ScalarKind {
	case ast.KindString:
		scalarType = "String"
	case ast.KindBytes, ast.KindInt8, ast.KindUint8:
		scalarType = "Byte"
	case ast.KindInt16, ast.KindUint16:
		scalarType = "Short"
	case ast.KindInt32, ast.KindUint32:
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

func formatScalar(val any) any {
	newVal := fmt.Sprintf("%#v", val)
	if len(strings.Split(newVal, ".")) > 1 {
		return val
	}
	return newVal
}
