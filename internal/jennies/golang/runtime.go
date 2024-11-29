package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	Tmpl   *template.Template
	Config Config
}

func (jenny Runtime) JennyName() string {
	return "GoRuntime"
}

func (jenny Runtime) Generate(_ languages.Context) (codejen.Files, error) {
	runtime, err := jenny.runtime()
	if err != nil {
		return nil, err
	}

	tools, err := jenny.tools()
	if err != nil {
		return nil, err
	}

	files := []codejen.File{
		*codejen.NewFile("cog/runtime.go", runtime, jenny),
		*codejen.NewFile("cog/errors.go", tools, jenny),
	}

	if jenny.Config.generateBuilders {
		files = append(files,
			*codejen.NewFile("cog/builder.go", jenny.builderInterface(), jenny),
		)
	}

	return files, nil
}

func (jenny Runtime) builderInterface() []byte {
	return []byte(`package cog

type Builder[ResourceT any] interface {
  Build() (ResourceT, error)
}

`)
}

func (jenny Runtime) runtime() ([]byte, error) {
	imports := NewImportMap(jenny.Config.PackageRoot)
	imports.Add("", jenny.Config.importPath("cog/variants"))
	imports.Add("json", "encoding/json")
	imports.Add("fmt", "fmt")
	imports.Add("reflect", "reflect")
	imports.Add("strings", "strings")

	return jenny.Tmpl.RenderAsBytes("runtime/runtime.tmpl", map[string]any{
		"imports": imports,
	})
}

func (jenny Runtime) tools() ([]byte, error) {
	return jenny.Tmpl.RenderAsBytes("runtime/tools.tmpl", map[string]any{})
}
