package typescript

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
)

//go:embed templates/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	tmpl, err := template.New(
		"typescript",

		// placeholder functions, will be overridden by jennies
		template.Funcs(template.FuncMap{
			"formatType": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"formatIdentifier": formatIdentifier,
			"typeIsDisjunctionOfBuilders": func(_ ast.Type) string {
				panic("typeIsDisjunctionOfBuilders() needs to be overridden by a jenny")
			},
			"defaultValueForType": func(_ ast.Type) string {
				panic("defaultValueForType() needs to be overridden by a jenny")
			},
			"formatValue": func(destinationType ast.Type, value any) string {
				panic("formatValue() needs to be overridden by a jenny")
			},
			"formatPath": func(_ ast.Path) string {
				panic("formatPath() needs to be overridden by a jenny")
			},
			"emptyValueForGuard": func(_ ast.Type) string {
				panic("emptyValueForGuard() needs to be overridden by a jenny")
			},
			"typeHasBuilder": func(_ ast.Type) bool {
				panic("typeHasBuilder() needs to be overridden by a jenny")
			},
			"resolvesToComposableSlot": func(_ ast.Type) bool {
				panic("resolvesToComposableSlot() needs to be overridden by a jenny")
			},
			"importPkg": func(pkg string) string {
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
