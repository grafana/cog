package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type Converter struct {
	Config Config
}

func (jenny *Converter) JennyName() string {
	return "GoConverter"
}

func (jenny *Converter) Generate(context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}

	for _, builder := range context.Builders {
		output, err := jenny.generateConverter(context, builder)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			formatPackageName(builder.Package),
			fmt.Sprintf("%s_converter_gen.go", strings.ToLower(builder.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Converter) generateConverter(context languages.Context, builder ast.Builder) ([]byte, error) {
	var buffer strings.Builder

	converter := (&languages.ConverterGenerator{}).FromBuilder(context, builder)

	imports := NewImportMap()
	typeImportMapper := func(pkg string) string {
		if imports.IsIdentical(pkg, builder.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.Config.importPath(pkg))
	}
	typeImportMapper("cog")
	typeFormatter := builderTypeFormatter(jenny.Config, context, typeImportMapper)

	formatRawRef := func(pkg string, ref string) string {
		return typeFormatter.formatRef(ast.NewRef(pkg, ref), false)
	}

	err := templates.
		Funcs(map[string]any{
			"typeHasEncoder": func(typeDef ast.Type) bool {
				return typeHasEncoder(context, typeDef)
			},
			"formatPath":   makePathFormatter(typeFormatter),
			"formatRawRef": formatRawRef,
		}).
		ExecuteTemplate(&buffer, "converters/converter.tmpl", map[string]any{
			"Imports":   imports,
			"Converter": converter,
		})
	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}
