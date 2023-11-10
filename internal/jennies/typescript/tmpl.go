package typescript

import (
	"embed"
	"text/template"

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
		Funcs(template.FuncMap{
			"formatScalar": formatScalar,
		})
	templates = template.Must(base.ParseFS(templatesFS, "templates/*.tmpl"))
}
