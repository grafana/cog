package golang

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/cog/internal/ast"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/tools"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/runtime/*.tmpl templates/builders/*.tmpl templates/builders/veneers/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("golang")
	base.
		Option("missingkey=error").
		Funcs(sprig.FuncMap()).
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
		}).
		Funcs(map[string]any{
			"formatPackageName": formatPackageName,
			"formatScalar":      formatScalar,
			"formatArgName": func(name string) string {
				return escapeVarName(tools.LowerCamelCase(name))
			},
			"maybeAsPointer": func(intoType ast.Type, variableName string) string {
				if intoType.Nullable && !(intoType.IsArray() || intoType.IsMap()) {
					return "&" + variableName
				}

				return variableName
			},
			"isNullableNonArray": func(typeDef ast.Type) bool {
				return typeDef.Nullable && !typeDef.IsArray()
			},
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
