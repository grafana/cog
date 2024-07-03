package typescript

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Converter struct {
	tmpl *template.Template
}

func (jenny *Converter) JennyName() string {
	return "TypescriptConverter"
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
			fmt.Sprintf("%sConverter.gen.ts", tools.LowerCamelCase(builder.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Converter) generateConverter(context languages.Context, builder ast.Builder) ([]byte, error) {
	converter := languages.NewConverterGenerator().FromBuilder(context, builder)

	imports := NewImportMap()
	typeImportMapper := func(pkg string) string {
		return imports.Add(pkg, fmt.Sprintf("../%s", pkg))
	}
	typeImportMapper("cog")

	typeFormatter := builderTypeFormatter(context, typeImportMapper)

	formatRawRef := func(pkg string, ref string) string {
		return typeFormatter.doFormatType(ast.NewRef(pkg, ref), false)
	}

	return jenny.tmpl.
		Funcs(map[string]any{
			"formatRawRef": formatRawRef,
			"formatValue":  formatValue,
		}).
		RenderAsBytes("converters/converter.tmpl", map[string]any{
			"Imports":   imports,
			"Converter": converter,
		})
}
