package java

import (
	"embed"
	"fmt"
	"text/template"

	"github.com/grafana/cog/internal/ast"
	cogtemplate "github.com/grafana/cog/internal/jennies/template"
)

//nolint:gochecknoglobals
var templates *template.Template

//go:embed templates/runtime/*.tmpl templates/types/*.tmpl templates/veneers/*.tmpl templates/marshalling/*.tmpl templates/gradle/*.*
//nolint:gochecknoglobals
var templatesFS embed.FS

//nolint:gochecknoinits
func init() {
	base := template.New("golang")
	base.
		Option("missingkey=error").
		Funcs(cogtemplate.Helpers(base)).
		Funcs(functions())

	templates = template.Must(cogtemplate.FindAndParseTemplates(templatesFS, base, "templates"))
}

func functions() template.FuncMap {
	return template.FuncMap{
		"escapeVar":             escapeVarName,
		"formatScalar":          formatScalar,
		"lastPathIdentifier":    lastPathIdentifier,
		"fillAnnotationPattern": fillAnnotationPattern,
		"lastItem": func(index int, values []EnumValue) bool {
			return len(values)-1 == index
		},
		"formatValue": func(_ ast.Type) string {
			panic("formatValue() needs to be overridden by a jenny")
		},
		"formatCastValue": func(_ ast.Type) string {
			panic("formatCastValue() needs to be overridden by a jenny")
		},
		"shouldCastNilCheck": func(_ ast.Type) string {
			panic("shouldCastNilCheck() needs to be overridden by a jenny")
		},
		"formatPath": func(_ ast.Type) string {
			panic("formatPath() needs to be overridden by a jenny")
		},
		"formatAssignmentPath": func(_ ast.Type) string {
			panic("formatAssignmentPath() needs to be overridden by a jenny")
		},
		"formatBuilderFieldType": func(_ ast.Type) string {
			panic("formatBuilderFieldType() needs to be overridden by a jenny")
		},
		"formatType": func(_ ast.Type) string {
			panic("formatType() needs to be overridden by a jenny")
		},
		"typeHasBuilder": func(_ ast.Type) bool {
			panic("typeHasBuilder() needs to be overridden by a jenny")
		},
		"emptyValueForType": func(_ ast.Type) string {
			panic("emptyValueForType() needs to be overridden by a jenny")
		},
		"resolvesToComposableSlot": func(_ ast.Type) bool {
			panic("resolvesToComposableSlot() needs to be overridden by a jenny")
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

	Fields     []Field
	Builders   []Builder
	HasBuilder bool

	Variant               string
	Annotation            string
	ToJSONFunction        string
	ShouldAddDeserialiser bool
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

type Builder struct {
	Package              string
	BuilderSignatureType string
	BuilderName          string
	ObjectName           string
	Imports              fmt.Stringer
	ImportAlias          string // alias to the pkg in which the object being built lives.
	Comments             []string
	Constructor          ast.Constructor
	Properties           []ast.StructField
	Options              []ast.Option
	Defaults             []OptionCall
}

type OptionCall struct {
	Initializers []string
	OptionName   string
	Args         []string
}

type DataquerySchema struct {
	Identifier string
	Class      string
}

type PanelSchema struct {
	Identifier  string
	Options     string
	FieldConfig string
}

type Unmarshalling struct {
	Package                   string
	Name                      string
	ShouldUnmarshallingPanels bool
	Imports                   []string
	DataqueryUnmarshalling    DataqueryUnmarshalling
	Fields                    []ast.StructField
	Hint                      any
}

type DataqueryUnmarshalling struct {
	DataqueryHint   string
	IsArray         bool
	DatasourceField string
	FieldName       string
}
