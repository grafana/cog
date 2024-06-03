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
		if value.Type.Kind == ast.KindRef {
			obj, toRemove := r.processObject(schema, value)
			schema.Objects.Set(key, obj)
			listToRemove[toRemove] = obj
		}
	})

	schema.Objects.Iterate(func(key string, value ast.Object) {
		if value.Type.Kind == ast.KindStruct {
			schema.Objects.Set(key, r.processStruct(value, listToRemove))
		}
	})

	// fmt.Println(listToRemove)
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
	newObject.Type.Hints[ast.HintImplementsVariant] = object.Type.ImplementedVariant()

	return newObject, locatedObject.Name
}

func (r RemoveIntersections) processStruct(object ast.Object, listToRemove map[string]ast.Object) ast.Object {
	str := object.Type.AsStruct()
	for i, field := range str.Fields {
		// fmt.Println(field.Name)
		if obj, ok := listToRemove[field.Name]; ok {
			object.Type.AsStruct().Fields[i].Name = obj.Name
			object.Type.AsStruct().Fields[i].Type = obj.Type
		}
	}

	return object
}
