package transforms

import (
	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

var _ Transform = (*FilterSchemas)(nil)

// FilterSchemas filters a schema to only include the allowed objects and their
// dependencies.
type FilterSchemas struct {
	AllowedObjects []ObjectReference
}

func (pass *FilterSchemas) Process(schemas ir.Schemas) (ir.Schemas, error) {
	allowList := pass.buildAllowList(schemas, pass.AllowedObjects)

	return tools.Map(schemas, func(schema *ir.Schema) *ir.Schema {
		return pass.processSchema(schema, allowList)
	}), nil
}

func (pass *FilterSchemas) processSchema(schema *ir.Schema, allowList *orderedmap.Map[string, struct{}]) *ir.Schema {
	schema.Objects = schema.Objects.Filter(func(_ string, object ir.Object) bool {
		return allowList.Has(object.SelfRef.String())
	})

	return schema
}

// buildAllowList returns the set of objects that should be included in the
// processed schemas. This set is built by recursively exploring the
// "entrypoint objects" and any object they might reference, each of these
// references contributing to the allow list.
func (pass *FilterSchemas) buildAllowList(schemas ir.Schemas, entrypoints []ObjectReference) *orderedmap.Map[string, struct{}] {
	allowList := orderedmap.New[string, struct{}]()
	rootObjects := orderedmap.New[string, ir.Object]()

	for _, allowedObj := range entrypoints {
		obj, found := schemas.GetObject(allowedObj.Package, allowedObj.Object)
		if !found {
			continue
		}

		rootObjects.Set(obj.SelfRef.String(), obj)
	}

	visitor := &Visitor{
		OnRef: func(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
			referredObj, found := schemas.GetObject(def.Ref.ReferredPkg, def.Ref.ReferredType)
			if !found {
				return def, nil
			}

			rootObjects.Set(def.Ref.String(), referredObj)

			return def, nil
		},
	}

	for {
		if rootObjects.Len() == 0 {
			break
		}

		objects := rootObjects
		rootObjects = orderedmap.New[string, ir.Object]()

		objects.Iterate(func(key string, object ir.Object) {
			if allowList.Has(object.SelfRef.String()) {
				return
			}

			allowList.Set(key, struct{}{})

			schema, found := schemas.Get(object.SelfRef.ReferredPkg)
			if !found {
				return
			}

			_, _ = visitor.VisitType(schema, object.Type)
		})
	}

	return allowList
}
