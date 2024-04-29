package compiler

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

var _ Pass = (*FilterSchemas)(nil)

// FilterSchemas filters a schema to only include the allowed objects and their
// dependencies.
type FilterSchemas struct {
	AllowedObjects []ObjectReference
}

func (pass *FilterSchemas) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	allowList := pass.buildAllowList(schemas, pass.AllowedObjects)

	return tools.Map(schemas, func(schema *ast.Schema) *ast.Schema {
		return pass.processSchema(schema, allowList)
	}), nil
}

func (pass *FilterSchemas) processSchema(schema *ast.Schema, allowList *orderedmap.Map[string, struct{}]) *ast.Schema {
	schema.Objects = schema.Objects.Filter(func(_ string, object ast.Object) bool {
		return allowList.Has(object.SelfRef.String())
	})

	return schema
}

// buildAllowList returns the set of objects that should be included in the
// processed schemas. This set is built by recursively exploring the
// "entrypoint objects" and any object they might reference, each of these
// references contributing to the allow list.
func (pass *FilterSchemas) buildAllowList(schemas ast.Schemas, entrypoints []ObjectReference) *orderedmap.Map[string, struct{}] {
	allowList := orderedmap.New[string, struct{}]()
	rootObjects := orderedmap.New[string, ast.Object]()

	for _, allowedObj := range entrypoints {
		obj, found := schemas.LocateObject(allowedObj.Package, allowedObj.Object)
		if !found {
			continue
		}

		rootObjects.Set(obj.SelfRef.String(), obj)
	}

	var exploreType func(typeDef ast.Type)
	exploreType = func(typeDef ast.Type) {
		if typeDef.IsRef() {
			ref := typeDef.Ref
			referredObj, found := schemas.LocateObject(ref.ReferredPkg, ref.ReferredType)
			if !found {
				return
			}

			rootObjects.Set(ref.String(), referredObj)
		}

		if typeDef.IsStruct() {
			for _, field := range typeDef.Struct.Fields {
				exploreType(field.Type)
			}
		}

		if typeDef.IsDisjunction() {
			for _, branch := range typeDef.Disjunction.Branches {
				exploreType(branch)
			}
		}

		if typeDef.IsIntersection() {
			for _, branch := range typeDef.Intersection.Branches {
				exploreType(branch)
			}
		}

		if typeDef.IsArray() {
			exploreType(typeDef.Array.ValueType)
		}

		if typeDef.IsMap() {
			exploreType(typeDef.Map.IndexType)
			exploreType(typeDef.Map.ValueType)
		}
	}

	for {
		if rootObjects.Len() == 0 {
			break
		}

		objects := rootObjects
		rootObjects = orderedmap.New[string, ast.Object]()

		objects.Iterate(func(key string, object ast.Object) {
			if allowList.Has(object.SelfRef.String()) {
				return
			}

			allowList.Set(key, struct{}{})

			exploreType(object.Type)
		})
	}

	return allowList
}
