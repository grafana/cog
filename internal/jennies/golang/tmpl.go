package golang

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
)

//go:embed templates/runtime/*.tmpl templates/builders/*.tmpl templates/converters/*.tmpl templates/types/*.tmpl
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
			"formatRawRef": func(_ ast.Type) string {
				panic("formatRawRef() needs to be overridden by a jenny")
			},
			"importStdPkg": func(_ string) string {
				panic("importStdPkg() needs to be overridden by a jenny")
			},
			"importPkg": func(_ string) string {
				panic("importPkg() needs to be overridden by a jenny")
			},
			"maybeUnptr": func(variableName string, intoType ast.Type) string {
				panic("maybeUnptr() needs to be overridden by a jenny")
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
			"typeHasEqualityFunc": func(_ ast.Type) bool {
				panic("typeHasEqualityFunc() needs to be overridden by a jenny")
			},
			"resolvesToScalar": func(typeDef ast.Type) bool {
				panic("refResolvesToScalar() needs to be overridden by a jenny")
			},
			"resolvesToMap": func(typeDef ast.Type) bool {
				panic("refResolvesToMap() needs to be overridden by a jenny")
			},
			"resolvesToArray": func(typeDef ast.Type) bool {
				panic("refResolvesToArray() needs to be overridden by a jenny")
			},
			"resolvesToEnum": func(typeDef ast.Type) bool {
				panic("refResolvesToEnum() needs to be overridden by a jenny")
			},
			"resolveRefs": func(typeDef ast.Type) ast.Type {
				panic("resolveRefs() needs to be overridden by a jenny")
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
			"maybeDereference": func(typeDef ast.Type) string {
				if typeDef.Nullable && !typeDef.IsAnyOf(ast.KindArray, ast.KindMap) {
					return "*"
				}

				return ""
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
