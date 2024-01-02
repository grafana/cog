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

func defaultTypeFormatter(packageMapper func(pkg string, class string) string) *typeFormatter {
	return &typeFormatter{
		packageMapper: packageMapper,
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
	case ast.KindDisjunction:
		return tf.formatDisjunction(def.AsDisjunction())
	}

	return "unknown"
}

func (tf *typeFormatter) formatReference(def ast.RefType) string {
	tf.packageMapper(def.ReferredPkg, def.ReferredType)
	return def.ReferredType
}

func (tf *typeFormatter) formatArray(def ast.ArrayType) string {
	return fmt.Sprintf("%s[]", tf.formatFieldType(def.ValueType))
}

func (tf *typeFormatter) formatDisjunction(def ast.DisjunctionType) string {
	return ""
}

func formatScalarType(def ast.ScalarType) string {
	scalarType := "unknown"

	switch def.ScalarKind {
	case ast.KindString:
		scalarType = "String"
	case ast.KindBytes, ast.KindInt8:
		scalarType = "Byte"
	case ast.KindInt16:
		scalarType = "Short"
	case ast.KindInt32:
		scalarType = "Integer"
	case ast.KindInt64:
		scalarType = "Long"
	case ast.KindFloat32:
		scalarType = "Float"
	case ast.KindFloat64:
		scalarType = "Double"
	case ast.KindBool:
		scalarType = "Boolean"
	}

	return scalarType
}
