package template

import (
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/grafana/cog/internal/ast"
)

//go:embed templates/*/builders/*.tmpl templates/*/builders/veneers/*.tmpl templates/*/runtime/*.tmpl templates/*/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

type Template struct {
	templates *template.Template
}

func NewTemplate(name string, functionMap template.FuncMap) *Template {
	base := template.New(name)
	base.Option("missingkey=error").
		Funcs(sprig.FuncMap()).
		Funcs(Helpers(base)).
		Funcs(template.FuncMap{
			"formatType": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"formatPath": func(_ ast.Path) string {
				panic("formatPath() needs to be overridden by a jenny")
			},
			"formatTypeNoBuilder": func(_ ast.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"defaultValueForType": func(_ ast.Type) string {
				panic("defaultValueForType() needs to be overridden by a jenny")
			},
			"formatArgName": func(name string) string {
				panic("defaultValueForType() needs to be overridden by a jenny")
			},
			"formatScalar": func(name string) string {
				panic("formatScalar() needs to be overridden by a jenny")
			},
		}).
		Funcs(functionMap)

	return &Template{
		templates: template.Must(FindAndParseTemplates(templatesFS, base, fmt.Sprintf("templates/%s", name))),
	}
}

func (t *Template) AddFuncMap(fnMap template.FuncMap) *Template {
	t.templates.Funcs(fnMap)
	return t
}

func (t *Template) Execute(templateName string, builder Builder) (string, error) {
	var buffer strings.Builder
	err := t.templates.ExecuteTemplate(&buffer, templateName, builder)
	return buffer.String(), err
}
