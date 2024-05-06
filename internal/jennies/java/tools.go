package java

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

func formatFieldPath(fieldPath ast.Path) []string {
	parts := tools.Map(fieldPath, func(fieldPath ast.PathItem) string {
		return tools.LowerCamelCase(fieldPath.Identifier)
	})
	return parts
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
		castedPath = fmt.Sprintf("%s.%s", castedPath, tools.LowerCamelCase(p.Identifier))
	}

	return CastPath{
		Class: fmt.Sprintf("%s.%s", refPkg, refType),
		Value: refType,
		Path:  castedPath,
	}
}

func formatScalar(val any) any {
	newVal := fmt.Sprintf("%#v", val)
	if len(strings.Split(newVal, ".")) > 1 {
		return val
	}
	return newVal
}

func formatAssignmentPath(fieldPath ast.Path) string {
	path := tools.LowerCamelCase(fieldPath[0].Identifier)

	if len(fieldPath[1:]) == 1 && fieldPath[0].TypeHint != nil && fieldPath[0].TypeHint.Kind == ast.KindRef {
		return tools.LowerCamelCase(path)
	}

	for i, p := range fieldPath[1:] {
		identifier := tools.LowerCamelCase(p.Identifier)
		if i == 0 && p.TypeHint != nil && p.TypeHint.Kind == ast.KindRef {
			return tools.LowerCamelCase(identifier)
		}

		path = fmt.Sprintf("%s.%s", path, identifier)
	}

	return path
}
