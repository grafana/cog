package terraform

import (
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}
func formatTerraformType(t ast.Type) string {
	if t.IsScalar() {
		tt := t.AsScalar()
		scalarType := "unknown"

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
		}
		return scalarType
	}
	return string(t.Kind)
}

func formatModelName(pkg, name string) string {
	pkgCamel := tools.UpperCamelCase(pkg)
	nameCamel := tools.UpperCamelCase(name)
	if strings.HasPrefix(nameCamel, pkgCamel) {
		return nameCamel
	}
	return pkgCamel + nameCamel
}
