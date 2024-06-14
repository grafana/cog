package php

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
)

var _ compiler.Pass = (*AddTypehintsComments)(nil)

type AddTypehintsComments struct {
	hinter *typehints
}

func (pass *AddTypehintsComments) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.hinter = &typehints{}

	visitor := &compiler.Visitor{
		OnStructField: pass.processStructField,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *AddTypehintsComments) processStructField(_ *compiler.Visitor, _ *ast.Schema, field ast.StructField) (ast.StructField, error) {
	if !pass.hinter.requiresHint(field.Type) {
		return field, nil
	}

	hint := pass.hinter.annotationForType(field.Type)
	if hint != "" {
		field.Comments = append(field.Comments, hint)
	}

	return field, nil
}
