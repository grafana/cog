package template

import (
	"encoding/json"

	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

func TypeResolvingHelpers(context languages.Context) FuncMap {
	return FuncMap{
		"resolvesToScalar": func(typeDef ir.Type) bool {
			return context.ResolveRefs(typeDef).IsScalar()
		},
		"resolvesToArray": func(typeDef ir.Type) bool {
			return context.ResolveRefs(typeDef).IsArray()
		},
		"resolvesToMap": func(typeDef ir.Type) bool {
			return context.ResolveRefs(typeDef).IsMap()
		},
		"resolvesToEnum": func(typeDef ir.Type) bool {
			return context.ResolveRefs(typeDef).IsEnum()
		},
		"resolvesToStruct": func(typeDef ir.Type) bool {
			return context.ResolveRefs(typeDef).IsStruct()
		},
		"resolvesToDisjunction": func(typeDef ir.Type) bool {
			return context.ResolveRefs(typeDef).IsDisjunction()
		},
		"resolvesToBuilder": context.ResolveToBuilder,
		"resolveRefs":       context.ResolveRefs,
		"resolveRefsChain":  context.ResolveRefsChain,
		"resolvesToComposableSlot": func(typeDef ir.Type) bool {
			_, found := context.ResolveToComposableSlot(typeDef)
			return found
		},
	}
}

func TypesHelpers(context languages.Context) FuncMap {
	return FuncMap{
		"dumpJson": func(input any) string {
			payload, err := json.Marshal(input)
			if err != nil {
				panic("could not marshal to JSON")
			}

			return string(payload)
		},
		"schemaHasObject": func(schema *ir.Schema, name string) bool {
			return schema.HasObject(name)
		},
		"objectExists": func(pkg string, name string) bool {
			_, ok := context.Schemas.GetObject(pkg, name)
			return ok
		},
	}
}
