package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*AddFields)(nil)

// AddFields rewrites the definition of an object to add new fields.
// Note: existing fields will not be overwritten.
type AddFields struct {
	Object ObjectReference
	Fields []ast.StructField
}

func (pass *AddFields) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *AddFields) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		if !pass.Object.Matches(object) {
			return object
		}

		return pass.processObject(object)
	})

	return schema
}

func (pass *AddFields) processObject(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	for _, field := range pass.Fields {
		// let's be safe: if a field with the same name already exists, we do not overwrite it.
		if _, exists := object.Type.AsStruct().FieldByName(field.Name); exists {
			continue
		}

		field.AddToPassesTrail("AddFields[created]")

		object.Type.Struct.Fields = append(object.Type.Struct.Fields, field)
	}

	return object
}
