package golang

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Converter struct {
	Config          Config
	NullableConfig  languages.NullableConfig
	Tmpl            *template.Template
	apiRefCollector *common.APIReferenceCollector
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
	converter := languages.NewConverterGenerator(jenny.NullableConfig).FromBuilder(context, builder)

	imports := NewImportMap(jenny.Config.PackageRoot)
	typeImportMapper := func(pkg string) string {
		if imports.IsIdentical(pkg, builder.Package) {
			return ""
		}

		return imports.Add(pkg, jenny.Config.importPath(pkg))
	}
	formatter := builderTypeFormatter(jenny.Config, context, imports, typeImportMapper)

	dummyImports := NewImportMap(jenny.Config.PackageRoot)
	dummyImportMapper := func(pkg string) string {
		return dummyImports.Add(pkg, jenny.Config.importPath(pkg))
	}

	formatRawRef := func(pkg string, ref string) string {
		return formatter.formatRef(ast.NewRef(pkg, ref), false)
	}

	jenny.apiRefCollector.RegisterFunction(builder.Package, common.FunctionReference{
		Name: fmt.Sprintf("%sConverter", tools.UpperCamelCase(converter.BuilderName)),
		Arguments: []common.ArgumentReference{
			{
				Name: "input",
				Type: formatRawRef(builder.For.SelfRef.ReferredPkg, builder.For.SelfRef.ReferredType),
			},
		},
		Return: "string",
	})

	return jenny.Tmpl.
		Funcs(map[string]any{
			"importStdPkg": func(pkg string) string {
				return imports.Add(pkg, pkg)
			},
			"importPkg":           typeImportMapper,
			"formatType":          builderTypeFormatter(jenny.Config, context, dummyImports, dummyImportMapper).formatType,
			"formatTypeNoBuilder": defaultTypeFormatter(jenny.Config, context, dummyImports, dummyImportMapper).formatType,
			"formatPath":          makePathFormatter(formatter),
			"formatRawRef":        formatRawRef,
			"maybeUnptr": func(variableName string, intoType ast.Type) string {
				if !intoType.Nullable || intoType.IsArray() || intoType.IsMap() || intoType.IsComposableSlot() {
					return variableName
				}
				typeImportMapper("cog")
				return "cog.Unptr(" + variableName + ")"
			},
		}).
		RenderAsBytes("converters/converter.tmpl", map[string]any{
			"Imports":   imports,
			"Converter": converter,
		})
}
