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
	tmpl := jenny.tmpl.Funcs(template.FuncMap{
		"typeHasEqualityFunc": func(typeDef ast.Type) bool {
			if !typeDef.IsRef() {
				return false
			}

			return context.ResolveToStruct(typeDef)
		},
		"refResolvesToEnum": func(typeDef ast.Type) bool {
			if !typeDef.IsRef() {
				return false
			}

			return context.ResolveRefs(typeDef).IsEnum()
		},
	})

	if object.Type.IsDataqueryVariant() && object.Type.IsStruct() {
		rendered, err := tmpl.Render("types/dataquery_equality_method.tmpl", map[string]any{
			"def": object,
		})
		if err != nil {
			return err
		}
		buffer.WriteString(rendered)

		return nil
	}

	if object.Type.IsStruct() {
		rendered, err := tmpl.Render("types/struct_equality_method.tmpl", map[string]any{
			"def": object,
		})
		if err != nil {
			return err
		}
		buffer.WriteString(rendered)

		return nil
	}

	return nil
}
