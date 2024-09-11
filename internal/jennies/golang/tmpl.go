package golang

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
)

//go:embed templates/runtime/*.tmpl templates/builders/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	tmpl, err := template.New(
		"golang",

		// placeholder functions, will be overridden by jennies
		template.Funcs(template.FuncMap{
			"formatPath": func(_ ast.Path) string {
				panic("formatPath() needs to be overridden by a jenny")
			},
			"formatType": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"formatTypeNoBuilder": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
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
		}),
		template.Funcs(map[string]any{
			"formatPackageName": formatPackageName,
			"formatScalar":      formatScalar,
			"formatArgName":     formatArgName,
			"maybeAsPointer": func(intoType ast.Type, variableName string) string {
				if intoType.Nullable && !(intoType.IsArray() || intoType.IsMap() || intoType.IsComposableSlot()) {
					return "&" + variableName
				}

				return variableName
			},
			"isNullableNonArray": func(typeDef ast.Type) bool {
				return typeDef.Nullable && !typeDef.IsArray()
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
