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
	newSchema := schema.DeepCopy()
	for i, object := range schema.Objects {
		newSchema.Objects[i] = pass.processObject(object)
	}

	return &newSchema
}

func (pass *NotRequiredFieldAsNullableType) processObject(object ast.Object) ast.Object {
	if object.Type.Kind != ast.KindStruct {
		return object
	}

	object.Type = pass.processType(object.Type)

	return object
}

func (pass *NotRequiredFieldAsNullableType) processType(def ast.Type) ast.Type {
	if def.Kind == ast.KindArray {
		return pass.processArray(def)
	}

	if def.Kind == ast.KindMap {
		return pass.processMap(def)
	}

	if def.Kind == ast.KindStruct {
		return pass.processStruct(def)
	}

	return def
}

func (pass *NotRequiredFieldAsNullableType) processArray(def ast.Type) ast.Type {
	def.Array.ValueType = pass.processType(def.Array.ValueType)

	return def
}

func (pass *NotRequiredFieldAsNullableType) processMap(def ast.Type) ast.Type {
	def.Map.IndexType = pass.processType(def.Map.IndexType)
	def.Map.ValueType = pass.processType(def.Map.ValueType)

	return def
}

func (pass *NotRequiredFieldAsNullableType) processStruct(def ast.Type) ast.Type {
	for i, field := range def.Struct.Fields {
		fieldType := pass.processType(field.Type)
		if !field.Required {
			fieldType.Nullable = true
		}

		newField := field
		newField.Type = fieldType

		def.Struct.Fields[i] = newField
	}

	return def
}
