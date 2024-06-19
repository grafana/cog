package terraform

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	cogtemplate "github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/types/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("terraform")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		// placeholder functions, will be overridden by jennies
		Funcs(template.FuncMap{
			"formatTerraformType": func(_ string) string {
				panic("formatTerraformType() needs to be overriden by a jenny")
			},
		}).
		Funcs(map[string]any{
			"formatObjectName":     tools.UpperCamelCase,
			"formatFieldName":      tools.UpperCamelCase,
			"formatFieldNameTFSDK": tools.SnakeCase,
		})

	templates = template.Must(cogtemplate.FindAndParseTemplates(veneersFS, base, "templates"))
}

func renderTemplate(templateFile string, data map[string]any) (string, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, templateFile, data); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}
