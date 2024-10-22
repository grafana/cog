package golang

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type equalityMethods struct {
	tmpl            *template.Template
	apiRefCollector *common.APIReferenceCollector
}

func newEqualityMethods(tmpl *template.Template, apiRefCollector *common.APIReferenceCollector) equalityMethods {
	return equalityMethods{
		tmpl:            tmpl,
		apiRefCollector: apiRefCollector,
	}
}

func (jenny equalityMethods) generateForObject(buffer *strings.Builder, context languages.Context, object ast.Object, imports *common.DirectImportMap) error {
	if !object.Type.IsStruct() {
		return nil
	}

	jenny.apiRefCollector.RegisterMethod(object, common.MethodReference{
		Name: "Equals",
		Arguments: []common.ArgumentReference{
			{Name: "other", Type: tools.UpperCamelCase(object.Name)},
		},
		Return: "bool",
	})

	tmpl := jenny.tmpl.
		Funcs(common.TypeResolvingTemplateHelpers(context)).
		Funcs(template.FuncMap{
			"typeHasEqualityFunc": func(typeDef ast.Type) bool {
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
