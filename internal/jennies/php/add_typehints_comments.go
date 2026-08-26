package php

import (
	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/ir/transforms"
	"github.com/grafana/cog/internal/languages"
)

var _ transforms.Transform = (*AddTypehintsComments)(nil)

type AddTypehintsComments struct {
	config Config
	hinter *typehints
}

func (pass *AddTypehintsComments) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.hinter = &typehints{
		config:  pass.config,
		context: languages.Context{Schemas: schemas},
	}

	visitor := &transforms.Visitor{
		OnStructField: pass.processStructField,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *AddTypehintsComments) processStructField(_ *transforms.Visitor, _ *ir.Schema, field ir.StructField) (ir.StructField, error) {
	if !pass.hinter.requiresHint(field.Type) {
		return field, nil
	}

	hint := pass.hinter.varAnnotationForType(field.Type)
	if hint != "" {
		field.Comments = append(field.Comments, hint)
	}

	return field, nil
}
