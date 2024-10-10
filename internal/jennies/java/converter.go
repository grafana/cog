package java

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
	return "JavaConverter"
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
			fmt.Sprintf("%sConverter.java", builder.Name),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny *Converter) generateConverter(context languages.Context, builder ast.Builder) ([]byte, error) {
	converter := languages.NewConverterGenerator(jenny.nullableConfig).FromBuilder(context, builder)

	schema, schemaFound := context.Schemas.Locate(builder.Package)

	imports := NewImportMap(jenny.config.PackagePath)
	packageMapper := func(pkg string, class string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(class, pkg)
	}

	inputIsDataquery := schemaFound && schema.Metadata.Variant == ast.SchemaVariantDataQuery && schema.EntryPoint == builder.For.Name
	typeFormatter := createFormatter(context, jenny.config).withPackageMapper(packageMapper)

	packageMapper("java.util", "List")
	packageMapper("java.util", "LinkedList")

	return jenny.tmpl.
		Funcs(map[string]any{
			"formatRawRef": func(pkg string, ref string) string {
				return typeFormatter.formatReference(ast.NewRef(pkg, ref).AsRef())
			},
			"formatPath": typeFormatter.formatFieldPath,
			"formatType": typeFormatter.formatFieldType,
		}).
		RenderAsBytes("converters/converter.tmpl", map[string]any{
			"Imports":          imports,
			"Converter":        converter,
			"InputIsDataquery": inputIsDataquery,
		})
}
