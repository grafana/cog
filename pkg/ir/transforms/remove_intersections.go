package transforms

import (
	"maps"

	"github.com/grafana/cog/pkg/ir"
)

type RemoveIntersections struct {
	objectsToRemove map[string]ir.Object
	arraysToFix     map[string]ir.Object
}

func (r RemoveIntersections) Process(schemas ir.Schemas) (ir.Schemas, error) {
	r.objectsToRemove = make(map[string]ir.Object)
	r.arraysToFix = make(map[string]ir.Object)
	visitor := Visitor{
		OnSchema: r.processSchema,
		OnObject: r.processObject,
		OnStruct: r.processStruct,
	}

	return visitor.VisitSchemas(schemas)
}

func (r RemoveIntersections) processSchema(v *Visitor, schema *ir.Schema) (*ir.Schema, error) {
	var foundErr error
	schema.Objects.Iterate(func(key string, value ir.Object) {
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

	schema.Objects.Iterate(func(key string, value ir.Object) {
		if value.Type.IsStruct() {
			if _, err := v.VisitStruct(schema, value.Type); err != nil {
				foundErr = err
			}
		}
	})

	if foundErr != nil {
		return nil, foundErr
	}

	for toRemove := range r.objectsToRemove {
		schema.Objects.Remove(toRemove)
	}

	return schema, nil
}

func (r RemoveIntersections) processObject(_ *Visitor, schema *ir.Schema, object ir.Object) (ir.Object, error) {
	ref := object.Type.AsRef()
	locatedObject, ok := schema.GetObject(ref.ReferredType)
	if !ok {
		return object, nil
	}

	if locatedObject.Type.IsStruct() {
		newObject := object
		newObject.Type = ir.NewStruct(locatedObject.Type.AsStruct().Fields...)
		if object.Type.ImplementsVariant() {
			newObject.Type.Hints[ir.HintImplementsVariant] = object.Type.ImplementedVariant()
		}
		maps.Copy(newObject.Type.Hints, locatedObject.Type.Hints)

		r.objectsToRemove[locatedObject.Name] = object
		return newObject, nil
	}

	if locatedObject.Type.IsArray() {
		r.objectsToRemove[object.Name] = object
		r.arraysToFix[object.Name] = locatedObject
	}

	// TODO: Check if a reference extends from a Map if necessary

	return object, nil
}

func (r RemoveIntersections) processStruct(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	str := def.AsStruct()
	for i, field := range str.Fields {
		if field.Type.IsRef() {
			if obj, ok := r.objectsToRemove[field.Type.AsRef().ReferredType]; ok {
				def.AsStruct().Fields[i] = ir.NewStructField(field.Name, ir.NewRef(obj.SelfRef.ReferredPkg, obj.SelfRef.ReferredType), ir.Comments(obj.Comments))
			}
			if obj, ok := r.arraysToFix[field.Type.AsRef().ReferredType]; ok {
				def.AsStruct().Fields[i] = ir.NewStructField(field.Name, ir.NewArray(obj.Type.AsArray().ValueType), ir.Comments(obj.Comments))
			}

			maps.Copy(field.Type.Hints, def.AsStruct().Fields[i].Type.Hints)
		}
	}

	return def, nil
}
