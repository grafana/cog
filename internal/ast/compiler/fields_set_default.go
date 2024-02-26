package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*FieldsSetDefault)(nil)

// FieldsSetDefault sets the default value for the given fields.
type FieldsSetDefault struct {
	DefaultValues map[FieldReference]any
}

func (pass *FieldsSetDefault) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *FieldsSetDefault) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(pass.processObject)

	return schema
}

func (pass *FieldsSetDefault) processObject(_ string, object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	for i, field := range object.Type.AsStruct().Fields {
		for fieldRef, value := range pass.DefaultValues {
			if !fieldRef.Matches(object, field) {
				continue
			}

			field.Type.Default = value

			object.Type.Struct.Fields[i] = field
		}
	}

	return object
}
