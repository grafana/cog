package testutils

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/orderedmap"
)

func ObjectsMap(objects ...ast.Object) *orderedmap.Map[string, ast.Object] {
	ordered := orderedmap.New[string, ast.Object]()
	for _, object := range objects {
		ordered.Set(object.Name, object)
	}
	return ordered
}
