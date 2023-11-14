package template

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	gotemplate "text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/cog/internal/ast"
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
		// placeholder functions, will be overridden by jennies
		"typeHasBuilder": func(_ ast.Type) bool {
			panic("typeHasBuilder() needs to be overridden by a jenny")
		},
		"resolvesToComposableSlot": func(_ ast.Type) bool {
			panic("resolvesToComposableSlot() needs to be overridden by a jenny")
		},

		"upperCamelCase": tools.UpperCamelCase,
		"lowerCamelCase": tools.LowerCamelCase,
		"include":        include,
		"include_if_exists": func(name string, data interface{}) (string, error) {
			if name == "builders/veneers/post_assignment_Dashboard_WithRow" {
				spew.Dump(name)
			}
			if tmpl := baseTemplate.Lookup(name); tmpl == nil {
				if name == "builders/veneers/post_assignment_Dashboard_WithRow" {
					spew.Dump("not found")
					spew.Dump(baseTemplate.DefinedTemplates())

					for _, known := range baseTemplate.Templates() {
						spew.Dump(known.Name())
					}
				}
				return "", nil
			}

			if name == "builders/veneers/post_assignment_Dashboard_WithRow" {
				spew.Dump("found")
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
