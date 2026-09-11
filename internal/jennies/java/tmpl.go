package java

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

//go:embed templates/runtime/*.tmpl templates/types/*.tmpl templates/marshalling/*.tmpl templates/converters/*.tmpl templates/builders/*.*
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(config Config, apiRefCollector *apiref.APIReferenceCollector) *template.Template {
	tmpl, err := template.New(
		"java",
		template.Funcs(template.TypesHelpers(languages.Context{})),
		template.Funcs(apiref.TemplateHelpers(apiRefCollector)),
		template.Funcs(config.OverridesTemplateFuncs),
		template.Funcs(functions()),
		template.Funcs(formattingTemplateFuncs()),

		// parse templates
		template.ParseFS(templatesFS, "templates"),
		template.ParseFS(config.OverridesTemplatesFS, "custom"),
		template.ParseDirectories(config.OverridesTemplatesDirectories...),
	)
	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}

	return tmpl
}

func formattingTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatPackageName": formatPackageName,
		"formatObjectName":  formatObjectName,
		"formatArgName":     formatArgName,
		"escapeVar":         escapeVarName,
		"formatScalar":      formatScalar,
		"cleanString":       cleanString,
		"formatIntegerLetter": func(t ir.Type) string {
			switch t.AsScalar().ScalarKind {
			case ir.KindInt64, ir.KindUint64:
				return "L"
			case ir.KindFloat32:
				return "f"
			}
			return ""
		},
	}
}

func functions() template.FuncMap {
	return template.FuncMap{
		"lastPathIdentifier":    lastPathIdentifier,
		"fillAnnotationPattern": fillAnnotationPattern,
		"containsValue":         containsValue,
		"getJavaFieldTypeCheck": getJavaFieldTypeCheck,
		"lastItem": func(index int, values []EnumValue) bool {
			return len(values)-1 == index
		},
		"refToType": func(ref ir.RefType) ir.Type {
			return ref.AsType()
		},
		"importStdPkg": func(_ ir.Type) string {
			panic("importStdPkg() needs to be overridden by a jenny")
		},
		"importPkg": func(_ string) string {
			panic("importPkg() needs to be overridden by a jenny")
		},
		"formatPackageName": func(_ ir.Type) string {
			panic("formatPackageName() needs to be overridden by a jenny")
		},
		"formatRawRef": func(_ ir.Type) string {
			panic("formatRawRef() needs to be overridden by a jenny")
		},
		"fillNullableAnnotationPattern": func(_ ir.Type) string {
			panic("fillNullableAnnotationPattern() needs to be overridden by a jenny")
		},
		"formatValue": func(_ ir.Type) string {
			panic("formatValue() needs to be overridden by a jenny")
		},
		"formatPathIndex": func(_ *ir.PathIndex) string {
			panic("formatPathIndex() needs to be overridden by a jenny")
		},
		"formatPath": func(_ ir.Type) string {
			panic("formatPath() needs to be overridden by a jenny")
		},
		"formatAssignmentPath": func(_ ir.Type) string {
			panic("formatAssignmentPath() needs to be overridden by a jenny")
		},
		"formatBuilderFieldType": func(_ ir.Type) string {
			panic("formatBuilderFieldType() needs to be overridden by a jenny")
		},
		"formatType": func(_ ir.Type) string {
			panic("formatType() needs to be overridden by a jenny")
		},
		"typeHasBuilder": func(_ ir.Type) bool {
			panic("typeHasBuilder() needs to be overridden by a jenny")
		},
		"emptyValueForType": func(_ ir.Type) string {
			panic("emptyValueForType() needs to be overridden by a jenny")
		},
		"resolvesToComposableSlot": func(_ ir.Type) bool {
			panic("resolvesToComposableSlot() needs to be overridden by a jenny")
		},
		"formatRefType": func(_ ir.Type, value any) string {
			panic("formatRefType() needs to be overridden by a jenny")
		},
		"formatGuardPath": func(_ ir.Path) string {
			panic("formatGuardPath() needs to be overridden by a jenny")
		},
		"enumFromConstantRef": func(_ ir.Type) string {
			panic("enumFromConstantRef() needs to be overridden by a jenny")
		},
		"factoryClassForPkg": func(_ string) string {
			panic("factoryClassForPkg() needs to be overridden by a jenny")
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
	Package            string
	RawPackage         string
	Imports            fmt.Stringer
	Name               string
	Extends            []string
	Comments           []string
	DeprecationMessage string

	Fields     []ir.StructField
	Builders   []BuilderTemplate
	HasBuilder bool

	Variant                 string
	Identifier              string
	Annotation              string
	ToJSONFunction          string
	ShouldAddSerializer     bool
	ShouldAddDeserializer   bool
	ShouldAddFactoryMethods bool
	Constructors            []ConstructorTemplate
	ExtraFunctionsBlock     string
}

type ConstructorTemplate struct {
	Args        []ir.Argument
	Assignments []ConstructorAssignmentTemplate
}

type ConstructorAssignmentTemplate struct {
	Name         string
	Type         ir.Type
	Value        any
	ValueFromArg string
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
	Package              string
	RawPackage           string
	BuilderSignatureType string
	BuilderName          string
	ObjectName           string
	Imports              fmt.Stringer
	ImportAlias          string // alias to the pkg in which the object being built lives.
	Comments             []string
	DeprecationMessage   string
	Constructor          ir.Constructor
	Properties           []ir.StructField
	Options              []ir.Option
	IsGenericPanel       bool
}

type Default struct {
	Name  string
	Value string
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
	Imports                   fmt.Stringer
	DataqueryUnmarshalling    []DataqueryUnmarshalling
	Fields                    []ir.StructField
	Hint                      any
}

type DataqueryUnmarshalling struct {
	DataqueryHint   string
	IsArray         bool
	DatasourceField string
	FieldName       string
}
