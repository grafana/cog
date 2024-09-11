package terraform

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

//go:embed templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates() *template.Template {
	tmpl, err := template.New(
		"terraform",

		template.Funcs(map[string]any{
			"formatObjectName":     tools.UpperCamelCase,
			"formatFieldName":      tools.UpperCamelCase,
			"formatFieldNameTFSDK": tools.SnakeCase,
			"formatTerraformType":  formatTerraformType,
			"filterScalars": func(fields []ast.StructField) []ast.StructField {
				return tools.Filter(fields, func(field ast.StructField) bool {
					return field.Type.IsScalar()
				})
			},
		}),

		// parse templates
		template.ParseFS(templatesFS, "templates"),
	)
	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}

	return tmpl
}
