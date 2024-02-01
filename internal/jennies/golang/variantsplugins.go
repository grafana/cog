package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
)

type VariantsPlugins struct {
	Config Config
}

func (jenny VariantsPlugins) JennyName() string {
	return "GoVariantsPlugins"
}

func (jenny VariantsPlugins) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	registries, err := jenny.variantPlugins(context)
	if err != nil {
		return nil, err
	}

	models, err := jenny.variantModels()
	if err != nil {
		return nil, err
	}

	files = append(files, *codejen.NewFile("cog/plugins/variants.go", []byte(registries), jenny))
	files = append(files, *codejen.NewFile("cog/variants/variants.go", []byte(models), jenny))

	return files, nil
}

func (jenny VariantsPlugins) variantModels() (string, error) {
	return renderTemplate("runtime/variant_models.tmpl", map[string]any{})
}

func (jenny VariantsPlugins) variantPlugins(context common.Context) (string, error) {
	imports := NewImportMap()
	initMap := make(map[string][]*ast.Schema) // variant to schemas

	imports.Add("cog", jenny.Config.importPath("cog"))

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Identifier == "" {
			continue
		}

		variant := string(schema.Metadata.Variant)
		imports.Add(schema.Package, jenny.Config.importPath(formatPackageName(schema.Package)))
		initMap[variant] = append(initMap[variant], schema)
	}

	return renderTemplate("runtime/variant_plugins.tmpl", map[string]any{
		"init_map": initMap,
		"imports":  imports,
	})
}
