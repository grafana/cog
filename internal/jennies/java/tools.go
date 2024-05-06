package java

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func formatFieldPath(fieldPath ast.Path) []string {
	parts := tools.Map(fieldPath, func(fieldPath ast.PathItem) string {
		return tools.UpperCamelCase(fieldPath.Identifier)
	})
	return parts
}
