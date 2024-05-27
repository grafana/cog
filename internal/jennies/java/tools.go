package java

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}

func formatFieldPath(fieldPath ast.Path) string {
	parts := tools.Map(fieldPath, func(fieldPath ast.PathItem) string {
		return tools.LowerCamelCase(fieldPath.Identifier)
	})
	return strings.Join(parts, ".")
}

func formatScalar(val any) any {
	newVal := fmt.Sprintf("%#v", val)
	if len(strings.Split(newVal, ".")) > 1 {
		return val
	}
	return newVal
}
