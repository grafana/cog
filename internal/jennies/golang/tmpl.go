package golang

import (
	"embed"
	"html/template"

	"github.com/grafana/cog/internal/tools"
)

var templates *template.Template

//go:embed veneers/*.tmpl
var veneersFS embed.FS

func init() {
	base := template.New("golang")
	base.Funcs(map[string]any{
		"formatIdentifier": tools.UpperCamelCase,
	})
	templates = template.Must(base.ParseFS(veneersFS, "veneers/*.tmpl"))
}
