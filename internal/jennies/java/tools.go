package java

import (
	"fmt"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func formatFieldPath(fieldPath ast.Path) []string {
	parts := tools.Map(fieldPath, func(fieldPath ast.PathItem) string {
		return tools.UpperCamelCase(fieldPath.Identifier)
	})
	return parts
}

func formatAssignmentPath(fieldPath ast.Path, method ast.AssignmentMethod) string {
	path := tools.LowerCamelCase(fieldPath[0].Identifier)

	if len(fieldPath[1:]) == 1 && fieldPath[0].TypeHint != nil && fieldPath[0].TypeHint.Kind == ast.KindRef {
		return tools.LowerCamelCase(path)
	}

	for i, p := range fieldPath[1:] {
		identifier := tools.UpperCamelCase(p.Identifier)
		if i == 0 && p.TypeHint != nil && p.TypeHint.Kind == ast.KindRef {
			return tools.LowerCamelCase(identifier)
		}

		if p.TypeHint != nil && p.TypeHint.Kind == ast.KindRef {
			path = fmt.Sprintf("%s.set%s", path, identifier)
			break
		}

		if i == len(fieldPath[1:])-1 && method != ast.AppendAssignment {
			path = fmt.Sprintf("%s.set%s", path, identifier)
			continue
		}

		path = fmt.Sprintf("%s.get%s()", path, identifier)
	}

	return path
}

type CastPath struct {
	Class string
	Value string
	Path  string
}

func formatCastValue(fieldPath ast.Path) CastPath {
	refPkg := ""
	refType := ""
	for _, path := range fieldPath {
		if path.TypeHint != nil && path.TypeHint.Kind == ast.KindRef {
			refPkg = path.TypeHint.AsRef().ReferredPkg
			refType = path.TypeHint.AsRef().ReferredType
		}
	}

	if refType == "" {
		return CastPath{}
	}

	castedPath := fieldPath[0].Identifier
	for _, p := range fieldPath[1 : len(fieldPath)-1] {
		castedPath = fmt.Sprintf("%s.get%s()", castedPath, tools.UpperCamelCase(p.Identifier))
	}

	return CastPath{
		Class: fmt.Sprintf("%s.%s", refPkg, refType),
		Value: refType,
		Path:  castedPath,
	}
}
