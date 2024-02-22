package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*FieldsSetRequired)(nil)

// FieldsSetRequired rewrites the definition of given fields to mark them as not nullable and required.
type FieldsSetRequired struct {
	Fields []FieldReference
}

func (pass *FieldsSetRequired) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *FieldsSetRequired) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(pass.processObject)

	return schema
}

func (pass *FieldsSetRequired) processObject(_ string, object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	for i, field := range object.Type.AsStruct().Fields {
		for _, fieldRef := range pass.Fields {
			if !fieldRef.Matches(object, field) {
				continue
			}

			field.Type.Nullable = false
			field.Required = true
			field.AddToPassesTrail("FieldsSetRequired[nullable=false, required=true]")

			object.Type.Struct.Fields[i] = field
		}
	}

	return object
}
