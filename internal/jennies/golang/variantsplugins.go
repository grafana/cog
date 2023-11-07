package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/context"
)

type VariantsPlugins struct {
}

func (jenny VariantsPlugins) JennyName() string {
	return "GoVariantsPlugins"
}

func (jenny VariantsPlugins) Generate(context context.Builders) (codejen.Files, error) {
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
	return renderTemplate("variant_models.tmpl", map[string]any{})
}

func (jenny VariantsPlugins) variantPlugins(context context.Builders) (string, error) {
	imports := newImportMap()
	initMap := make(map[string][]*ast.Schema) // variant to schemas

	imports.Add("cog", "github.com/grafana/cog/generated/cog")

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Identifier == "" {
			continue
		}

		variant := string(schema.Metadata.Variant)
		imports.Add(schema.Package, "github.com/grafana/cog/generated/"+schema.Package)
		initMap[variant] = append(initMap[variant], schema)
	}

	return renderTemplate("variant_plugins.tmpl", map[string]any{
		"init_map": initMap,
		"imports":  imports.Format(),
	})
}
