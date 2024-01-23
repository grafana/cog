package java

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func formatFieldPath(fieldPath ast.Path) []string {
	parts := tools.Map(fieldPath, func(part ast.PathItem) string {
		return tools.UpperCamelCase(part.Identifier)
	})

	return parts
}

func formatAssignmentPath(fieldPath ast.Path, method ast.AssignmentMethod) string {
	path := fieldPath[0].Identifier
	for i, p := range fieldPath[1:] {
		identifier := tools.UpperCamelCase(p.Identifier)
		if i == len(fieldPath[1:])-1 && method != ast.AppendAssignment {
			path = fmt.Sprintf("%s.set%s", path, identifier)
			continue
		}
		path = fmt.Sprintf("%s.get%s()", path, identifier)
	}

	return path
}
