package python

import (
	"sort"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	tmpl *template.Template
}

func (jenny Runtime) JennyName() string {
	return "PythonRuntime"
}

func (jenny Runtime) Generate(context languages.Context) (codejen.Files, error) {
	builder, err := jenny.tmpl.RenderAsBytes("runtime/builder.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	encoder, err := jenny.tmpl.RenderAsBytes("runtime/encoder.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	models, err := jenny.tmpl.RenderAsBytes("runtime/variant_models.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	runtime, err := jenny.tmpl.RenderAsBytes("runtime/runtime.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	plugins, err := jenny.variantPlugins(context)
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("cog/builder.py", builder, jenny),
		*codejen.NewFile("cog/encoder.py", encoder, jenny),
		*codejen.NewFile("cog/variants.py", models, jenny),
		*codejen.NewFile("cog/runtime.py", runtime, jenny),
		*codejen.NewFile("cog/plugins.py", plugins, jenny),
	}, nil
}

func (jenny Runtime) variantPlugins(context languages.Context) ([]byte, error) {
	imports := NewImportMap()
	var panelSchemas []string
	var dataquerySchemas []string

	for _, schema := range context.Schemas {
		if schema.Metadata.Kind != ast.SchemaKindComposable || schema.Metadata.Identifier == "" {
			continue
		}

		importAlias := imports.AddModule(schema.Package, "..models", schema.Package)

		if schema.Metadata.Variant == ast.SchemaVariantPanel {
			panelSchemas = append(panelSchemas, importAlias)
		} else if schema.Metadata.Variant == ast.SchemaVariantDataQuery {
			dataquerySchemas = append(dataquerySchemas, importAlias)
		}
	}

	// to guarantee a consistent output for this jenny
	sort.Strings(panelSchemas)
	sort.Strings(dataquerySchemas)

	rendered, err := jenny.tmpl.Render("runtime/plugins.tmpl", map[string]any{
		"panel_schemas":     panelSchemas,
		"dataquery_schemas": dataquerySchemas,
		"imports":           imports,
	})
	if err != nil {
		return nil, err
	}

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n"
	}

	return []byte(importStatements + rendered), nil
}
