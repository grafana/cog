package typescript

import (
	"fmt"
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/common"
	"strings"
)

type Runtime struct {
	RuntimePath *string
}

func (jenny Runtime) JennyName() string {
	return "TypescriptRuntime"
}

func (jenny Runtime) Generate(_ common.Context) (codejen.Files, error) {
	path := getPath(jenny.RuntimePath)
	return codejen.Files{
		*codejen.NewFile(fmt.Sprintf("%scog/variants_gen.ts", path), []byte(jenny.generateVariantsFile()), jenny),
		*codejen.NewFile(fmt.Sprintf("%scog/builder_gen.ts", path), []byte(jenny.generateOptionsBuilderFile()), jenny),
		*codejen.NewFile(fmt.Sprintf("%scog/index.ts", path), []byte(jenny.generateIndexFile()), jenny),
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

func getPath(runtimePath *string) string {
	if runtimePath == nil {
		return "src"
	}

	path := *runtimePath
	if path != "" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path
}
