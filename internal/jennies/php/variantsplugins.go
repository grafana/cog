package php

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/languages"
)

type VariantsPlugins struct {
	config Config
}

func (jenny VariantsPlugins) JennyName() string {
	return "PHPVariantsPlugins"
}

func (jenny VariantsPlugins) Generate(_ languages.Context) (codejen.Files, error) {
	dataqueryInterface, err := jenny.dataqueryVariant()
	if err != nil {
		return nil, err
	}

	panelcfgInterface, err := jenny.panelcfgInterface()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		dataqueryInterface,
		panelcfgInterface,
	}, nil
}

func (jenny VariantsPlugins) dataqueryVariant() (codejen.File, error) {
	rendered, err := renderTemplate("runtime/dataquery_variant.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Runtime/Dataquery.php", []byte(rendered), jenny), nil
}

func (jenny VariantsPlugins) panelcfgInterface() (codejen.File, error) {
	rendered, err := renderTemplate("runtime/panelcfg_variant.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Runtime/Panelcfg.php", []byte(rendered), jenny), nil
}
