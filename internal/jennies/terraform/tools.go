package terraform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}
func formatTerraformType(t ast.Type) string {
	if t.IsScalar() {
		tt := t.AsScalar()
		var scalarType string

		switch tt.ScalarKind {
		case ast.KindString, ast.KindBytes:
			scalarType = "types.String"
		case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
			scalarType = "types.Int64"
		case ast.KindInt32, ast.KindUint32:
			scalarType = "types.Int64"
		case ast.KindInt64, ast.KindUint64:
			scalarType = "types.Int64"
		case ast.KindFloat32:
			scalarType = "types.Float64"
		case ast.KindFloat64:
			scalarType = "types.Float64"
		case ast.KindBool:
			scalarType = "types.Bool"
		case ast.KindAny:
			scalarType = "types.Object"
		default:
			scalarType = "types.Any"
		}
		return scalarType
	}
	return "types.Any"
}

func formatGolangType(t ast.Type) string {
	if t.IsScalar() {

		tt := t.AsScalar()
		var scalarType string

		switch tt.ScalarKind {
		case ast.KindString, ast.KindBytes:
			scalarType = "string"
		case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
			scalarType = "int64"
		case ast.KindInt32, ast.KindUint32:
			scalarType = "int64"
		case ast.KindInt64, ast.KindUint64:
			scalarType = "int64"
		case ast.KindFloat32:
			scalarType = "float64"
		case ast.KindFloat64:
			scalarType = "float64"
		case ast.KindBool:
			scalarType = "bool"
		case ast.KindAny:
			scalarType = "any"
		default:
			scalarType = "any"
		}
		if t.Nullable {
			scalarType = "*" + scalarType
		}
		return scalarType
	}
	return "types.Any"
}
func formatJSONField(f ast.StructField) string {
	if f.Type.Nullable {
		return fmt.Sprintf("%s,omitempty", f.Name)
	}
	return f.Name
}

func formatTypeValue(t ast.Type) string {
	if t.IsScalar() {

		tt := t.AsScalar()
		var scalarType string

		switch tt.ScalarKind {
		case ast.KindString, ast.KindBytes:
			scalarType = "String"
		case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
			scalarType = "Int64"
		case ast.KindInt32, ast.KindUint32:
			scalarType = "Int64"
		case ast.KindInt64, ast.KindUint64:
			scalarType = "Int64"
		case ast.KindFloat32:
			scalarType = "Float64"
		case ast.KindFloat64:
			scalarType = "Float64"
		case ast.KindBool:
			scalarType = "Bool"
		case ast.KindAny:
			scalarType = "Any"
		default:
			scalarType = "Any"
		}
		if t.Nullable {
			return fmt.Sprintf("Value%sPointer()", scalarType)
		}
		return fmt.Sprintf("Value%s()", scalarType)
	}
	return "ValueAny()"
}

func formatTypeValueNoPointers(dflt any, t ast.Type) string {
	if t.IsScalar() {

		switch t.AsScalar().ScalarKind {
		case ast.KindString, ast.KindBytes:
			return fmt.Sprintf("StringValue(`%s`)", dflt)
		case ast.KindInt8, ast.KindUint8, ast.KindInt16, ast.KindUint16:
			return fmt.Sprintf("StringValue(%d)", dflt)
		case ast.KindInt32, ast.KindUint32:
			return fmt.Sprintf("Int64Value(%d)", dflt)
		case ast.KindInt64, ast.KindUint64:
			return fmt.Sprintf("Int64Value(%d)", dflt)
		case ast.KindFloat32:
			return fmt.Sprintf("Float64Value(%f)", dflt)
		case ast.KindFloat64:
			return fmt.Sprintf("Float64Value(%f)", dflt)
		case ast.KindBool:
			return fmt.Sprintf("BoolValue(%t)", dflt)
		case ast.KindAny:
			return fmt.Sprintf("AnyValue(%v)", dflt)
		default:
			return fmt.Sprintf("AnyValue(%v)", dflt)
		}
	}
	return fmt.Sprintf("AnyValue(%v)", dflt)
}
