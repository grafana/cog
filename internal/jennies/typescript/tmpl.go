package typescript

import (
	"embed"
	"text/template"

	"github.com/grafana/cog/internal/ast"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
)

//go:embed templates/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	base := template.New("ts")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		// placeholder functions, will be overridden by jennies
		Funcs(template.FuncMap{
			"formatType": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"formatIdentifier": formatIdentifier,
			"typeIsDisjunctionOfBuilders": func(_ ast.Type) string {
				panic("typeIsDisjunctionOfBuilders() needs to be overridden by a jenny")
			},
			"defaultValueForType": func(_ ast.Type) string {
				panic("defaultValueForType() needs to be overridden by a jenny")
			},
			"formatValue": func(destinationType ast.Type, value any) string {
				panic("formatValue() needs to be overridden by a jenny")
			},
			"formatPath": func(_ ast.Path) string {
				panic("formatPath() needs to be overridden by a jenny")
			},
			"emptyValueForGuard": func(_ ast.Type) string {
				panic("emptyValueForGuard() needs to be overridden by a jenny")
			},
			"typeHasBuilder": func(_ ast.Type) bool {
				panic("typeHasBuilder() needs to be overridden by a jenny")
			},
			"resolvesToComposableSlot": func(_ ast.Type) bool {
				panic("resolvesToComposableSlot() needs to be overridden by a jenny")
			},
		})

	templates := template.Must(cogtemplate.FindAndParseTemplatesFromFS(templatesFS, base, "templates"))

	return template.Must(cogtemplate.FindAndParseTemplates(templates, extraTemplatesDirectories...))
}
