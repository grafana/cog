package languages

import (
	"sort"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/tools"
)

//nolint:musttag
type Context struct {
	Schemas         ir.Schemas
	Builders        ir.Builders
	ConverterConfig ConverterConfig
}

func (context *Context) LocateObject(pkg string, name string) (ir.Object, bool) {
	return context.Schemas.GetObject(pkg, name)
}

func (context *Context) LocateObjectByRef(ref ir.RefType) (ir.Object, bool) {
	return context.Schemas.GetObjectByRef(ref)
}

func (context *Context) ResolveToBuilder(def ir.Type) bool {
	if def.IsArray() {
		return context.ResolveToBuilder(def.AsArray().ValueType)
	}

	if def.IsMap() {
		return context.ResolveToBuilder(def.AsMap().ValueType)
	}

	if def.IsDisjunction() {
		for _, branch := range def.AsDisjunction().Branches {
			if found := context.ResolveToBuilder(branch); found {
				return true
			}
		}

		return false
	}

	if !def.IsRef() {
		return false
	}

	resolvedRef := context.ResolveRefs(def)
	if resolvedRef.IsDisjunction() {
		return context.ResolveToBuilder(resolvedRef)
	}

	return len(context.Builders.LocateAllByRef(def.AsRef())) != 0
}

func (context *Context) IsDisjunctionOfBuilders(def ir.Type) bool {
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

func (context *Context) IsArrayOfKinds(def ir.Type, kinds ...ir.Kind) bool {
	def = context.ResolveRefs(def)
	if !def.IsArray() {
		return false
	}

	valueType := context.ResolveRefs(def.AsArray().ValueType)
	if valueType.IsArray() {
		return context.IsArrayOfKinds(valueType, kinds...)
	}

	return valueType.IsAnyOf(kinds...)
}

func (context *Context) IsMapOfKinds(def ir.Type, kinds ...ir.Kind) bool {
	def = context.ResolveRefs(def)
	if !def.IsMap() {
		return false
	}

	valueType := context.ResolveRefs(def.AsMap().ValueType)
	if valueType.IsMap() {
		return context.IsMapOfKinds(valueType, kinds...)
	}

	return valueType.IsAnyOf(kinds...)
}

func (context *Context) ResolveToComposableSlot(def ir.Type) (ir.Type, bool) {
	if def.IsComposableSlot() {
		return def, true
	}

	if def.IsArray() {
		return context.ResolveToComposableSlot(def.AsArray().ValueType)
	}

	if def.IsRef() {
		referredObj, found := context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
		if !found {
			return ir.Type{}, false
		}

		return context.ResolveToComposableSlot(referredObj.Type)
	}

	return ir.Type{}, false
}

func (context *Context) ResolveToStruct(def ir.Type) bool {
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

func (context *Context) ResolveRefsChain(def ir.Type) ir.Type {
	if !def.IsRef() {
		return def
	}

	referredObj, found := context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
	if !found {
		return def
	}

	if !referredObj.Type.IsRef() {
		return def
	}

	return context.ResolveRefsChain(referredObj.Type)
}

func (context *Context) ResolveRefs(def ir.Type) ir.Type {
	if !def.IsRef() {
		return def
	}

	referredObj, found := context.LocateObject(def.AsRef().ReferredPkg, def.AsRef().ReferredType)
	if !found {
		return def
	}

	return context.ResolveRefs(referredObj.Type)
}

func (context *Context) BuildersForType(typeDef ir.Type) ir.Builders {
	var candidateBuilders ir.Builders

	var search func(def ir.Type)
	search = func(def ir.Type) {
		if def.IsArray() {
			search(def.AsArray().ValueType)
			return
		}
		if def.IsMap() {
			search(def.AsMap().ValueType)
			return
		}

		if def.IsDisjunction() {
			for _, branch := range def.AsDisjunction().Branches {
				search(branch)
			}

			return
		}

		if !def.IsRef() {
			return
		}

		candidateBuilders = append(candidateBuilders, context.Builders.LocateAllByRef(def.AsRef())...)
	}

	search(typeDef)

	return candidateBuilders
}

func (context Context) PackagesForVariant(variant string) []string {
	return tools.Map(context.SchemasForVariant(variant), func(schema *ir.Schema) string {
		return schema.Package
	})
}

func (context Context) SchemasForVariant(variant string) []*ir.Schema {
	schemas := tools.Filter(context.Schemas, func(schema *ir.Schema) bool {
		return schema.Metadata.Kind == ir.SchemaKindComposable && string(schema.Metadata.Variant) == variant && schema.Metadata.Identifier != ""
	})

	sort.Slice(schemas, func(i int, j int) bool {
		return schemas[i].Package < schemas[j].Package
	})

	return schemas
}
