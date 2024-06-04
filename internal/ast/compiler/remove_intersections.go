package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

type RemoveIntersections struct {
	listToRemove map[string]ast.Object
}

func (r RemoveIntersections) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	r.listToRemove = make(map[string]ast.Object)
	visitor := Visitor{
		OnSchema: r.processSchema,
		OnObject: r.processObject,
		OnStruct: r.processStruct,
	}

	return visitor.VisitSchemas(schemas)
}

func (r RemoveIntersections) processSchema(v *Visitor, schema *ast.Schema) (*ast.Schema, error) {
	var foundErr error
	schema.Objects.Iterate(func(key string, value ast.Object) {
		if value.Type.IsRef() {
			obj, err := v.VisitObject(schema, value)
			if err != nil {
				foundErr = err
			}
			schema.Objects.Set(key, obj)
		}
	})

	if foundErr != nil {
		return nil, foundErr
	}

	schema.Objects.Iterate(func(key string, value ast.Object) {
		if value.Type.IsStruct() {
			if _, err := v.VisitStruct(schema, value.Type); err != nil {
				foundErr = err
			}
		}
	})

	if foundErr != nil {
		return nil, foundErr
	}

	for toRemove, _ := range r.listToRemove {
		schema.Objects.Remove(toRemove)
	}

	return schema, nil
}

func (r RemoveIntersections) processObject(_ *Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	ref := object.Type.AsRef()
	locatedObject, ok := schema.LocateObject(ref.ReferredType)
	if !ok {
		return object, nil
	}

	newObject := object
	newObject.Type = ast.NewStruct(locatedObject.Type.AsStruct().Fields...)
	if object.Type.ImplementsVariant() {
		newObject.Type.Hints[ast.HintImplementsVariant] = object.Type.ImplementedVariant()
	}

	r.listToRemove[locatedObject.Name] = object
	return newObject, nil
}

func (r RemoveIntersections) processStruct(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	str := def.AsStruct()
	for i, field := range str.Fields {
		// TODO: Add Map/List checks if necessary
		if field.Type.IsRef() {
			if obj, ok := r.listToRemove[field.Type.AsRef().ReferredType]; ok {
				def.AsStruct().Fields[i] = ast.NewStructField(obj.Name, ast.NewRef(obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType), ast.Comments(obj.Comments))
			}
		}
	}

	return def, nil
}
