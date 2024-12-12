package python

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
)

//go:embed templates/*/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	tmpl, err := template.New(
		"python",

		template.Funcs(formattingTemplateFuncs()),
		// placeholder functions, will be overridden by jennies
		template.Funcs(template.FuncMap{
			"formatType": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"formatFullyQualifiedRef": func(_ ast.Type) string {
				panic("formatFullyQualifiedRef() needs to be overridden by a jenny")
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
			"importModule": func(alias string, pkg string, module string) string {
				panic("importModule() needs to be overridden by a jenny")
			},
			"importPkg": func(alias string, pkg string) string {
				panic("importPkg() needs to be overridden by a jenny")
			},
		}),

		// parse templates
		template.ParseFS(templatesFS, "templates"),
		template.ParseDirectories(extraTemplatesDirectories...),
	)
	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}

	return tmpl
}

func formattingTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatIdentifier": formatIdentifier,
		"formatPath":       formatFieldPath,
	}
}
