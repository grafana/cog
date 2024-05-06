package typescript

import (
	"embed"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/cog/internal/ast"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*.tmpl templates/veneers/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("ts")
	base.
		Option("missingkey=error").
		Funcs(sprig.FuncMap()).
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
		})
	templates = template.Must(cogtemplate.FindAndParseTemplates(templatesFS, base, "templates"))
}
