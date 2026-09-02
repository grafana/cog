package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
)

var _ Transform = (*DisjunctionPropagateVariant)(nil)

type DisjunctionPropagateVariant struct {
	schemas ir.Schemas
}

func (pass *DisjunctionPropagateVariant) Process(schemas ir.Schemas) (ir.Schemas, error) {
	pass.schemas = schemas

	visitor := &Visitor{
		OnObject: pass.processObject,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *DisjunctionPropagateVariant) processObject(_ *Visitor, schema *ir.Schema, obj ir.Object) (ir.Object, error) {
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

		referredObj, found := schema.GetObject(def.Ref.ReferredType)
		if !found {
			return referredObj, fmt.Errorf("could not resolve reference %s", def.Ref)
		}

		referredObj.Type.Hints[ir.HintImplementsVariant] = obj.Type.Hints[ir.HintImplementsVariant]
		referredObj.AddToPassesTrail("DisjunctionPropagateVariant")
	}

	return obj, nil
}
