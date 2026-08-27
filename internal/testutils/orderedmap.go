package testutils

import (
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/pkg/ir"
)

func ObjectsMap(objects ...ir.Object) *orderedmap.Map[string, ir.Object] {
	ordered := orderedmap.New[string, ir.Object]()
	for _, object := range objects {
		ordered.Set(object.Name, object)
	}
	return ordered
}
