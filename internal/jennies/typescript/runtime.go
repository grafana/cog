package typescript

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/template"
	"github.com/grafana/cog/internal/languages"
)

type Runtime struct {
	tmpl    *template.Template
	targets languages.Config
}

func (jenny Runtime) JennyName() string {
	return "TypescriptRuntime"
}

func (jenny Runtime) Generate(_ languages.Context) (codejen.Files, error) {
	runtime, err := jenny.generateRuntime()
	if err != nil {
		return nil, err
	}

	return codejen.Files{
		*codejen.NewFile("src/cog/variants_gen.ts", []byte(jenny.generateVariantsFile()), jenny),
		*codejen.NewFile("src/cog/builder_gen.ts", []byte(jenny.generateOptionsBuilderFile()), jenny),
		*codejen.NewFile("src/cog/runtime.ts", runtime, jenny),
		*codejen.NewFile("src/cog/index.ts", []byte(jenny.generateIndexFile()), jenny),
	}, nil
}

func (jenny Runtime) generateIndexFile() string {
	index := `export * from './variants_gen';
export * from './builder_gen';
`
	if jenny.targets.Converters {
		index += "export * from './runtime';\n"
	}

	return index
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

export function isBuilder<T>(input: Builder<T> | any): input is Builder<T> {
  if (input === null) {
    return false;
  }
  if (!input?.build) {
    return false;
  }

  return typeof input.build === "function";
}
`
}

func (jenny Runtime) generateRuntime() ([]byte, error) {
	return jenny.tmpl.RenderAsBytes("runtime/runtime.ts.tmpl", nil)
}
