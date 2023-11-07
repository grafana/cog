package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
)

type Runtime struct {
}

func (jenny Runtime) JennyName() string {
	return "GoRuntime"
}

func (jenny Runtime) Generate(context context.Builders) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	runtime, err := jenny.Runtime()
	if err != nil {
		return nil, err
	}

	files = append(files, *codejen.NewFile("cog/runtime.go", []byte(runtime), jenny))

	return files, nil
}

func (jenny Runtime) Runtime() (string, error) {
	imports := newImportMap()
	imports.Add("cogvariants", "github.com/grafana/cog/generated/cog/variants")

	return renderTemplate("runtime.tmpl", map[string]any{
		"imports": imports.Format(),
	})
}
