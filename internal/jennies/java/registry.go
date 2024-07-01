package java

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Registry struct {
	config Config
}

func (jenny Registry) JennyName() string {
	return "JavaRegistry"
}

func (jenny Registry) Generate(context languages.Context) (codejen.Files, error) {
	panelRegistry, err := jenny.renderPanelRegistry()
	if err != nil {
		return nil, err
	}

	registry, err := jenny.renderRegistry(context)
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/PanelRegistry.java"), panelRegistry, jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/Registry.java"), registry, jenny),
	}, nil
}

func (jenny Registry) renderPanelRegistry() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, "runtime/panel_registry.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
	}); err != nil {
		return nil, fmt.Errorf("failed executing template: %w", err)
	}

	return buf.Bytes(), nil
}

func (jenny Registry) renderRegistry(context languages.Context) ([]byte, error) {
	imports := NewImportMap(jenny.config.PackagePath)
	var panelSchemas []DataqueryPanelSchema
	var dataquerySchemas []DataqueryPanelSchema

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Identifier == "" {
			continue
		}

		if schema.Metadata.Variant == ast.SchemaVariantDataQuery {
			dataquerySchemas = append(dataquerySchemas, DataqueryPanelSchema{
				Identifier: strings.ToLower(schema.Metadata.Identifier),
				Class:      jenny.formatPackage(fmt.Sprintf("%s.%s", schema.Package, jenny.findDataqueryClass(schema))),
			})
		} else {
			panelSchemas = append(panelSchemas, DataqueryPanelSchema{
				Identifier: strings.ToLower(schema.Metadata.Identifier),
				Class:      jenny.formatPackage(fmt.Sprintf("%s.Options", schema.Package)),
			})
		}

	}

	sort.SliceStable(dataquerySchemas, func(i, j int) bool {
		return dataquerySchemas[i].Identifier < dataquerySchemas[j].Identifier
	})

	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, "runtime/registry.tmpl", map[string]any{
		"Package":          jenny.formatPackage("cog.variants"),
		"Imports":          imports,
		"PanelSchemas":     panelSchemas,
		"DataquerySchemas": dataquerySchemas,
	}); err != nil {
		return nil, fmt.Errorf("failed executing template: %w", err)
	}

	return buf.Bytes(), nil
}

func (jenny Registry) findDataqueryClass(schema *ast.Schema) string {
	name := ""
	schema.Objects.Iterate(func(key string, object ast.Object) {
		if object.Type.ImplementedVariant() == string(ast.SchemaVariantDataQuery) && !object.Type.HasHint(ast.HintSkipVariantPluginRegistration) {
			name = tools.UpperCamelCase(object.Name)
		}
	})

	return name
}

func (jenny Registry) formatPackage(pkg string) string {
	if jenny.config.PackagePath != "" {
		return fmt.Sprintf("%s.%s", jenny.config.PackagePath, pkg)
	}

	return pkg
}
