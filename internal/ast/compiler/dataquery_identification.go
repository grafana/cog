package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DataqueryIdentification)(nil)

type DataqueryIdentification struct {
}

func (pass *DataqueryIdentification) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	commonDataquery, found := ast.Schemas(schemas).LocateObject("common", "DataQuery")
	if !found {
		return schemas, nil
	}

	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		newSchemas = append(newSchemas, pass.processSchema(schema, commonDataquery))
	}

	return newSchemas, nil
}

func (pass *DataqueryIdentification) processSchema(schema *ast.Schema, commonDataquery ast.Object) *ast.Schema {
	newSchema := schema.DeepCopy()
	for i, object := range schema.Objects {
		newSchema.Objects[i] = pass.processObject(object, commonDataquery)
	}

	return &newSchema
}

func (pass *DataqueryIdentification) processObject(object ast.Object, commonDataquery ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	typeDef := object.Type

	// this object is already identified as a variant: nothing to do.
	if typeDef.ImplementsVariant() {
		return object
	}

	if pass.structsIntersect(typeDef, commonDataquery.Type) {
		object.Type.Hints[ast.HintImplementsVariant] = string(ast.SchemaVariantDataQuery)
		object.Type.AddCompilerPassTrail("DataqueryIdentification")
	}

	return object
}

func (pass *DataqueryIdentification) structsIntersect(def ast.Type, base ast.Type) bool {
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
