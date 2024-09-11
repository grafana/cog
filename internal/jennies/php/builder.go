package php

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Builder struct {
	config        Config
	tmpl          *template.Template
	typeFormatter *typeFormatter
}

func (jenny *Builder) JennyName() string {
	return "PHPBuilder"
}

func (jenny *Builder) Generate(context languages.Context) (codejen.Files, error) {
	var err error
	jenny.typeFormatter = builderTypeFormatter(jenny.config, context)

	files := make([]codejen.File, 0, len(context.Builders))

	// Add argument typehints and ensure arguments are not nullable
	hinter := &typehints{config: jenny.config, context: context, resolveBuilders: true}
	visitor := ast.BuilderVisitor{
		OnOption: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, option ast.Option) (ast.Option, error) {
			option.Args = tools.Map(option.Args, func(arg ast.Argument) ast.Argument {
				newArg := arg.DeepCopy()
				newArg.Type.Nullable = false

				if !hinter.requiresHint(newArg.Type) {
					return newArg
				}

				typehint := hinter.paramAnnotationForType(newArg.Name, newArg.Type)
				if typehint != "" {
					option.Comments = append(option.Comments, typehint)
				}

				return newArg
			})

			return option, nil
		},
	}
	context.Builders, err = visitor.Visit(context.Schemas, context.Builders)
	if err != nil {
		return nil, err
	}

	for _, builder := range context.Builders {
		output, err := jenny.generateBuilder(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			"src",
			formatPackageName(builder.Package),
			fmt.Sprintf("%sBuilder.php", formatObjectName(builder.Name)),
		)
		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Builder) generateBuilder(context languages.Context, builder ast.Builder) ([]byte, error) {
	builder.For.Comments = append(
		builder.For.Comments,
		fmt.Sprintf("@implements %s<%s>", jenny.config.fullNamespaceRef("Cog\\Builder"), jenny.typeFormatter.doFormatType(builder.For.SelfRef.AsType(), false)),
	)

	hinter := &typehints{config: jenny.config, context: context, resolveBuilders: false}

	return jenny.tmpl.
		Funcs(map[string]any{
			"fullNamespaceRef": jenny.config.fullNamespaceRef,
			"formatPath":       jenny.formatFieldPath,
			"formatType":       jenny.typeFormatter.formatType,
			"formatRawType": func(def ast.Type) string {
				return jenny.typeFormatter.doFormatType(def, false)
			},
			"formatRawTypeNotNullable": func(def ast.Type) string {
				typeDef := def.DeepCopy()
				typeDef.Nullable = false

				return jenny.typeFormatter.doFormatType(typeDef, false)
			},
			"typeHasBuilder":          context.ResolveToBuilder,
			"isDisjunctionOfBuilders": context.IsDisjunctionOfBuilders,
			"typeHint": func(def ast.Type) string {
				clone := def.DeepCopy()
				clone.Nullable = false

				return hinter.forType(clone, false)
			},
			"resolvesToComposableSlot": func(typeDef ast.Type) bool {
				_, found := context.ResolveToComposableSlot(typeDef)
				return found
			},
			"formatValue": func(destinationType ast.Type, value any) string {
				if destinationType.IsRef() {
					referredObj, found := context.LocateObject(destinationType.AsRef().ReferredPkg, destinationType.AsRef().ReferredType)
					if found && referredObj.Type.IsEnum() {
						return jenny.typeFormatter.formatEnumValue(referredObj, value)
					}
				}

				return formatValue(value)
			},
			"defaultForType": func(typeDef ast.Type) string {
				return formatValue(defaultValueForType(jenny.config, context.Schemas, typeDef, nil))
			},
		}).
		RenderAsBytes("builders/builder.tmpl", map[string]any{
			"NamespaceRoot": jenny.config.NamespaceRoot,
			"Builder":       builder,
		})
}

func (jenny *Builder) formatFieldPath(fieldPath ast.Path) string {
	return strings.Join(tools.Map(fieldPath, func(item ast.PathItem) string {
		return formatFieldName(item.Identifier)
	}), "->")
}
