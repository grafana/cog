package terraform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/tools"
)

func formatPackageName(pkg string) string {
	splitPath := strings.Split(pkg, "/")
	if len(splitPath) > 1 {
		pkg = splitPath[len(splitPath)-1]
	}
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return strings.ToLower(rgx.ReplaceAllString(pkg, ""))
}

func formatObjectName(name string) string {
	return tools.UpperCamelCase(name)
}

func formatFieldName(name string) string {
	return tools.UpperCamelCase(name)
}

func formatScalar(val any) string {
	if val == nil {
		return "nil"
	}

	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatScalar(item))
		}

		// FIXME: this is wrong, we can't just assume a list of strings.
		return fmt.Sprintf("[]string{%s}", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}
