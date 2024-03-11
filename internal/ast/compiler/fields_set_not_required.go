package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*FieldsSetNotRequired)(nil)

// FieldsSetNotRequired rewrites the definition of given fields to mark them as nullable and not required.
type FieldsSetNotRequired struct {
	Fields []FieldReference
}

func (pass *FieldsSetNotRequired) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *FieldsSetNotRequired) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(pass.processObject)

	return schema
}

func (pass *FieldsSetNotRequired) processObject(_ string, object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	for i, field := range object.Type.AsStruct().Fields {
		for _, fieldRef := range pass.Fields {
			if !fieldRef.Matches(object, field) {
				continue
			}

			field.Type.Nullable = true
			field.Required = false
			field.AddToPassesTrail("FieldsSetNotRequired[nullable=true, required=false]")

			object.Type.Struct.Fields[i] = field
		}
	}

	return object
}
