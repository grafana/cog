package testutils

import (
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/pkg/ast"
)

func ObjectsMap(objects ...ast.Object) *orderedmap.Map[string, ast.Object] {
	ordered := orderedmap.New[string, ast.Object]()
	for _, object := range objects {
		ordered.Set(object.Name, object)
	}
	return ordered
}
