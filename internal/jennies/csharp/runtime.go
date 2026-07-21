package csharp

import (
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

// Runtime emits the small set of runtime helpers that generated
// builders depend on (currently just `Cog.IBuilder<T>`). The runtime
// lives under `<ProjectPath>/Cog/` in the `<NamespaceRoot>.Cog`
// namespace.
type Runtime struct {
	config Config
	tmpl   *template.Template
}

func (jenny Runtime) JennyName() string {
	return "CSharpRuntime"
}

func (jenny Runtime) Generate(_ languages.Context) (codejen.Files, error) {
	builder, err := jenny.tmpl.RenderAsBytes("runtime/builder.tmpl", map[string]any{
		"Namespace": jenny.config.formatNamespace("cog"),
	})
	if err != nil {
		return nil, err
	}
	return codejen.Files{
		*codejen.NewFile(filepath.Join(jenny.config.ProjectPath, "Cog", "IBuilder.cs"), builder, jenny),
	}, nil
}
