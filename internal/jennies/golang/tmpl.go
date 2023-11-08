package golang

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"github.com/pkg/errors"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*.tmpl templates/veneers/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

const recursionMaxNums = 1000

//nolint:gochecknoinits
func init() {
	includedNames := make(map[string]int)

	base := template.New("golang")
	base.
		Funcs(sprig.FuncMap()).
		Funcs(map[string]any{
			"formatPath":     func(_ ast.Path) string { return "" },  // placeholder function, will be overridden by jennies
			"formatType":     func(_ ast.Type) string { return "" },  // placeholder function, will be overridden by jennies
			"typeHasBuilder": func(_ ast.Type) bool { return false }, // placeholder function, will be overridden by jennies

			"formatIdentifier": tools.UpperCamelCase,
			"lowerCamelCase":   tools.LowerCamelCase,
			"formatScalar":     formatScalar,
			"formatArgName": func(name string) string {
				return escapeVarName(tools.LowerCamelCase(name))
			},
			"trimPrefix": strings.TrimPrefix,
			"maybeAsPointer": func(intoNullable bool, variableName string) string {
				if intoNullable {
					return "&" + variableName
				}

				return variableName
			},
			"isNullableNonArray": func(typeDef ast.Type) bool {
				return typeDef.Nullable && typeDef.Kind != ast.KindArray
			},
			"include": func(name string, data interface{}) (string, error) {
				var buf strings.Builder
				if v, ok := includedNames[name]; ok {
					if v > recursionMaxNums {
						return "", errors.Wrapf(fmt.Errorf("unable to execute template"), "rendering template has a nested reference name: %s", name)
					}
					includedNames[name]++
				} else {
					includedNames[name] = 1
				}
				err := base.ExecuteTemplate(&buf, name, data)
				includedNames[name]--
				return buf.String(), err
			},
		})
	templates = template.Must(base.ParseFS(veneersFS, "templates/*.tmpl", "templates/veneers/*.tmpl")).Option("missingkey=error")
}
