package golang

import (
	"embed"
	"html/template"
	"strings"

	"github.com/grafana/cog/internal/tools"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed veneers/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("golang")
	base.Funcs(map[string]any{
		"formatIdentifier": tools.UpperCamelCase,
		"lowerCamelCase":   tools.LowerCamelCase,
		"formatType":       formatType,
		"trimPrefix":       strings.TrimPrefix,
	})
	templates = template.Must(base.ParseFS(veneersFS, "veneers/*.tmpl"))
}
