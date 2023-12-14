package common

import (
	"github.com/grafana/cog/internal/ast"
)

//nolint:musttag
type Context struct {
	Schemas  ast.Schemas
	Builders ast.Builders
}

func (context *Context) LocateObject(pkg string, name string) (ast.Object, bool) {
	return context.Schemas.LocateObject(pkg, name)
}

func (context *Context) ResolveToBuilder(def ast.Type) bool {
	if def.Kind == ast.KindArray {
		return context.ResolveToBuilder(def.AsArray().ValueType)
	}

	if def.Kind == ast.KindDisjunction {
		for _, branch := range def.AsDisjunction().Branches {
			if !context.ResolveToBuilder(branch) {
				return false
			}
		}

		return true
	}

	if def.Kind != ast.KindRef {
		return false
	}

	ref := def.AsRef()
	_, found := context.Builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)

	return found
}

func (context *Context) ResolveToComposableSlot(def ast.Type) (ast.Type, bool) {
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

func (context *Context) ResolveToStruct(def ast.Type) bool {
	if def.IsStruct() {
		return true
	}

	if !def.IsRef() {
		return false
	}

	referredObj, found := context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
	if !found {
		return false
	}

	return context.ResolveToStruct(referredObj.Type)
}

type BuildOptions struct {
	Languages []string
}
