package golang

import (
	"embed"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/cog/internal/tools"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*.tmpl templates/veneers/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("golang")
	base.
		Funcs(sprig.FuncMap()).
		Funcs(map[string]any{
			"formatIdentifier": tools.UpperCamelCase,
			"lowerCamelCase":   tools.LowerCamelCase,
			"formatType":       formatType,
			"trimPrefix":       strings.TrimPrefix,
			"maybeAsPointer": func(intoNullable bool, variableName string) string {
				if intoNullable {
					return "&" + variableName
				}

				return variableName
			},
		})
	templates = template.Must(base.ParseFS(veneersFS, "templates/*.tmpl", "templates/veneers/*.tmpl"))
}
