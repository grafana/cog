package golang

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/grafana/cog/internal/tools"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("golang")
	base.Funcs(map[string]any{
		"formatIdentifier": tools.UpperCamelCase,
		"lowerCamelCase":   tools.LowerCamelCase,
		"toLower":          strings.ToLower,
		"formatType":       formatType,
		"trimPrefix":       strings.TrimPrefix,
	})
	templates = template.Must(base.ParseFS(veneersFS, "templates/*.tmpl"))
}

func renderTemplate(templateFile string, data map[string]any) (string, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, templateFile, data); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}
