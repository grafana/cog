package context

import (
	"github.com/grafana/cog/internal/ast"
)

type Builders struct {
	Schemas  ast.Schemas
	Builders ast.Builders
}

func (context *Builders) BuilderForType(t ast.Type) (ast.Builder, bool) {
	if t.Kind != ast.KindRef {
		return ast.Builder{}, false
	}

	ref := t.AsRef()
	return context.Builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)
}

func (context *Builders) LocateObject(pkg string, name string) (ast.Object, bool) {
	return context.Schemas.LocateObject(pkg, name)
}

func (context *Builders) ResolveToComposableSlot(def ast.Type) (ast.Type, bool) {
	if def.Kind == ast.KindComposableSlot {
		return def, true
	}

	if def.Kind == ast.KindArray {
		return context.ResolveToComposableSlot(def.AsArray().ValueType)
	}

	if def.Kind == ast.KindRef {
		referredObj, found := context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
		if !found {
			return ast.Type{}, false
		}

		return context.ResolveToComposableSlot(referredObj.Type)
	}

	return ast.Type{}, false
}
