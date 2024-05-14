package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*RetypeField)(nil)

type RetypeField struct {
	Field    FieldReference
	As       ast.Type
	Comments []string
}

func (pass *RetypeField) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *RetypeField) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(object)
	})

	return schema
}

func (pass *RetypeField) processObject(object ast.Object) ast.Object {
	if !object.Type.IsStruct() {
		return object
	}

	for i, field := range object.Type.Struct.Fields {
		if !pass.Field.Matches(object, field) {
			continue
		}

		object.Type.Struct.Fields[i].AddToPassesTrail(fmt.Sprintf("RetypeField[%s → %s]", ast.TypeName(field.Type), ast.TypeName(pass.As)))
		object.Type.Struct.Fields[i].Type = pass.As

		if pass.Comments != nil {
			object.Type.Struct.Fields[i].Comments = pass.Comments
		}

		break
	}

	return object
}
