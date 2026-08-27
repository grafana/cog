package python

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
)

//go:embed templates/*/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(config Config, apiRefCollector *common.APIReferenceCollector) *template.Template {
	tmpl, err := template.New(
		"python",

		template.Funcs(common.TypeResolvingTemplateHelpers(languages.Context{})),
		template.Funcs(common.TypesTemplateHelpers(languages.Context{})),
		template.Funcs(common.APIRefTemplateHelpers(apiRefCollector)),
		template.Funcs(config.OverridesTemplateFuncs),
		template.Funcs(formattingTemplateFuncs()),
		// placeholder functions, will be overridden by jennies
		template.Funcs(template.FuncMap{
			"isDisjunctionOfBuilders": func(_ ir.Type) string {
				panic("isDisjunctionOfBuilders() needs to be overridden by a jenny")
			},
			"formatType": func(_ ir.Type) string {
				panic("formatType() needs to be overridden by a jenny")
			},
			"formatTypeNotNullable": func(_ ir.Type) string {
				panic("formatTypeNotNullable() needs to be overridden by a jenny")
			},
			"formatFullyQualifiedRef": func(_ ir.Type) string {
				panic("formatFullyQualifiedRef() needs to be overridden by a jenny")
			},
			"formatRawType": func(_ ir.Type) string {
				panic("formatRawType() needs to be overridden by a jenny")
			},
			"formatRawTypeNotNullable": func(_ ir.Type) string {
				panic("formatRawTypeNotNullable() needs to be overridden by a jenny")
			},
			"formatValue": func(_ ir.Type, _ any) string {
				panic("formatValue() needs to be overridden by a jenny")
			},
			"defaultForType": func(_ ir.Type) string {
				panic("defaultForType() needs to be overridden by a jenny")
			},
			"importModule": func(alias string, pkg string, module string) string {
				panic("importModule() needs to be overridden by a jenny")
			},
			"importPkg": func(alias string, pkg string) string {
				panic("importPkg() needs to be overridden by a jenny")
			},
			"unmarshalForType": func(typeDef ir.Type, inputVar string, hint string) fromJSONCode {
				panic("unmarshalForType() needs to be overridden by a jenny")
			},
		}),

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
		"refToType": func(ref ir.RefType) ir.Type {
			return ref.AsType()
		},
		"formatIdentifier":   formatIdentifier,
		"formatFunctionName": formatFunctionName,
		"formatPath":         formatFieldPath,
		"formatObjectName":   formatObjectName,
		"formatConcrete":     formatValue,
	}
}
