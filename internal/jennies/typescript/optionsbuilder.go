package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
)

type OptionsBuilder struct {
}

func (jenny OptionsBuilder) JennyName() string {
	return "TypescriptOptionsBuilder"
}

func (jenny OptionsBuilder) Generate(_ context.Builders) (*codejen.File, error) {
	output := jenny.generateFile()

	return codejen.NewFile("options_builder_gen.ts", []byte(output), jenny), nil
}

func (jenny OptionsBuilder) generateFile() string {
	return `export interface CogOptionsBuilder<T> {
  build: () => T;
}
`
}
