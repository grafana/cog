package loaders

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
)

type KindRegistryInput struct {
	Path    string `yaml:"path"`
	Version string `yaml:"version"`
}

func (input *KindRegistryInput) InterpolateParameters(interpolator ParametersInterpolator) {
	input.Path = interpolator(input.Path)
	input.Version = interpolator(input.Version)
}

func (input *KindRegistryInput) LoadSchemas(config Config) (ast.Schemas, error) {
	var allSchemas ast.Schemas
	var cueImports []string
	var cueEntrypoints []string

	if input.Path == "" {
		return nil, nil
	}

	coreKindEntrypoints, err := locateEntrypoints(config, input, "core")
	if err != nil {
		return nil, fmt.Errorf("could not locate core kind entrypoints: %w", err)
	}

	composableKindEntrypoints, err := locateEntrypoints(config, input, "composable")
	if err != nil {
		return nil, fmt.Errorf("could not locate composable kind entrypoints: %w", err)
	}

	commonPkgPath := kindRegistryKindPath(config, input, "common")
	commonPkgExists, err := dirExists(commonPkgPath)
	if err != nil {
		return nil, fmt.Errorf("could not locate common package: %w", err)
	}
	if commonPkgExists {
		cueEntrypoints = append(cueEntrypoints, commonPkgPath)
		cueImports = append(cueImports, fmt.Sprintf("%s:%s", commonPkgPath, "github.com/grafana/grafana/packages/grafana-schema/src/common"))
	}

	kindLoader := func(loader func(config Config, input CueInput) (ast.Schemas, error), entrypoints []string) error {
		for _, entrypoint := range entrypoints {
			schemas, err := loader(config, CueInput{
				Entrypoint: entrypoint,
				CueImports: cueImports,
			})
			if err != nil {
				return err
			}
			allSchemas = append(allSchemas, schemas...)
		}

		return nil
	}

	// CUE entrypoints
	if err := kindLoader(cueLoader, cueEntrypoints); err != nil {
		return nil, err
	}

	// Core kinds
	if err := kindLoader(kindsysCoreLoader, coreKindEntrypoints); err != nil {
		return nil, err
	}

	// Composable kinds
	if err := kindLoader(kindsysComposableLoader, composableKindEntrypoints); err != nil {
		return nil, err
	}

	return allSchemas, nil
}

func kindRegistryRoot(config Config, input *KindRegistryInput) string {
	return filepath.Join(config.Path(input.Path), "grafana")
}

func kindRegistryKindPath(config Config, input *KindRegistryInput, kind string) string {
	return filepath.Join(kindRegistryRoot(config, input), input.Version, kind)
}

func locateEntrypoints(config Config, input *KindRegistryInput, kind string) ([]string, error) {
	directory := kindRegistryKindPath(config, input, kind)
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("could not open directory '%s': %w", directory, err)
	}

	results := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		results = append(results, filepath.Join(directory, file.Name()))
	}

	return results, nil
}

func dirExists(dir string) (bool, error) {
	stat, err := os.Stat(dir)
	//nolint:gocritic
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if !stat.IsDir() {
		return false, fmt.Errorf("'%s' is not a directory", dir)
	}

	return true, nil
}
