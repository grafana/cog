package compiler

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
)

var _ Pass = (*DisjunctionPropagateVariant)(nil)

type DisjunctionPropagateVariant struct {
	schemas ast.Schemas
}

func (pass *DisjunctionPropagateVariant) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	pass.schemas = schemas

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DisjunctionPropagateVariant) processObject(_ *Visitor, schema *ast.Schema, obj ast.Object) (ast.Object, error) {
	if !obj.Type.ImplementsVariant() {
		return obj, nil
	}

	if !obj.Type.IsDisjunction() {
		return obj, nil
	}

	for _, def := range obj.Type.Disjunction.Branches {
		if !def.IsRef() {
			continue
		}

		referredObj, found := schema.LocateObject(def.Ref.ReferredType)
		if !found {
			return referredObj, fmt.Errorf("could not resolve reference %s", def.Ref)
		}

		referredObj.Type.Hints[ast.HintImplementsVariant] = obj.Type.Hints[ast.HintImplementsVariant]
		referredObj.AddToPassesTrail("DisjunctionPropagateVariant")
	}

	return obj, nil
}
