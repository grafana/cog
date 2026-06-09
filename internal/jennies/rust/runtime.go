package rust

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

// Runtime emits the hand-written Rust runtime crate modules that every
// generated SDK depends on. The modules live under src/cog/ to match a normal
// Rust crate layout.
type Runtime struct {
	tmpl *template.Template
}

func (jenny Runtime) JennyName() string {
	return "RustRuntime"
}

func (jenny Runtime) Generate(_ languages.Context) (codejen.Files, error) {
	modules := map[string]string{
		"src/cog/mod.rs":      "runtime/mod.tmpl",
		"src/cog/builder.rs":  "runtime/builder.tmpl",
		"src/cog/error.rs":    "runtime/error.tmpl",
		"src/cog/variants.rs": "runtime/variants.tmpl",
	}

	files := make(codejen.Files, 0, len(modules))
	for output, tmplFile := range modules {
		rendered, err := jenny.tmpl.RenderAsBytes(tmplFile, map[string]any{})
		if err != nil {
			return nil, err
		}
		files = append(files, *codejen.NewFile(output, rendered, jenny))
	}

	return files, nil
}
