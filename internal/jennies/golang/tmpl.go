package golang

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/runtime/*.tmpl templates/builders/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

const recursionMaxNums = 1000

//nolint:gochecknoinits
func init() {
	includedNames := make(map[string]int)

	base := template.New("golang")
	base.
		Option("missingkey=error").
		Funcs(sprig.FuncMap()).
		Funcs(map[string]any{
			// placeholder functions, will be overridden by jennies
			"formatPath":               func(_ ast.Path) string { return "" },
			"formatType":               func(_ ast.Type) string { return "" },
			"typeHasBuilder":           func(_ ast.Type) bool { return false },
			"resolvesToComposableSlot": func(_ ast.Type) bool { return false },

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
						return "", fmt.Errorf("unable to execute template: rendering template has a nested reference name: %s", name)
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

	templates = template.Must(findAndParseTemplates(veneersFS, base, "templates"))
}

func findAndParseTemplates(vfs fs.FS, rootTmpl *template.Template, rootDir string) (*template.Template, error) {
	err := fs.WalkDir(vfs, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileHandle, err := vfs.Open(path)
		if err != nil {
			return err
		}

		contents, err := io.ReadAll(fileHandle)
		if err != nil {
			return err
		}

		templateName := strings.TrimPrefix(strings.TrimPrefix(path, rootDir), "/")
		t := rootTmpl.New(templateName)
		_, err = t.Parse(string(contents))

		return err
	})

	return rootTmpl, err
}

func renderTemplate(templateFile string, data map[string]any) (string, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, templateFile, data); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}
