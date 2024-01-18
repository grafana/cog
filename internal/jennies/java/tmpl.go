package java

import (
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

//go:embed templates/runtime/*.tmpl templates/types/*.tmpl templates/builders/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("golang")
	base.
		Option("missingkey=error").
		Funcs(sprig.FuncMap()).
		Funcs(cogtemplate.Helpers(base)).
		Funcs(functions())

	templates = template.Must(cogtemplate.FindAndParseTemplates(templatesFS, base, "templates"))
}

func functions() template.FuncMap {
	return template.FuncMap{
		"lastItem": func(index int, values []EnumValue) bool {
			return len(values)-1 == index
		},
		"escapeVar":      escapeVarName,
		"formatScalar":   formatScalar,
		"lowerCamelCase": tools.LowerCamelCase,
		"formatType": func(_ ast.Type) string {
			panic("formatType() needs to be overridden by a jenny")
		},
		"isBuilder": func(name string, def ast.Type) bool {
			return def.Kind == ast.KindComposableSlot || def.Kind == ast.KindRef
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

	GenGettersAndSetters  bool
	GenBuilderConstructor bool
	Variant               string
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

type BuilderTemplate struct {
	Package         string
	Imports         fmt.Stringer
	Name            string
	ObjectSignature string
	Fields          []Field
	Options         []Option
}

type Option struct {
	Name        string
	Args        []ast.Argument
	Assignments []Assignment
	Type        any
}

type Assignment struct {
	Path           ast.Path
	InitSafeguards []string
	Constraints    []Constraint
	Method         ast.AssignmentMethod
	Value          ast.AssignmentValue
}

type Constraint struct {
	ArgName   string
	Op        ast.Op
	Parameter any
}
