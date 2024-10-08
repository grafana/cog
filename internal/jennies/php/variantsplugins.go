package php

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type VariantsPlugins struct {
	config Config
	tmpl   *template.Template
}

func (jenny VariantsPlugins) JennyName() string {
	return "PHPVariantsPlugins"
}

func (jenny VariantsPlugins) Generate(_ languages.Context) (codejen.Files, error) {
	dataqueryInterface, err := jenny.dataqueryVariant()
	if err != nil {
		return nil, err
	}
	dataqueryConfig, err := jenny.dataqueryConfig()
	if err != nil {
		return nil, err
	}

	panelcfgInterface, err := jenny.panelcfgInterface()
	if err != nil {
		return nil, err
	}
	panelcfgConfig, err := jenny.panelcfgConfig()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		dataqueryInterface,
		dataqueryConfig,
		panelcfgInterface,
		panelcfgConfig,
	}, nil
}

func (jenny VariantsPlugins) dataqueryVariant() (codejen.File, error) {
	rendered, err := jenny.tmpl.RenderAsBytes("runtime/dataquery_variant.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Cog/Dataquery.php", rendered, jenny), nil
}

func (jenny VariantsPlugins) dataqueryConfig() (codejen.File, error) {
	rendered, err := jenny.tmpl.RenderAsBytes("runtime/dataquery_config.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Cog/DataqueryConfig.php", rendered, jenny), nil
}

func (jenny VariantsPlugins) panelcfgInterface() (codejen.File, error) {
	rendered, err := jenny.tmpl.RenderAsBytes("runtime/panelcfg_variant.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Cog/Panelcfg.php", rendered, jenny), nil
}

func (jenny VariantsPlugins) panelcfgConfig() (codejen.File, error) {
	rendered, err := jenny.tmpl.RenderAsBytes("runtime/panelcfg_config.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Cog/PanelcfgConfig.php", rendered, jenny), nil
}
