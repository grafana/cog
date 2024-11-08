package java

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

//go:embed templates/runtime/*.tmpl templates/types/*.tmpl templates/marshalling/*.tmpl templates/gradle/*.* templates/converters/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	tmpl, err := template.New(
		"java",
		template.Funcs(common.TypeResolvingTemplateHelpers(languages.Context{})),
		template.Funcs(functions()),
		template.ParseFS(templatesFS, "templates"),
		template.ParseDirectories(extraTemplatesDirectories...),
	)
	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}

	return tmpl
}

func functions() template.FuncMap {
	return template.FuncMap{
		"escapeVar":             escapeVarName,
		"formatScalar":          formatScalar,
		"cleanString":           cleanString,
		"lastPathIdentifier":    lastPathIdentifier,
		"fillAnnotationPattern": fillAnnotationPattern,
		"containsValue":         containsValue,
		"getJavaFieldTypeCheck": getJavaFieldTypeCheck,
		"isScalarType": func(p ast.Path, scalar string) bool {
			if !p.Last().Type.IsScalar() {
				return false
			}

			return p.Last().Type.AsScalar().ScalarKind == ast.ScalarKind(scalar)
		},
		"is": is,
		"importStdPkg": func(_ ast.Type) string {
			panic("formatRawRef() needs to be overridden by a jenny")
		},
		"formatPackageName": func(_ ast.Type) string {
			panic("formatRawRef() needs to be overridden by a jenny")
		},
		"formatRawRef": func(_ ast.Type) string {
			panic("formatRawRef() needs to be overridden by a jenny")
		},
		"fillNullableAnnotationPattern": func(_ ast.Type) string {
			panic("fillNullableAnnotationPattern() needs to be overridden by a jenny")
		},
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
		"formatRefType": func(_ ast.Type, value any) string {
			panic("formatRefType() needs to be overridden by a jenny")
		},
		"formatGuardPath": func(_ ast.Path) string {
			panic("formatGuardPath() needs to be overridden by a jenny")
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

	Fields     []ast.StructField
	Builders   []Builder
	HasBuilder bool

	Variant                 string
	Identifier              string
	Annotation              string
	ToJSONFunction          string
	ShouldAddSerializer     bool
	ShouldAddDeserializer   bool
	ShouldAddFactoryMethods bool
	DefaultConstructorArgs  []ast.Argument
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
	IsGeneric            bool
}

type OptionCall struct {
	Initializers []string
	OptionName   string
	Args         []string
}

type DataquerySchema struct {
	Identifier string
	Class      string
	Converter  string
}

type PanelSchema struct {
	Identifier  string
	Options     string
	FieldConfig string
	Converter   string
}

type Unmarshalling struct {
	Package                   string
	Name                      string
	ShouldUnmarshallingPanels bool
	Imports                   []string
	DataqueryUnmarshalling    []DataqueryUnmarshalling
	Fields                    []ast.StructField
	Hint                      any
}

type DataqueryUnmarshalling struct {
	DataqueryHint   string
	IsArray         bool
	DatasourceField string
	FieldName       string
}
