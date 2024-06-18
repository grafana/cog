package php

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

//go:embed templates/builders/*.tmpl templates/runtime/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("php")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		// placeholder functions, will be overridden by jennies
		Funcs(template.FuncMap{
			"fullNamespaceRef": func(_ string) string {
				panic("fullNamespaceRef() needs to be overridden by a jenny")
			},
			"formatPath": func(_ ast.Path) string {
				panic("formatPath() needs to be overridden by a jenny")
			},
			"formatType":               func(_ ast.Type) string { panic("formatType() needs to be overridden by a jenny") },
			"formatRawType":            func(_ ast.Type) string { panic("formatRawType() needs to be overridden by a jenny") },
			"formatRawTypeNotNullable": func(_ ast.Type) string { panic("formatRawTypeNotNullable() needs to be overridden by a jenny") },
			"typeHasBuilder": func(_ ast.Type) bool {
				panic("typeHasBuilder() needs to be overridden by a jenny")
			},
			"typeHint": func(_ ast.Type) string {
				panic("typeHint() needs to be overridden by a jenny")
			},
			"isDisjunctionOfBuilders": func(_ ast.Type) bool {
				panic("isDisjunctionOfBuilders() needs to be overridden by a jenny")
			},
			"resolvesToComposableSlot": func(_ ast.Type) bool {
				panic("resolvesToComposableSlot() needs to be overridden by a jenny")
			},
			"defaultForType": func(_ ast.Type) bool {
				panic("defaultForType() needs to be overridden by a jenny")
			},
			"formatValue": func(_ ast.Type, _ any) bool {
				panic("formatValue() needs to be overridden by a jenny")
			},
		}).
		Funcs(map[string]any{
			"formatPackageName":    formatPackageName,
			"formatObjectName":     formatObjectName,
			"formatOptionName":     formatOptionName,
			"formatEnumMemberName": formatEnumMemberName,
			"formatArgName":        formatArgName,
			"formatScalar":         formatValue,
			"formatDocsBlock":      formatCommentsBlock,
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
