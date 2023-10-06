package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*NotRequiredFieldAsNullableType)(nil)

type NotRequiredFieldAsNullableType struct {
}

func (pass *NotRequiredFieldAsNullableType) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	newSchemas := make([]*ast.Schema, 0, len(schemas))

	for _, schema := range schemas {
		newSchemas = append(newSchemas, pass.processSchema(schema))
	}

	return newSchemas, nil
}

func (pass *NotRequiredFieldAsNullableType) processSchema(schema *ast.Schema) *ast.Schema {
	processedObjects := make([]ast.Object, 0, len(schema.Objects))
	for _, object := range schema.Objects {
		processedObjects = append(processedObjects, pass.processObject(object))
	}

	newSchema := schema.DeepCopy()
	newSchema.Objects = processedObjects

	return &newSchema
}

func (pass *NotRequiredFieldAsNullableType) processObject(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	newObject := object
	newObject.Type = pass.processType(object.Type)

	return newObject
}

func (pass *NotRequiredFieldAsNullableType) processType(def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(def.AsArray())
	}

	if def.Kind == ast.KindMap {
		return pass.processMap(def.AsMap())
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(def.AsStruct())
	}

	return def
}

func (pass *NotRequiredFieldAsNullableType) processArray(def ast.ArrayType) ast.Type {
	return ast.NewArray(pass.processType(def.ValueType))
}

func (pass *NotRequiredFieldAsNullableType) processMap(def ast.MapType) ast.Type {
	return ast.NewMap(
		pass.processType(def.IndexType),
		pass.processType(def.ValueType),
	)
}

func (pass *NotRequiredFieldAsNullableType) processStruct(def ast.StructType) ast.Type {
	processedFields := make([]ast.StructField, 0, len(def.Fields))
	for _, field := range def.Fields {
		fieldType := pass.processType(field.Type)
		if !field.Required {
			fieldType.Nullable = true
		}

		newField := field
		newField.Type = fieldType

		processedFields = append(processedFields, newField)
	}

	return ast.NewStruct(processedFields...)
}
