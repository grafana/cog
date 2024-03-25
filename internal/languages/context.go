package languages

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

func (context *Context) LocateObjectByRef(ref ast.RefType) (ast.Object, bool) {
	return context.Schemas.LocateObjectByRef(ref)
}

func (context *Context) ResolveToBuilder(def ast.Type) bool {
	_, ok := context.ResolveAsBuilder(def)

	return ok
}

func (context *Context) ResolveAsBuilder(def ast.Type) (ast.Builder, bool) {
	if def.IsArray() {
		return context.ResolveAsBuilder(def.AsArray().ValueType)
	}

	if def.IsDisjunction() {
		for _, branch := range def.AsDisjunction().Branches {
			if builder, found := context.ResolveAsBuilder(branch); found {
				return builder, true
			}
		}

		return ast.Builder{}, false
	}

	if !def.IsRef() {
		return ast.Builder{}, false
	}

	ref := def.AsRef()
	return context.Builders.LocateByObject(ref.ReferredPkg, ref.ReferredType)
}

func (context *Context) IsDisjunctionOfBuilders(def ast.Type) bool {
	if !def.IsDisjunction() {
		return false
	}

	for _, branch := range def.AsDisjunction().Branches {
		if !context.ResolveToBuilder(branch) {
			return false
		}
	}

	return true
}

func (context *Context) ResolveToComposableSlot(def ast.Type) (ast.Type, bool) {
	if def.IsComposableSlot() {
		return def, true
	}

	if def.IsArray() {
		return context.ResolveToComposableSlot(def.AsArray().ValueType)
	}

	if def.IsRef() {
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

func (context *Context) ResolveRefs(def ast.Type) ast.Type {
	if !def.IsRef() {
		return def
	}

	referredObj, found := context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
	if !found {
		return def
	}

	return context.ResolveRefs(referredObj.Type)
}
