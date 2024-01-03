package java

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
)

type typeFormatter struct {
	packageMapper func(pkg string, class string) string
	context       common.Context
}

func defaultTypeFormatter(ctx common.Context, packageMapper func(pkg string, class string) string) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
		context:       ctx,
	}
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
	}

	return "unknown"
}

func (tf *typeFormatter) formatReference(def ast.RefType) string {
	tf.packageMapper(def.ReferredPkg, def.ReferredType)
	object, _ := tf.context.LocateObject(def.ReferredPkg, def.ReferredType)
	if object.Type.Kind == ast.KindScalar {
		return formatScalarType(object.Type.AsScalar())
	}
	return def.ReferredType
}

func (tf *typeFormatter) formatArray(def ast.ArrayType) string {
	return fmt.Sprintf("%s[]", tf.formatFieldType(def.ValueType))
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
