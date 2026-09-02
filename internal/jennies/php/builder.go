package php

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/template"
)

type Builder struct {
	config          Config
	tmpl            *template.Template
	typeFormatter   *typeFormatter
	apiRefCollector *common.APIReferenceCollector
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
	visitor := ir.BuilderVisitor{
		OnOption: func(visitor *ir.BuilderVisitor, schemas ir.Schemas, builder ir.Builder, option ir.Option) (ir.Option, error) {
			option.Args = tools.Map(option.Args, func(arg ir.Argument) ir.Argument {
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
		OnFactory: func(visitor *ir.BuilderVisitor, schemas ir.Schemas, builder ir.Builder, factory ir.BuilderFactory) (ir.BuilderFactory, error) {
			factory.Args = tools.Map(factory.Args, func(arg ir.Argument) ir.Argument {
				newArg := arg.DeepCopy()
				newArg.Type.Nullable = false

				if !hinter.requiresHint(newArg.Type) {
					return newArg
				}

				typehint := hinter.paramAnnotationForType(newArg.Name, newArg.Type)
				if typehint != "" {
					factory.Comments = append(factory.Comments, typehint)
				}

				return newArg
			})

			return factory, nil
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

func (jenny *Builder) generateBuilder(context languages.Context, builder ir.Builder) ([]byte, error) {
	builder.For.Comments = append(
		builder.For.Comments,
		fmt.Sprintf("@implements %s<%s>", jenny.config.fullNamespaceRef("Cog\\Builder"), jenny.typeFormatter.doFormatType(builder.For.SelfRef.AsType(), false)),
	)

	if builder.DeprecationMessage != "" {
		builder.For.Comments = append(builder.For.Comments, "@deprecated "+builder.DeprecationMessage)
	}

	jenny.apiRefCollector.BuilderMethod(builder, common.MethodReference{
		Name: "build",
		Comments: []string{
			"Builds the object.",
		},
		Return: jenny.typeFormatter.doFormatType(builder.For.SelfRef.AsType(), false),
	})

	return jenny.tmpl.
		Funcs(templateHelpers(templateDeps{
			config:  jenny.config,
			context: context,
		})).
		RenderAsBytes("builders/builder.tmpl", map[string]any{
			"NamespaceRoot": jenny.config.NamespaceRoot,
			"Builder":       builder,
		})
}
