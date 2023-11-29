package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
)

type Runtime struct {
}

func (jenny Runtime) JennyName() string {
	return "TypescriptRuntime"
}

func (jenny Runtime) Generate(_ common.Context) (codejen.Files, error) {
	return codejen.Files{
		*codejen.NewFile("cog/variants_gen.ts", []byte(jenny.generateVariantsFile()), jenny),
		*codejen.NewFile("cog/builder_gen.ts", []byte(jenny.generateOptionsBuilderFile()), jenny),
		*codejen.NewFile("cog/index.ts", []byte(jenny.generateIndexFile()), jenny),
	}, nil
}

func (jenny Runtime) generateIndexFile() string {
	return `export * from './variants_gen';
export * from './builder_gen';
`
}

func (jenny Runtime) generateVariantsFile() string {
	return `export interface Dataquery {
	_implementsDataqueryVariant(): void;
}

`
}

func (jenny Runtime) generateOptionsBuilderFile() string {
	return `export interface Builder<T> {
  build: () => T;
}
`
}
