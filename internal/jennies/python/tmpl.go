package python

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/grafana/cog/internal/ast"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*/*.tmpl templates/builders/veneers/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("python")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		// placeholder functions, will be overridden by jennies
		Funcs(template.FuncMap{
			"formatType": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"formatRawType": func(_ ast.Type) string {
				panic("formatRawType() needs to be overridden by a jenny")
			},
			"formatRawTypeNotNullable": func(_ ast.Type) string {
				panic("formatRawTypeNotNullable() needs to be overridden by a jenny")
			},
			"formatValue": func(_ ast.Type, _ any) string {
				panic("formatValue() needs to be overridden by a jenny")
			},
			"defaultForType": func(_ ast.Type) string {
				panic("defaultForType() needs to be overridden by a jenny")
			},
			"typeHasBuilder": func(_ ast.Type) bool {
				panic("typeHasBuilder() needs to be overridden by a jenny")
			},
			"resolvesToComposableSlot": func(_ ast.Type) bool {
				panic("resolvesToComposableSlot() needs to be overridden by a jenny")
			},
		}).
		Funcs(template.FuncMap{
			"formatIdentifier": formatIdentifier,
			"formatPath":       formatFieldPath,
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
