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

type CastPath struct {
	Class string
	Value string
	Path  string
}

func formatScalar(val any) any {
	newVal := fmt.Sprintf("%#v", val)
	if len(strings.Split(newVal, ".")) > 1 {
		return val
	}
	return newVal
}

// formatAssignmentPath generates the pad to assign the value. When the value is a generic one (Object) like Custom or FieldConfig
// we should return until this pad to set the object to it.
func formatAssignmentPath(fieldPath ast.Path) string {
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
