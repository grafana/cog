package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
)

type Runtime struct {
}

func (jenny Runtime) JennyName() string {
	return "TypescriptRuntime"
}

func (jenny Runtime) Generate(_ context.Builders) (codejen.Files, error) {
	return codejen.Files{
		*codejen.NewFile("cog/variants_gen.ts", []byte(jenny.generateVariantsFile()), jenny),
		*codejen.NewFile("cog/options_builder_gen.ts", []byte(jenny.generateOptionsBuilderFile()), jenny),
		*codejen.NewFile("cog/index.ts", []byte(jenny.generateIndexFile()), jenny),
	}, nil
}

func (jenny Runtime) generateIndexFile() string {
	return `export * from './variants_gen';
export * from './options_builder_gen';
`
}

func (jenny Runtime) generateVariantsFile() string {
	return `export interface Dataquery {
	implementsDataqueryVariant(): void;
}

export interface Panelcfg {
	implementsPanelcfgVariant(): void;
}
`
}

func (jenny Runtime) generateOptionsBuilderFile() string {
	return `export interface OptionsBuilder<T> {
  build: () => T;
}
`
}
