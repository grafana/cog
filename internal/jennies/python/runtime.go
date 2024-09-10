package python

import (
	"sort"
	gotemplate "text/template"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	tmpl *gotemplate.Template
}

func (jenny Runtime) JennyName() string {
	return "PythonRuntime"
}

func (jenny Runtime) Generate(context languages.Context) (codejen.Files, error) {
	builder, err := template.Render(jenny.tmpl, "runtime/builder.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	encoder, err := template.Render(jenny.tmpl, "runtime/encoder.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	models, err := template.Render(jenny.tmpl, "runtime/variant_models.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	runtime, err := template.Render(jenny.tmpl, "runtime/runtime.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	plugins, err := jenny.variantPlugins(context)
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("cog/builder.py", []byte(builder), jenny),
		*codejen.NewFile("cog/encoder.py", []byte(encoder), jenny),
		*codejen.NewFile("cog/variants.py", []byte(models), jenny),
		*codejen.NewFile("cog/runtime.py", []byte(runtime), jenny),
		*codejen.NewFile("cog/plugins.py", []byte(plugins), jenny),
	}, nil
}

func (jenny Runtime) variantPlugins(context languages.Context) (string, error) {
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

	rendered, err := template.Render(jenny.tmpl, "runtime/plugins.tmpl", map[string]any{
		"panel_schemas":     panelSchemas,
		"dataquery_schemas": dataquerySchemas,
		"imports":           imports,
	})
	if err != nil {
		return "", err
	}

	importStatements := imports.String()
	if importStatements != "" {
		importStatements += "\n"
	}

	return importStatements + rendered, nil
}
