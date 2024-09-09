package golang

import (
	"embed"
	"text/template"

	"github.com/grafana/cog/internal/ast"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
)

//go:embed templates/runtime/*.tmpl templates/builders/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	base := template.New("golang")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		// placeholder functions, will be overridden by jennies
		Funcs(template.FuncMap{
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
		}).
		Funcs(map[string]any{
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
		})

	templates := template.Must(cogtemplate.FindAndParseTemplatesFromFS(veneersFS, base, "templates"))

	return template.Must(cogtemplate.FindAndParseTemplates(templates, extraTemplatesDirectories...))
}
