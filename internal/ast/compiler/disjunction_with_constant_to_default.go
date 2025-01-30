package compiler

import (
	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DisjunctionWithConstantToDefault)(nil)

type DisjunctionWithConstantToDefault struct {
}

func (pass *DisjunctionWithConstantToDefault) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	visitor := &Visitor{
		OnDisjunction: pass.processDisjunction,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DisjunctionWithConstantToDefault) processDisjunction(_ *Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	branches := def.Disjunction.Branches

	if len(branches) != 2 {
		return def, nil
	}

	if branches[0].Kind != branches[1].Kind {
		return def, nil
	}

	if !branches[0].IsScalar() {
		return def, nil
	}

	if branches[0].Scalar.ScalarKind != branches[1].Scalar.ScalarKind {
		return def, nil
	}

	if branches[0].IsConcrete() == branches[1].IsConcrete() {
		return def, nil
	}

	if branches[0].IsConcrete() {
		def = branches[1]
		def.Default = branches[0].Value
	} else {
		def = branches[0]
		def.Default = branches[1].Value
	}

	def.AddToPassesTrail("DisjunctionWithConstantToDefault")

	return def, nil
}
