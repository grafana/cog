package python

import (
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/languages"
)

var _ codejen.OneToMany[languages.Context] = ModuleInit{}

type ModuleInit struct {
}

func (jenny ModuleInit) JennyName() string {
	return "PythonModuleInit"
}

func (jenny ModuleInit) Generate(context languages.Context) (codejen.Files, error) {
	module := func(name string) []byte {
		return []byte(fmt.Sprintf(`"""%s module"""
`, name))
	}

	boilerplate := []struct {
		name    string
		content []byte
	}{
		{"py.typed", nil},
		{"__init__.py", module("root")},
		{"builders/__init__.py", module("builders")},
		{"models/__init__.py", module("models")},
		{"cog/__init__.py", module("runtime")},
	}

	files := make(codejen.Files, 0, len(boilerplate)+len(context.Schemas))

	for _, file := range boilerplate {
		files = append(files, *codejen.NewFile(file.name, file.content, jenny))
	}

	return files, nil
}
