package php

import (
	"sort"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	config Config
}

func (jenny Runtime) JennyName() string {
	return "PHPRuntime"
}

func (jenny Runtime) Generate(context languages.Context) (codejen.Files, error) {
	runtime, err := jenny.runtime(context)
	if err != nil {
		return nil, err
	}

	unknownDataquery, err := jenny.unknownDataquery()
	if err != nil {
		return nil, err
	}

	builderInterface, err := jenny.builderInterface()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		runtime,
		builderInterface,
		unknownDataquery,
	}, nil
}

func (jenny Runtime) builderInterface() (codejen.File, error) {
	rendered, err := renderTemplate("runtime/builder.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Cog/Builder.php", []byte(rendered), jenny), nil
}

func (jenny Runtime) runtime(context languages.Context) (codejen.File, error) {
	var panelSchemas []*ast.Schema
	var dataquerySchemas []*ast.Schema

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Identifier == "" {
			continue
		}

		if schema.Metadata.Variant == ast.SchemaVariantPanel {
			panelSchemas = append(panelSchemas, schema)
		} else if schema.Metadata.Variant == ast.SchemaVariantDataQuery {
			dataquerySchemas = append(dataquerySchemas, schema)
		}
	}

	// to guarantee a consistent output for this jenny
	sort.SliceStable(panelSchemas, func(i, j int) bool {
		return panelSchemas[i].Package < panelSchemas[j].Package
	})
	sort.SliceStable(dataquerySchemas, func(i, j int) bool {
		return dataquerySchemas[i].Package < dataquerySchemas[j].Package
	})

	rendered, err := renderTemplate("runtime/runtime.tmpl", map[string]any{
		"PanelSchemas":     panelSchemas,
		"DataquerySchemas": dataquerySchemas,
		"NamespaceRoot":    jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Cog/Runtime.php", []byte(rendered), jenny), nil
}

func (jenny Runtime) unknownDataquery() (codejen.File, error) {
	rendered, err := renderTemplate("runtime/unknown_dataquery.tmpl", map[string]any{
		"NamespaceRoot": jenny.config.NamespaceRoot,
	})
	if err != nil {
		return codejen.File{}, err
	}

	return *codejen.NewFile("src/Cog/UnknownDataquery.php", []byte(rendered), jenny), nil
}
