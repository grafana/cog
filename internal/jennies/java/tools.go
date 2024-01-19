package java

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"strings"
)

func formatFieldPath(fieldPath ast.Path) string {
	parts := tools.Map(fieldPath, func(part ast.PathItem) string {
		return tools.LowerCamelCase(part.Identifier)
	})

	return strings.Join(parts, ".")
}
