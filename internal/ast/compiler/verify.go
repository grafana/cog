package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*Verify)(nil)

type Verify struct {
	schemas ast.Schemas
}

func (pass *Verify) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.schemas = schemas

	visitor := &Visitor{
		OnRef: pass.processRef,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *Verify) processRef(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	_, found := pass.schemas.LocateObjectByRef(def.AsRef())
	if !found {
		return ast.Type{}, fmt.Errorf("could not resolve reference to %s", def.Ref.String())
	}

	return def, nil
}
