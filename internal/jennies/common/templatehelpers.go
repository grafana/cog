package common

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

func TypeResolvingTemplateHelpers(context languages.Context) template.FuncMap {
	return template.FuncMap{
		"resolvesToScalar": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsScalar()
		},
		"resolvesToArray": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsArray()
		},
		"resolvesToMap": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsMap()
		},
		"resolvesToEnum": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsEnum()
		},
		"resolvesToStruct": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsStruct()
		},
		"resolvesToComposableSlot": func(typeDef ast.Type) bool {
			_, found := context.ResolveToComposableSlot(typeDef)
			return found
		},
		"resolveRefs": context.ResolveRefs,
	}
}
