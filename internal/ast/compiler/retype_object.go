package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*RetypeObject)(nil)

type RetypeObject struct {
	Object ObjectReference
	As     ast.Type
}

func (pass *RetypeObject) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	for i, schema := range schemas {
		schemas[i] = pass.processSchema(schema)
	}

	return schemas, nil
}

func (pass *RetypeObject) processSchema(schema *ast.Schema) *ast.Schema {
	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		return pass.processObject(object)
	})

	return schema
}

func (pass *RetypeObject) processObject(object ast.Object) ast.Object {
	if !pass.Object.Matches(object) {
		return object
	}

	trailMessage := fmt.Sprintf("RetypeObject[%s → %s]", ast.TypeName(object.Type), ast.TypeName(pass.As))

	object.Type = pass.As
	object.AddToPassesTrail(trailMessage)

	return object
}
