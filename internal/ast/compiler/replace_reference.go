package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

// ReplaceReference replaces any usage of the `From` reference by the one given in `To`.
type ReplaceReference struct {
	From     ObjectReference
	To       ObjectReference
	refFound bool
}

func (pass *ReplaceReference) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.refFound = false

	visitor := Visitor{
		OnRef: pass.processRef,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *ReplaceReference) processRef(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if !pass.From.MatchesRef(def.AsRef()) {
		return def, nil
	}

	pass.refFound = true

	return ast.NewRef(pass.To.Package, pass.To.Object, ast.Trail(fmt.Sprintf("ReplaceReference[%s → %s]", def.Ref, pass.To))), nil
}

func (pass *ReplaceReference) Diagnostics() []string {
	if pass.refFound {
		return nil
	}

	return []string{
		fmt.Sprintf("reference '%s' not found", pass.From),
	}
}
