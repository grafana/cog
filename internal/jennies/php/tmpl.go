package php

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	cogtemplate "github.com/grafana/cog/internal/jennies/template"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/runtime/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("php")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		Funcs(map[string]any{
			"formatPackageName":    formatPackageName,
			"formatObjectName":     formatObjectName,
			"formatEnumMemberName": formatEnumMemberName,
			"formatValue":          formatValue,
		})

	templates = template.Must(cogtemplate.FindAndParseTemplates(templatesFS, base, "templates"))
}

func renderTemplate(templateFile string, data map[string]any) (string, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, templateFile, data); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}
