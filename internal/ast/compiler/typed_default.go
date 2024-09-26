package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*TypedDefaults)(nil)

type TypedDefaults struct {
}

func (t *TypedDefaults) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := Visitor{
		OnStructField: t.processField,
	}

	return visitor.VisitSchemas(schemas)
}

func (t *TypedDefaults) processField(_ *Visitor, _ *ast.Schema, f ast.StructField) (ast.StructField, error) {
	if f.Type.Default == nil {
		return f, nil
	}

	f.Type.TypedDefault = &f.Type
	return f, nil
}
