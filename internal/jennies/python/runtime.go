package python

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
)

type Runtime struct {
}

func (jenny Runtime) JennyName() string {
	return "PythonRuntime"
}

func (jenny Runtime) Generate(_ common.Context) (codejen.Files, error) {
	builder, err := renderTemplate("runtime/builder.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	encoder, err := renderTemplate("runtime/encoder.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	models, err := renderTemplate("runtime/variant_models.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("cog/builder.py", []byte(builder), jenny),
		*codejen.NewFile("cog/encoder.py", []byte(encoder), jenny),
		*codejen.NewFile("cog/variants.py", []byte(models), jenny),
	}, nil
}
