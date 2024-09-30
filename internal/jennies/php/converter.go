package php

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Converter struct {
	config         Config
	nullableConfig languages.NullableConfig
	tmpl           *template.Template
}

func (jenny *Converter) JennyName() string {
	return "PHPConverter"
}

func (jenny *Converter) Generate(context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range context.Builders {
		output, err := jenny.generateConverter(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			"src",
			formatPackageName(builder.Package),
			fmt.Sprintf("%sConverter.php", formatObjectName(builder.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Converter) generateConverter(context languages.Context, builder ast.Builder) ([]byte, error) {
	converter := languages.NewConverterGenerator(jenny.nullableConfig).FromBuilder(context, builder)
	formatter := builderTypeFormatter(jenny.config, context)

	return jenny.tmpl.
		Funcs(map[string]any{
			"formatType": builderTypeFormatter(jenny.config, context).formatType,
			"formatPath": formatFieldPath,
			"formatRawRef": func(pkg string, ref string) string {
				return formatter.formatRef(ast.NewRef(pkg, ref), false)
			},
		}).
		RenderAsBytes("converters/converter.tmpl", map[string]any{
			"NamespaceRoot": jenny.config.NamespaceRoot,
			"Converter":     converter,
		})
}
