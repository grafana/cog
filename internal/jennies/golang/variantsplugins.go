package golang

import (
	"sort"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type VariantsPlugins struct {
	Tmpl   *template.Template
	Config Config
}

func (jenny VariantsPlugins) JennyName() string {
	return "GoVariantsPlugins"
}

func (jenny VariantsPlugins) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	registries, err := jenny.variantPlugins(context)
	if err != nil {
		return nil, err
	}

	models, err := jenny.variantModels()
	if err != nil {
		return nil, err
	}

	files = append(files, *codejen.NewFile("cog/plugins/variants.go", registries, jenny))
	files = append(files, *codejen.NewFile("cog/variants/variants.go", models, jenny))

	return files, nil
}

func (jenny VariantsPlugins) variantModels() ([]byte, error) {
	return jenny.Tmpl.RenderAsBytes("runtime/variant_models.tmpl", map[string]any{})
}

func (jenny VariantsPlugins) variantPlugins(context languages.Context) ([]byte, error) {
	imports := NewImportMap(jenny.Config.PackageRoot)
	var panelSchemas []*ast.Schema
	var dataquerySchemas []*ast.Schema

	imports.Add("cog", jenny.Config.importPath("cog"))

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Identifier == "" {
			continue
		}

		if schema.Metadata.Variant == ast.SchemaVariantPanel {
			panelSchemas = append(panelSchemas, schema)
		} else if schema.Metadata.Variant == ast.SchemaVariantDataQuery {
			dataquerySchemas = append(dataquerySchemas, schema)
		}

		imports.Add(schema.Package, jenny.Config.importPath(formatPackageName(schema.Package)))
	}

	// to guarantee a consistent output for this jenny
	sort.SliceStable(panelSchemas, func(i, j int) bool {
		return panelSchemas[i].Package < panelSchemas[j].Package
	})
	sort.SliceStable(dataquerySchemas, func(i, j int) bool {
		return dataquerySchemas[i].Package < dataquerySchemas[j].Package
	})

	return jenny.Tmpl.RenderAsBytes("runtime/variant_plugins.tmpl", map[string]any{
		"panel_schemas":     panelSchemas,
		"dataquery_schemas": dataquerySchemas,
		"imports":           imports,
	})
}
