package java

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type DataqueryConverter struct {
	config         Config
	nullableConfig languages.NullableConfig
	tmpl           *template.Template
}

func (jenny *DataqueryConverter) JennyName() string {
	return "JavaDataqueryConverter"
}

func (jenny *DataqueryConverter) Generate(context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}

	var err error
	for _, schema := range context.Schemas {
		schema.Objects = schema.Objects.Filter(func(key string, obj ast.Object) bool {
			if obj.Type.ImplementedVariant() != string(ast.SchemaVariantDataQuery) {
				return false
			}

			return !obj.Type.HasHint(ast.HintSkipVariantPluginRegistration)
		})

		schema.Objects.Iterate(func(key string, obj ast.Object) {
			output, genErr := jenny.generateSchema(context, schema, obj)
			if genErr != nil {
				err = genErr
			} else {
				filename := filepath.Join(
					jenny.config.ProjectPath,
					formatPackageName(schema.Package),
					fmt.Sprintf("%sMapperConverter.java", tools.UpperCamelCase(obj.Name)),
				)

				files = append(files, *codejen.NewFile(filename, output, jenny))
			}
		})
	}

	return files, err
}

func (jenny *DataqueryConverter) generateSchema(context languages.Context, schema *ast.Schema, obj ast.Object) ([]byte, error) {
	var disjunctionStruct *ast.StructType

	if obj.Type.IsRef() {
		resolved, _ := schema.Resolve(obj.Type)
		if resolved.IsStructGeneratedFromDisjunction() {
			disjunctionStruct = resolved.Struct
		}
	}

	imports := NewImportMap(jenny.config.PackagePath)
	packageMapper := func(pkg string, class string) string {
		if imports.IsIdentical(pkg, schema.Package) {
			return ""
		}

		return imports.Add(class, pkg)
	}

	typeFormatter := createFormatter(context, jenny.config).withPackageMapper(packageMapper)

	imports.Add("Converter", "cog")

	return jenny.tmpl.Funcs(map[string]any{
		"formatPackageName": typeFormatter.formatPackage,
	}).
		RenderAsBytes("converters/dataquery_converter.tmpl", map[string]any{
			"Package":     obj.SelfRef.ReferredPkg,
			"Imports":     imports.String(),
			"Name":        obj.Name,
			"Input":       jenny.formatPackage("cog.variants.Dataquery"),
			"Disjunction": disjunctionStruct,
		})
}

func (jenny *DataqueryConverter) formatPackage(pkg string) string {
	if jenny.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", jenny.config.PackagePath, pkg)
	}

	return pkg
}
