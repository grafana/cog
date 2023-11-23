package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

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

	var obj ast.Object
	var keyObject ast.Object

	for i, object := range schema.Objects {
		if object.Name != testDataSimulatorQueryObject {
			continue
		}

		obj, keyObject = t.processObject(object)
		schema.Objects[i] = obj
	}

	if keyObject.Name != "" {
		schema.Objects = append(schema.Objects, keyObject)
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
