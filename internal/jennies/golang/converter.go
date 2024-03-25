package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
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
	typeFormatter := builderTypeFormatter(jenny.Config, context, typeImportMapper)

	formatFieldPath := func(fieldPath ast.Path) string {
		parts := make([]string, len(fieldPath))

		for i := range fieldPath {
			output := tools.UpperCamelCase(fieldPath[i].Identifier)

			// don't generate type hints if:
			// * there isn't one defined
			// * the type isn't "any"
			// * as a trailing element in the path
			if !fieldPath[i].Type.IsAny() || fieldPath[i].TypeHint == nil || i == len(fieldPath)-1 {
				parts[i] = output
				continue
			}

			formattedTypeHint := typeFormatter.formatType(*fieldPath[i].TypeHint)
			parts[i] = output + fmt.Sprintf(".(*%s)", formattedTypeHint)
		}

		return strings.Join(parts, ".")
	}

	err := templates.
		Funcs(map[string]any{
			"formatPath": formatFieldPath,
		}).
		ExecuteTemplate(&buffer, "converters/converter.tmpl", map[string]any{
			"Converter": converter,
		})
	if err != nil {
		return nil, err
	}

	return []byte(buffer.String()), nil
}
