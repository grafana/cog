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
	models, err := renderTemplate("runtime/variant_models.tmpl", map[string]any{})
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("cog/variants.py", []byte(models), jenny),
	}, nil
}
