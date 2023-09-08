package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
)

type TypescriptOptionsBuilder struct {
}

func (jenny TypescriptOptionsBuilder) JennyName() string {
	return "TypescriptOptionsBuilder"
}

func (jenny TypescriptOptionsBuilder) Generate(files []*ast.File) (*codejen.File, error) {
	output := jenny.generateFile()

	return codejen.NewFile("options_builder_gen.ts", []byte(output), jenny), nil
}

func (jenny TypescriptOptionsBuilder) generateFile() string {
	return `export interface OptionsBuilder<T> {
  build: () => T;
}
`
}
