package typescript

import (
	"embed"
	"text/template"

	"github.com/grafana/cog/internal/ast"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("ts")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		// placeholder functions, will be overridden by jennies
		Funcs(template.FuncMap{
			"formatType": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
		}).
		Funcs(template.FuncMap{
			"formatScalar": formatScalar,
		})
	templates = template.Must(base.ParseFS(templatesFS, "templates/*.tmpl"))
}
