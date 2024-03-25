package template

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	gotemplate "text/template"

	"github.com/grafana/cog/internal/tools"
)

const recursionMaxNums = 1000

func Helpers(baseTemplate *gotemplate.Template) gotemplate.FuncMap {
	includedNames := make(map[string]int)
	include := func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", fmt.Errorf("unable to execute template: rendering template has a nested reference name: %s", name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := baseTemplate.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}

	return gotemplate.FuncMap{
		"add1": func(i int) int { return i + 1 },
		// https://github.com/Masterminds/sprig/blob/581758eb7d96ae4d113649668fa96acc74d46e7f/dict.go#L76
		"dict": func(v ...any) map[string]any {
			dict := map[string]any{}
			lenv := len(v)
			for i := 0; i < lenv; i += 2 {
				key := v[i].(string)
				if i+1 >= lenv {
					dict[key] = ""
					continue
				}
				dict[key] = v[i+1]
			}
			return dict
		},
		"sub": func(a, b int) int { return a - b },
		"ternary": func(valTrue any, valFalse any, condition bool) any {
			if condition {
				return valTrue
			}

			return valFalse
		},

		// ------- \\
		// Strings \\
		// ------- \\
		"indent": func(spaces int, input string) string {
			return tools.Indent(input, spaces)
		},
		// Parameter order is reversed to stay compatible with sprig: https://github.com/Masterminds/sprig/blob/581758eb7d96ae4d113649668fa96acc74d46e7f/strings.go#L199
		"join": func(separator string, input []string) string {
			return strings.Join(input, separator)
		},
		"lower":          strings.ToLower,
		"lowerCamelCase": tools.LowerCamelCase,
		// Parameter order is reversed to stay compatible with sprig: https://github.com/Masterminds/sprig/blob/581758eb7d96ae4d113649668fa96acc74d46e7f/functions.go#L135
		"trimPrefix":     func(a, b string) string { return strings.TrimPrefix(b, a) },
		"upperCamelCase": tools.UpperCamelCase,

		// --------- \\
		// Templates \\
		// --------- \\
		"include": include,
		"includeIfExists": func(name string, data interface{}) (string, error) {
			if tmpl := baseTemplate.Lookup(name); tmpl == nil {
				return "", nil
			}

			return include(name, data)
		},
	}
}

func FindAndParseTemplates(vfs fs.FS, rootTmpl *gotemplate.Template, rootDir string) (*gotemplate.Template, error) {
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
