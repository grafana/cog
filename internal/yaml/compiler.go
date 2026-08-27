package yaml

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/grafana/cog/pkg/ir/transforms"
)

type Transforms struct {
	Passes []Transform `yaml:"passes"`
}

type TransformsLoader struct {
}

func NewTransformsLoader() *TransformsLoader {
	return &TransformsLoader{}
}

func (loader *TransformsLoader) LoadFiles(filenames []string) (transforms.Transforms, error) {
	allPasses := make(transforms.Transforms, 0, len(filenames))

	for _, filename := range filenames {
		passes, err := loader.load(filename)
		if err != nil {
			return nil, err
		}

		allPasses = append(allPasses, passes...)
	}

	return allPasses, nil
}

func (loader *TransformsLoader) load(file string) (transforms.Transforms, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	compilerConfig := &Transforms{}
	if err := yaml.UnmarshalWithOptions(contents, compilerConfig, yaml.DisallowUnknownField()); err != nil {
		return nil, fmt.Errorf("can not load compiler passes: %s\n%s", file, yaml.FormatError(err, true, true))
	}

	passes := make(transforms.Transforms, 0, len(compilerConfig.Passes))

	// convert compiler passes
	for i, passConfig := range compilerConfig.Passes {
		pass, err := passConfig.AsTransform()
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
