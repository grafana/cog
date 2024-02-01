package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

// TestData creates a new Key struct and builder for anonymous SimulationQuery.key field, and it changes the builder signature.
//
// Original signature output (Go):
//
//	```
//
//	func (builder *SimulationQueryBuilder) Key(key struct {
//		Type string  `json:"type"`
//		Tick float64 `json:"tick"`
//		Uid  *string `json:"uid,omitempty"`
//	}) *SimulationQueryBuilder {
//
//		builder.internal.Key = key
//
//		return builder
//	}
//
// ```
//
// New signature (Go):
//
//	```
//
//	func (builder *SimulationQueryBuilder) Key(key cog.Builder[Key]) *SimulationQueryBuilder {
//		keyResource, err := key.Build()
//		if err != nil {
//			builder.errors["key"] = err.(cog.BuildErrors)
//			return builder
//		}
//		builder.internal.Key = keyResource
//
//		return builder
//	}
//
//	```
type TestData struct {
}

func (t TestData) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, len(schemas))

	for i, schema := range schemas {
		newSchemas[i] = t.processSchema(schema)
	}

	return newSchemas, nil
}

func (t TestData) processSchema(schema *ast.Schema) *ast.Schema {
	if schema.Package != testDataPackage {
		return schema
	}

	var processed ast.Object
	var keyObject ast.Object

	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		if object.Name != testDataSimulatorQueryObject {
			return object
		}

		processed, keyObject = t.processObject(object)
		return processed
	})

	if keyObject.Name != "" {
		schema.AddObject(keyObject)
	}

	return schema
}

func (t TestData) processObject(object ast.Object) (ast.Object, ast.Object) {
	var keyObject ast.Object

	for i, field := range object.Type.AsStruct().Fields {
		if field.Name != testDataSimulatorQueryKeyField {
			continue
		}

		if field.Type.Kind != ast.KindStruct {
			continue
		}

		keyObject = ast.NewObject(object.SelfRef.ReferredPkg, testDataKeyObject, field.Type)
		object.Type.AsStruct().Fields[i].Type = ast.NewRef(object.SelfRef.ReferredPkg, testDataKeyObject)
	}

	return object, keyObject
}
