package transforms

import (
	"github.com/grafana/cog/pkg/ir"
)

var _ Transform = (*DataqueryIdentification)(nil)

type DataqueryIdentification struct {
}

func (pass *DataqueryIdentification) Process(schemas ir.Schemas) (ir.Schemas, error) {
	commonDataquery, found := schemas.GetObject("common", "DataQuery")
	if !found {
		return schemas, nil
	}

	newSchemas := make(ir.Schemas, 0, len(schemas))

	for _, schema := range schemas {
		newSchemas = append(newSchemas, pass.processSchema(schema, commonDataquery))
	}

	return newSchemas, nil
}

func (pass *DataqueryIdentification) processSchema(schema *ir.Schema, commonDataquery ir.Object) *ir.Schema {
	var variantObjects []string
	schema.Objects = schema.Objects.Map(func(_ string, object ir.Object) ir.Object {
		if object.SelfRef.String() == commonDataquery.SelfRef.String() {
			return object
		}

		obj, implementsVariant := pass.processObject(object, commonDataquery)

		if implementsVariant {
			variantObjects = append(variantObjects, obj.Name)
		}

		return obj
	})

	if len(variantObjects) != 0 {
		schema.Metadata.Kind = ir.SchemaKindComposable
		schema.Metadata.Variant = ir.SchemaVariantDataQuery
	}

	if schema.EntryPoint == "" && len(variantObjects) == 1 {
		schema.EntryPoint = variantObjects[0]
		schema.EntryPointType = schema.Objects.Get(variantObjects[0]).SelfRef.AsType()
	}

	return schema
}

func (pass *DataqueryIdentification) processObject(object ir.Object, commonDataquery ir.Object) (ir.Object, bool) {
	if !object.Type.IsStruct() {
		return object, false
	}

	typeDef := object.Type

	// this object is already identified as a variant: nothing to do.
	if typeDef.ImplementsVariant() {
		return object, true
	}

	if !pass.structsIntersect(typeDef, commonDataquery.Type) {
		return object, false
	}

	object.Type.Hints[ir.HintImplementsVariant] = string(ir.SchemaVariantDataQuery)
	object.AddToPassesTrail("DataqueryIdentification[hint.ImplementsVariant=VariantDataQuery]")

	return object, true
}

func (pass *DataqueryIdentification) structsIntersect(def ir.Type, base ir.Type) bool {
	structDef := def.AsStruct()

	for _, baseField := range base.AsStruct().Fields {
		// ginormous assumption here: if we find fields with the same name, then we assume their types
		// to be identical too.
		if _, found := structDef.FieldByName(baseField.Name); !found {
			return false
		}
	}

	return true
}
