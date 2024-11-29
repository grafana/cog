package java

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/tools"
)

type Registry struct {
	config Config
	tmpl   *template.Template
}

func (jenny Registry) JennyName() string {
	return "JavaRegistry"
}

func (jenny Registry) Generate(context languages.Context) (codejen.Files, error) {
	panelConfig, err := jenny.renderPanelConfig()
	if err != nil {
		return nil, err
	}

	dataqueryConfig, err := jenny.renderDataqueryConfig()
	if err != nil {
		return nil, err
	}

	registry, err := jenny.renderRegistry(context)
	if err != nil {
		return nil, err
	}

	unknownDataquery, err := jenny.unknownDataquery()
	if err != nil {
		return nil, err
	}

	unknownDataquerySerializer, err := jenny.unknownDataquerySerializer()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/PanelConfig.java"), []byte(panelConfig), jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/DataqueryConfig.java"), []byte(dataqueryConfig), jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/Registry.java"), []byte(registry), jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/UnknownDataquery.java"), []byte(unknownDataquery), jenny),
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "cog/variants/UnknownDataquerySerializer.java"), []byte(unknownDataquerySerializer), jenny),
	}, nil
}

func (jenny Registry) renderPanelConfig() (string, error) {
	imports := NewImportMap(jenny.config.PackagePath)
	if jenny.config.generateConverters {
		imports.Add("Converter", "cog")
		imports.Add("Panel", "dashboard")
	}

	return jenny.tmpl.Render("runtime/panel_config.tmpl", map[string]any{
		"Package":             jenny.formatPackage("cog.variants"),
		"Imports":             imports.String(),
		"ShouldAddConverters": jenny.config.generateConverters,
	})
}

func (jenny Registry) renderDataqueryConfig() (string, error) {
	imports := NewImportMap(jenny.config.PackagePath)
	if jenny.config.generateConverters {
		imports.Add("Converter", "cog")
	}

	return jenny.tmpl.Render("runtime/dataquery_config.tmpl", map[string]any{
		"Package":             jenny.formatPackage("cog.variants"),
		"Imports":             imports.String(),
		"ShouldAddConverters": jenny.config.generateConverters,
	})
}

func (jenny Registry) renderRegistry(context languages.Context) (string, error) {
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
				Class:      jenny.formatPackage(fmt.Sprintf("%s.%s.class", schema.Package, jenny.findDataqueryClass(schema))),
				Converter:  jenny.formatPackage(fmt.Sprintf("%s.%sMapperConverter()", schema.Package, jenny.findDataqueryClass(schema))),
			})
		} else if schema.Metadata.Variant == ast.SchemaVariantPanel {
			panelSchemas = append(panelSchemas, PanelSchema{
				Identifier:  strings.ToLower(schema.Metadata.Identifier),
				Options:     jenny.formatPackage(fmt.Sprintf("%s.Options.class", schema.Package)),
				FieldConfig: jenny.findFieldConfig(schema),
				Converter:   jenny.formatPackage(fmt.Sprintf("%s.PanelConverter()", schema.Package)),
			})
		}
	}

	sort.SliceStable(panelSchemas, func(i, j int) bool {
		return panelSchemas[i].Identifier < panelSchemas[j].Identifier
	})

	if jenny.config.generateConverters {
		imports.Add("Panel", "dashboard")
		imports.Add("reflect.InvocationTargetException", "java.lang")
		imports.Add("Converter", "cog")
	}

	sort.SliceStable(dataquerySchemas, func(i, j int) bool {
		return dataquerySchemas[i].Identifier < dataquerySchemas[j].Identifier
	})

	return jenny.tmpl.Render("runtime/registry.tmpl", map[string]any{
		"Package":             jenny.formatPackage("cog.variants"),
		"Imports":             imports.String(),
		"PanelSchemas":        panelSchemas,
		"DataquerySchemas":    dataquerySchemas,
		"ShouldAddConverters": jenny.config.generateConverters,
	})
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

func (jenny Registry) unknownDataquery() (string, error) {
	return jenny.tmpl.Render("runtime/unknown_dataquery.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
	})
}

func (jenny Registry) unknownDataquerySerializer() (string, error) {
	return jenny.tmpl.Render("runtime/unknown_dataquery_serializer.tmpl", map[string]any{
		"Package": jenny.formatPackage("cog.variants"),
	})
}
