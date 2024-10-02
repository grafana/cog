package php

import (
	"embed"
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
)

//go:embed templates/builders/*.tmpl templates/converters/*.tmpl templates/runtime/*.tmpl templates/types/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

func initTemplates(extraTemplatesDirectories []string) *template.Template {
	tmpl, err := template.New(
		"php",

		// placeholder functions, will be overridden by jennies
		template.Funcs(template.FuncMap{
			"fullNamespaceRef": func(_ string) string {
				panic("fullNamespaceRef() needs to be overridden by a jenny")
			},
			"formatType":    func(_ ast.Type) string { panic("formatType() needs to be overridden by a jenny") },
			"formatRawType": func(_ ast.Type) string { panic("formatRawType() needs to be overridden by a jenny") },
			"formatRawRef": func(pkg string, ref string) string {
				panic("formatRawRef() needs to be overridden by a jenny")
			},
			"formatRawTypeNotNullable": func(_ ast.Type) string { panic("formatRawTypeNotNullable() needs to be overridden by a jenny") },
			"typeHasBuilder": func(_ ast.Type) bool {
				panic("typeHasBuilder() needs to be overridden by a jenny")
			},
			"typeHint": func(_ ast.Type) string {
				panic("typeHint() needs to be overridden by a jenny")
			},
			"isDisjunctionOfBuilders": func(_ ast.Type) bool {
				panic("isDisjunctionOfBuilders() needs to be overridden by a jenny")
			},
			"resolvesToComposableSlot": func(_ ast.Type) bool {
				panic("resolvesToComposableSlot() needs to be overridden by a jenny")
			},
			"defaultForType": func(_ ast.Type) bool {
				panic("defaultForType() needs to be overridden by a jenny")
			},
			"formatValue": func(_ ast.Type, _ any) bool {
				panic("formatValue() needs to be overridden by a jenny")
			},
			"disjunctionCaseForType": func(input string, typeDef ast.Type) string {
				panic("disjunctionCaseForType() needs to be overridden by a jenny")
			},
			"resolvesToEnum": func(typeDef ast.Type) bool {
				panic("refResolvesToEnum() needs to be overridden by a jenny")
			},
			"resolvesToStruct": func(typeDef ast.Type) bool {
				panic("resolvesToStruct() needs to be overridden by a jenny")
			},
			"resolvesToMap": func(typeDef ast.Type) bool {
				panic("resolvesToMap() needs to be overridden by a jenny")
			},
			"resolveRefs": func(typeDef ast.Type) ast.Type {
				panic("resolveRefs() needs to be overridden by a jenny")
			},
		}),
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
