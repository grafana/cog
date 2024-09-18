package golang

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type equalityMethods struct {
	tmpl *template.Template
}

func newEqualityMethods(tmpl *template.Template) equalityMethods {
	return equalityMethods{
		tmpl: tmpl,
	}
}

func (jenny equalityMethods) generateForObject(buffer *strings.Builder, context languages.Context, schema *ast.Schema, object ast.Object) error {
	if !object.Type.IsStruct() {
		return nil
	}

	tmpl := jenny.tmpl.Funcs(template.FuncMap{
		"typeHasEqualityFunc": func(typeDef ast.Type) bool {
			if !typeDef.IsRef() {
				return false
			}

			return context.ResolveToStruct(typeDef)
		},
		"resolvesToScalar": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsScalar()
		},
		"resolvesToMap": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsMap()
		},
		"resolvesToArray": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsArray()
		},
		"resolvesToEnum": func(typeDef ast.Type) bool {
			return context.ResolveRefs(typeDef).IsEnum()
		},
		"resolveRefs": context.ResolveRefs,
	})

	templateFile := "types/struct_equality_method.tmpl"
	if object.Type.IsDataqueryVariant() {
		templateFile = "types/dataquery_equality_method.tmpl"
	}

	rendered, err := tmpl.Render(templateFile, map[string]any{
		"def": object,
	})
	if err != nil {
		return err
	}
	buffer.WriteString(rendered)

	return nil
}
