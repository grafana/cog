package yaml

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/grafana/cog/internal/ast/compiler"
)

type Compiler struct {
	Passes []CompilerPass `yaml:"passes"`
}

type CompilerLoader struct {
}

func NewCompilerLoader() *CompilerLoader {
	return &CompilerLoader{}
}

func (loader *CompilerLoader) PassesFrom(filenames []string) (compiler.Passes, error) {
	allPasses := make(compiler.Passes, 0, len(filenames))

	for _, filename := range filenames {
		passes, err := loader.load(filename)
		if err != nil {
			return nil, err
		}

		allPasses = append(allPasses, passes...)
	}

	return allPasses, nil
}

func (loader *CompilerLoader) load(file string) (compiler.Passes, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	compilerConfig := &Compiler{}
	if err := yaml.UnmarshalWithOptions(contents, compilerConfig, yaml.DisallowUnknownField()); err != nil {
		return nil, fmt.Errorf("can not load compiler passes: %s\n%s", file, yaml.FormatError(err, true, true))
	}

	passes := make(compiler.Passes, 0, len(compilerConfig.Passes))

	// convert compiler passes
	for i, passConfig := range compilerConfig.Passes {
		pass, err := passConfig.AsCompilerPass()
		if err != nil {
			path, innerErr := yaml.PathString(fmt.Sprintf("$.passes[%d]", i))
			if innerErr != nil {
				return nil, err
			}
			source, innerErr := path.AnnotateSource(contents, true)
			if innerErr != nil {
				return nil, err
			}

			return nil, fmt.Errorf("%w in %s\n%s", err, file, string(source))
		}

		passes = append(passes, pass)
	}

	return passes, nil
}
