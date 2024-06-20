package terraform

import (
	"embed"
	"text/template"

	"github.com/grafana/cog/internal/ast"
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
		Funcs(map[string]any{
			"formatObjectName":     tools.UpperCamelCase,
			"formatFieldName":      tools.UpperCamelCase,
			"formatFieldNameTFSDK": tools.SnakeCase,
			"formatTerraformType":  formatTerraformType,
			"filterScalars": func(list []ast.StructField) []ast.StructField {
				newList := []ast.StructField{}
				for _, f := range list {
					if f.Type.IsScalar() {
						newList = append(newList, f)
					}
				}
				return newList
			},
		})

	templates = template.Must(cogtemplate.FindAndParseTemplates(veneersFS, base, "templates"))
}
