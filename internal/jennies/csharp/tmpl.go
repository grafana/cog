package csharp

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

//go:embed templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(config Config, apiRefCollector *common.APIReferenceCollector) *template.Template {
	tmpl, err := template.New(
		"csharp",
		template.Funcs(common.TypeResolvingTemplateHelpers(languages.Context{})),
		template.Funcs(common.TypesTemplateHelpers(languages.Context{})),
		template.Funcs(common.APIRefTemplateHelpers(apiRefCollector)),
		template.Funcs(config.OverridesTemplateFuncs),
		template.Funcs(formattingTemplateFuncs()),

		template.ParseFS(templatesFS, "templates"),
		template.ParseFS(config.OverridesTemplatesFS, "custom"),
		template.ParseDirectories(config.OverridesTemplatesDirectories...),
	)
	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}

	return tmpl
}

// formattingTemplateFuncs returns helpers that are pure formatting
// utilities — they don't depend on a typeFormatter and so can be added
// to the template root without an active jenny.
func formattingTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatFieldName":   formatFieldName,
		"formatArgName":     formatArgName,
		"formatPackageName": formatPackageName,
		"formatObjectName":  formatObjectName,
		"escapeVar":         escapeVarName,
		"lastValueIndex": func(index int, values []EnumValue) bool {
			return len(values)-1 == index
		},
	}
}

// ---------------------------------------------------------------------
// Template data structures
// ---------------------------------------------------------------------

// classTemplate is the model bound to templates/types/class.tmpl.
type classTemplate struct {
	Namespace          string
	Name               string
	Imports            fmt.Stringer
	Comments           []string
	Extends            []string
	Fields             []classField
	DefaultAssignments []assignment
	HasArgsConstructor bool
	Args               []classArg
	ArgAssignments     []assignment
}

// classField mirrors ast.StructField but pre-formats the name and type
// so the template stays simple and so type formatting (which mutates
// the importMap as a side-effect) happens before the imports block is
// rendered.
type classField struct {
	Name     string
	Type     string
	Comments []string
}

type classArg struct {
	Name string
	Type string
}

type assignment struct {
	Name         string // already-formatted field name
	Value        string // literal expression, used by the parameterless ctor
	ValueFromArg string // already-formatted arg name, used by the all-args ctor
}

// enumTemplate is the model bound to templates/types/enum.tmpl.
type enumTemplate struct {
	Namespace       string
	Name            string
	Comments        []string
	IsString        bool
	NeedsEnumMember bool
	Values          []EnumValue
}

// EnumValue is exported so it can be used by the lastValueIndex helper.
type EnumValue struct {
	// Name is the C# member identifier (PascalCase, escaped).
	Name string
	// RawValue is the underlying schema value:
	//   - for integer enums: the int64 itself, formatted as `{n}` in
	//     templates;
	//   - for string enums: the original string, used inside
	//     [EnumMember(Value = "…")].
	RawValue any
}

// constantsTemplate is the model bound to templates/types/constants.tmpl.
type constantsTemplate struct {
	Namespace string
	Name      string
	Constants []constant
}

type constant struct {
	Name  string
	Type  string
	Value string
}
