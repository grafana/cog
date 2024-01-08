package java

import (
	"embed"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
	"text/template"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/*.tmpl
//nolint:gochecknoglobals
var veneersFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("golang")
	base.
		Option("missingkey=error").
		Funcs(sprig.FuncMap()).
		Funcs(cogtemplate.Helpers(base)).
		Funcs(functions())

	templates = template.Must(cogtemplate.FindAndParseTemplates(veneersFS, base, "templates"))
}

func functions() template.FuncMap {
	return template.FuncMap{
		"lastItem": func(index int, values []EnumValue) bool {
			return len(values)-1 == index
		},
	}
}

type EnumTemplate struct {
	Package  string
	Name     string
	Values   []EnumValue
	Type     string
	Comments []string
}

type EnumValue struct {
	Name  string
	Value any
}

type ClassTemplate struct {
	Package  string
	Imports  fmt.Stringer
	Name     string
	Extends  []string
	Comments []string

	Fields       []Field
	InnerClasses []ClassTemplate

	GenGettersAndSetters bool
}

type Field struct {
	Name     string
	Type     string
	Comments []string
}

type ConstantTemplate struct {
	Package   string
	Name      string
	Constants []Constant
}

type Constant struct {
	Name  string
	Type  string
	Value any
}
