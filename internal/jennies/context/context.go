package context

import (
	"github.com/grafana/cog/internal/ast"
)

//nolint:musttag
type Builders struct {
	Schemas  ast.Schemas
	Builders ast.Builders
}

func (context *Builders) LocateObject(pkg string, name string) (ast.Object, bool) {
	return context.Schemas.LocateObject(pkg, name)
}

func (context *Builders) ResolveToBuilder(def ast.Type) bool {
	if def.Kind == ast.KindArray {
		return context.ResolveToBuilder(def.AsArray().ValueType)
	}

	if def.Kind != ast.KindRef {
		return false
	}

	ref := def.AsRef()
	_, found := context.Builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)

	return found
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
