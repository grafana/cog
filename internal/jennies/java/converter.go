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
			jenny.config.ProjectPath,
			formatPackageName(builder.Package),
			fmt.Sprintf("%sConverter.java", tools.UpperCamelCase(builder.Name)),
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	var err error
	for _, schema := range context.Schemas {
		schema.Objects = schema.Objects.Filter(func(key string, obj ast.Object) bool {
			if obj.Type.ImplementedVariant() != string(ast.SchemaVariantDataQuery) {
				return false
			}

			return !obj.Type.HasHint(ast.HintSkipVariantPluginRegistration)
		})

		schema.Objects.Iterate(func(key string, obj ast.Object) {
			output, genErr := jenny.generateDataqueryConverter(context, schema, obj)
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
	typeFormatter := createFormatter(context, jenny.config).withPackageMapper(packageMapper)

	return jenny.tmpl.
		Funcs(map[string]any{
			"formatRawRef": func(pkg string, ref string) string {
				return typeFormatter.formatReference(ast.NewRef(pkg, ref).AsRef())
			},
			"formatPath":        typeFormatter.formatFieldPath,
			"formatType":        typeFormatter.formatFieldType,
			"formatRefType":     typeFormatter.formatRefType,
			"formatGuardPath":   typeFormatter.formatGuardPath,
			"formatPackageName": typeFormatter.formatPackage,
			"importStdPkg":      packageMapper,
		}).
		RenderAsBytes("converters/converter.tmpl", map[string]any{
			"Imports":   imports,
			"Converter": converter,
			"IsPanel":   schemaFound && schema.Metadata.Variant == ast.SchemaVariantPanel && builder.Name == "Panel",
		})
}

func (jenny *Converter) generateDataqueryConverter(context languages.Context, schema *ast.Schema, obj ast.Object) ([]byte, error) {
	var disjunctionStruct *ast.StructType

	if obj.Type.IsDisjunctionOfRefs() {
		disjunctionStruct = obj.Type.Struct
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

func (jenny *Converter) formatPackage(pkg string) string {
	if jenny.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", jenny.config.PackagePath, pkg)
	}

	return pkg
}
