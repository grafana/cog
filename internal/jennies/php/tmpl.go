package php

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

//go:embed templates/builders/*.tmpl templates/converters/*.tmpl templates/runtime/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	tmpl, err := template.New(
		"php",

		// "dummy"/unimplemented helpers, to be able to parse the templates before jennies are initialized.
		// Jennies will override these with proper dependencies.
		template.Funcs(templateHelpers(templateDeps{})),
		template.Funcs(map[string]any{
			"formatPath":           formatFieldPath,
			"formatPackageName":    formatPackageName,
			"formatObjectName":     formatObjectName,
			"formatOptionName":     formatOptionName,
			"formatEnumMemberName": formatEnumMemberName,
			"formatArgName":        formatArgName,
			"formatScalar":         formatValue,
			"formatDocsBlock":      formatCommentsBlock,
		}),

		// parse templates
		template.ParseFS(templatesFS, "templates"),
		template.ParseDirectories(extraTemplatesDirectories...),
	)
	if err != nil {
		panic(fmt.Errorf("could not initialize templates: %w", err))
	}

	return tmpl
}

type templateDeps struct {
	config           Config
	context          languages.Context
	unmarshalForType func(typeDef ast.Type, inputVar string) string
}

func templateHelpers(deps templateDeps) template.FuncMap {
	typesFormatter := builderTypeFormatter(deps.config, deps.context)
	hinter := &typehints{config: deps.config, context: deps.context, resolveBuilders: false}
	shaper := &shape{context: deps.context}

	return template.FuncMap{
		"fullNamespaceRef":        deps.config.fullNamespaceRef,
		"typeHasBuilder":          deps.context.ResolveToBuilder,
		"isDisjunctionOfBuilders": deps.context.IsDisjunctionOfBuilders,

		"formatType": typesFormatter.formatType,
		"formatRawType": func(def ast.Type) string {
			return typesFormatter.doFormatType(def, false)
		},
		"formatRawRef": func(pkg string, ref string) string {
			return typesFormatter.formatRef(ast.NewRef(pkg, ref), false)
		},
		"formatRawTypeNotNullable": func(def ast.Type) string {
			typeDef := def.DeepCopy()
			typeDef.Nullable = false

			return typesFormatter.doFormatType(typeDef, false)
		},
		"formatValue": func(destinationType ast.Type, value any) string {
			if destinationType.IsRef() {
				referredObj, found := deps.context.LocateObjectByRef(destinationType.AsRef())
				if found && referredObj.Type.IsEnum() {
					return typesFormatter.formatEnumValue(referredObj, value)
				}
			}

			return formatValue(value)
		},

		"typeHint": func(def ast.Type) string {
			clone := def.DeepCopy()
			clone.Nullable = false

			return hinter.forType(clone, false)
		},
		"typeShape": shaper.typeShape,
		"defaultForType": func(typeDef ast.Type) string {
			return formatValue(defaultValueForType(deps.config, deps.context.Schemas, typeDef, nil))
		},
		"disjunctionCaseForType": func(input string, typeDef ast.Type) string {
			return disjunctionCaseForType(typesFormatter, input, typeDef)
		},

		"unmarshalForType": deps.unmarshalForType,

		"resolveRefs": deps.context.ResolveRefs,
		"resolvesToStruct": func(typeDef ast.Type) bool {
			return deps.context.ResolveRefs(typeDef).IsStruct()
		},
		"resolvesToMap": func(typeDef ast.Type) bool {
			return deps.context.ResolveRefs(typeDef).IsMap()
		},
		"resolvesToEnum": func(typeDef ast.Type) bool {
			return deps.context.ResolveRefs(typeDef).IsEnum()
		},
		"resolvesToComposableSlot": func(typeDef ast.Type) bool {
			_, found := deps.context.ResolveToComposableSlot(typeDef)
			return found
		},
	}
}
