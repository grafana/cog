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
	panelRegistry, err := jenny.renderPanelConfig()
	if err != nil {
		return nil, err
	}

	registry, err := jenny.renderRegistry(context)
	if err != nil {
		return nil, err
	}

	emptyDataquery, err := jenny.emptyDataquery()
	if err != nil {
		return nil, err
	}

	emptyDataquerySerializer, err := jenny.emptyDataquerySerializer()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/PanelConfig.java"), panelRegistry, jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/Registry.java"), registry, jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/EmptyDataquery.java"), emptyDataquery, jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/EmptyDataquerySerializer.java"), emptyDataquerySerializer, jenny),
	}, nil
}

func (jenny Registry) renderPanelConfig() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, "runtime/panel_config.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
	}); err != nil {
		return nil, fmt.Errorf("failed executing template: %w", err)
	}

	return buf.Bytes(), nil
}

func (jenny Registry) renderRegistry(context languages.Context) ([]byte, error) {
	imports := NewImportMap(jenny.config.PackagePath)
	var panelSchemas []PanelSchema
	var dataquerySchemas []DataquerySchema

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Identifier == "" {
			continue
		}

		if schema.Metadata.Variant == ast.SchemaVariantDataQuery {
			dataquerySchemas = append(dataquerySchemas, DataquerySchema{
				Identifier: strings.ToLower(schema.Metadata.Identifier),
				Class:      jenny.formatPackage(fmt.Sprintf("%s.%s", schema.Package, jenny.findDataqueryClass(schema))),
			})
		} else if schema.Metadata.Variant == ast.SchemaVariantPanel {
			panelSchemas = append(panelSchemas, PanelSchema{
				Identifier:  strings.ToLower(schema.Metadata.Identifier),
				Options:     jenny.formatPackage(fmt.Sprintf("%s.Options.class", schema.Package)),
				FieldConfig: jenny.findFieldConfig(schema),
			})
		}
	}

	sort.SliceStable(panelSchemas, func(i, j int) bool {
		return panelSchemas[i].Identifier < panelSchemas[j].Identifier
	})

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

func (jenny Registry) findFieldConfig(schema *ast.Schema) string {
	name := "null"
	schema.Objects.Iterate(func(key string, value ast.Object) {
		if key == "FieldConfig" {
			name = jenny.formatPackage(fmt.Sprintf("%s.FieldConfig.class", schema.Package))
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

func (jenny Registry) emptyDataquery() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, "runtime/empty_dataquery.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
	}); err != nil {
		return nil, fmt.Errorf("failed executing template: %w", err)
	}

	return buf.Bytes(), nil
}

func (jenny Registry) emptyDataquerySerializer() ([]byte, error) {
	buf := bytes.Buffer{}
	if err := templates.ExecuteTemplate(&buf, "runtime/empty_dataquery_serializer.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
	}); err != nil {
		return nil, fmt.Errorf("failed executing template: %w", err)
	}

	return buf.Bytes(), nil
}
