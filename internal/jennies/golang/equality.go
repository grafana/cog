package golang

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/imports"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

type equalityMethods struct {
	tmpl            *template.Template
	apiRefCollector *apiref.APIReferenceCollector
}

func newEqualityMethods(tmpl *template.Template, apiRefCollector *apiref.APIReferenceCollector) equalityMethods {
	return equalityMethods{
		tmpl:            tmpl,
		apiRefCollector: apiRefCollector,
	}
}

func (jenny equalityMethods) generateForObject(buffer *strings.Builder, context languages.Context, object ir.Object, imports *imports.DirectImportMap) error {
	if !object.Type.IsStruct() {
		return nil
	}

	jenny.apiRefCollector.ObjectMethod(object, apiref.MethodReference{
		Name: "Equals",
		Arguments: []apiref.ArgumentReference{
			{Name: "other", Type: formatObjectName(object.Name)},
		},
		Comments: []string{
			fmt.Sprintf("Equals tests the equality of two `%s` objects.", formatObjectName(object.Name)),
		},
		Return: "bool",
	})

	tmpl := jenny.tmpl.
		Funcs(template.TypeResolvingHelpers(context)).
		Funcs(template.FuncMap{
			"typeHasEqualityFunc": func(typeDef ir.Type) bool {
				if !typeDef.IsRef() {
					return false
				}

				return context.ResolveToStruct(typeDef)
			},
			"resolveRefs": context.ResolveRefs,
			"importStdPkg": func(pkg string) string {
				return imports.Add(pkg, pkg)
			},
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
