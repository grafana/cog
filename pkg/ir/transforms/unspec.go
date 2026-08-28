package transforms

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*Unspec)(nil)

// Unspec removes the Kubernetes-style envelope added by kindsys.
//
// Objects named "spec" will be renamed, using the package as new name.
type Unspec struct {
}

func (pass *Unspec) Process(schemas ir.Schemas) (ir.Schemas, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *Unspec) processSchema(schema *ir.Schema) *ir.Schema {
	schema.Objects = schema.Objects.Filter(func(_ string, object ir.Object) bool {
		return !strings.EqualFold(object.Name, "metadata")
	})

	originalObjects := schema.Objects
	schema.Objects = orderedmap.New[string, ir.Object]()

	originalObjects.Iterate(func(name string, object ir.Object) {
		if strings.EqualFold(object.Name, "spec") && object.Type.IsStruct() {
			object.Name = schema.Package
			if schema.Metadata.Identifier != "" {
				object.Name = schema.Metadata.Identifier
			}

			object.SelfRef.ReferredType = object.Name
			object.AddToPassesTrail(fmt.Sprintf("Unspec[%s → %s]", name, object.Name))
		}

		schema.AddObject(object)
	})

	return schema
}
