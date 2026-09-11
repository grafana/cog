package php

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

//go:embed templates/builders/*.tmpl templates/converters/*.tmpl templates/runtime/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(config Config, apiRefCollector *apiref.APIReferenceCollector) *template.Template {
	tmpl, err := template.New(
		"php",

		// "dummy"/unimplemented helpers, to be able to parse the templates before jennies are initialized.
		// Jennies will override these with proper dependencies.
		template.Funcs(template.TypesHelpers(languages.Context{})),
		template.Funcs(apiref.TemplateHelpers(apiRefCollector)),
		template.Funcs(common.DynamicFilesTemplateHelpers()),

		template.Funcs(templateHelpers(templateDeps{})),
		template.Funcs(formattingTemplateFuncs()),
		template.Funcs(config.OverridesTemplateFuncs),

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
		"formatPath":           formatFieldPath,
		"formatPackageName":    formatPackageName,
		"formatObjectName":     formatObjectName,
		"formatOptionName":     formatOptionName,
		"formatEnumMemberName": formatEnumMemberName,
		"formatArgName":        formatArgName,
		"formatFieldName":      formatFieldName,
		"formatScalar":         formatValue,
		"formatDocsBlock":      formatCommentsBlock,
	}
}

type templateDeps struct {
	config                   Config
	context                  languages.Context
	unmarshalForType         func(typeDef ir.Type, inputVar string) string
	unmarshalDisjunctionFunc func(typeDef ir.Type) string
	convertDisjunctionFunc   func(typeDef ir.Type) string
}

func templateHelpers(deps templateDeps) template.FuncMap {
	typesFormatter := builderTypeFormatter(deps.config, deps.context)
	hinter := &typehints{config: deps.config, context: deps.context, resolveBuilders: false}
	shaper := &shape{context: deps.context}

	funcs := template.FuncMap{
		"fullNamespace":           deps.config.fullNamespace,
		"fullNamespaceRef":        deps.config.fullNamespaceRef,
		"typeHasBuilder":          deps.context.ResolveToBuilder,
		"isDisjunctionOfBuilders": deps.context.IsDisjunctionOfBuilders,

		"formatType": typesFormatter.formatType,
		"formatRawType": func(def ir.Type) string {
			return typesFormatter.doFormatType(def, false)
		},
		"formatRawRef": func(pkg string, ref string) string {
			return typesFormatter.formatRef(ir.NewRef(pkg, ref), false)
		},
		"formatRawTypeNotNullable": func(def ir.Type) string {
			typeDef := def.DeepCopy()
			typeDef.Nullable = false

			return typesFormatter.doFormatType(typeDef, false)
		},
		"formatValue": func(destinationType ir.Type, value any) string {
			if destinationType.IsRef() {
				referredObj, found := deps.context.GetObjectByRef(destinationType.AsRef())
				if found && referredObj.Type.IsEnum() {
					return typesFormatter.formatEnumValue(referredObj, value)
				}
			}

			if destinationType.IsScalar() && (destinationType.Scalar.ScalarKind == ir.KindFloat32 || destinationType.Scalar.ScalarKind == ir.KindFloat64) {
				return fmt.Sprintf("(float) %s", formatValue(value))
			}

			return formatValue(value)
		},

		"typeHint": func(def ir.Type) string {
			clone := def.DeepCopy()
			clone.Nullable = false

			return hinter.forType(clone, false)
		},
		"typeShape": shaper.typeShape,
		"defaultForType": func(typeDef ir.Type) string {
			return formatValue(defaultValueForType(deps.config, deps.context.Schemas, typeDef, nil))
		},
		"disjunctionCaseForType": func(input string, typeDef ir.Type) string {
			return disjunctionCaseForType(typesFormatter, input, typeDef)
		},

		"factoryClassForPkg": func(pkg string) string {
			return deps.config.builderFactoryClassForPackage(pkg)
		},

		"unmarshalForType":         deps.unmarshalForType,
		"unmarshalDisjunctionFunc": deps.unmarshalDisjunctionFunc,
		"convertDisjunctionFunc":   deps.convertDisjunctionFunc,
	}

	return funcs.MergeWith(template.TypesHelpers(deps.context))
}
