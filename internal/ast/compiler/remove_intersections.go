package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

type RemoveIntersections struct {
}

func (r RemoveIntersections) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0)
	for _, schema := range schemas {
		if sch, ok := r.processSchema(schema); ok {
			newSchemas = append(newSchemas, sch)
		}
	}

	return newSchemas, nil
}

func (r RemoveIntersections) processSchema(schema *ast.Schema) (*ast.Schema, bool) {
	listToRemove := make(map[string]ast.Object)
	schema.Objects.Iterate(func(key string, value ast.Object) {
		if value.Type.IsRef() {
			obj, toRemove := r.processObject(schema, value)
			schema.Objects.Set(key, obj)
			listToRemove[toRemove] = obj
		}
	})

	schema.Objects.Iterate(func(key string, value ast.Object) {
		if value.Type.IsStruct() {
			schema.Objects.Set(key, r.processStruct(value, listToRemove))
		}
	})

	for toRemove, _ := range listToRemove {
		schema.Objects.Remove(toRemove)
	}

	return schema, true
}

func (r RemoveIntersections) processObject(schema *ast.Schema, object ast.Object) (ast.Object, string) {
	ref := object.Type.AsRef()
	locatedObject, ok := schema.LocateObject(ref.ReferredType)
	if !ok {
		return object, ""
	}

	newObject := object
	newObject.Type = ast.NewStruct(locatedObject.Type.AsStruct().Fields...)
	if object.Type.ImplementsVariant() {
		newObject.Type.Hints[ast.HintImplementsVariant] = object.Type.ImplementedVariant()
	}

	return newObject, locatedObject.Name
}

func (r RemoveIntersections) processStruct(object ast.Object, listToRemove map[string]ast.Object) ast.Object {
	str := object.Type.AsStruct()
	for i, field := range str.Fields {
		// TODO: Add Map/List checks if necessary
		if field.Type.IsRef() {
			if obj, ok := listToRemove[field.Type.AsRef().ReferredType]; ok {
				object.Type.AsStruct().Fields[i] = ast.NewStructField(obj.Name, ast.NewRef(obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType), ast.Comments(obj.Comments))
			}
		}
	}

	return object
}
