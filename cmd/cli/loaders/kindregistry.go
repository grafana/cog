package loaders

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
)

func kindRegistryLoader(opts Options) ([]*ast.Schema, error) {
	var allSchemas []*ast.Schema

	if opts.KindRegistryPath == "" {
		return nil, nil
	}

	coreKindEntrypoints, err := locateEntrypoints(opts, "core")
	if err != nil {
		return nil, fmt.Errorf("could not locate core kind entrypoints: %w", err)
	}

	composableKindEntrypoints, err := locateEntrypoints(opts, "composable")
	if err != nil {
		return nil, fmt.Errorf("could not locate composable kind entrypoints: %w", err)
	}

	newOpts := opts
	newOpts.KindsysCoreEntrypoints = coreKindEntrypoints
	newOpts.KindsysComposableEntrypoints = composableKindEntrypoints

	commonPkgPath := kindRegistryKindPath(opts, "common")
	commonPkgExists, err := dirExists(commonPkgPath)
	if err != nil {
		return nil, fmt.Errorf("could not locate common package: %w", err)
	}
	if commonPkgExists {
		newOpts.CueEntrypoints = append(newOpts.CueEntrypoints, commonPkgPath)
		newOpts.CueImports = append(newOpts.CueImports, fmt.Sprintf("%s:%s", commonPkgPath, "github.com/grafana/grafana/packages/grafana-schema/src/common"))
	}

	cueSchemas, err := cueLoader(newOpts)
	if err != nil {
		return nil, err
	}

	coreSchemas, err := kindsysCoreLoader(newOpts)
	if err != nil {
		return nil, err
	}

	composableSchemas, err := kindsysComposableLoader(newOpts)
	if err != nil {
		return nil, err
	}

	allSchemas = append(allSchemas, cueSchemas...)
	allSchemas = append(allSchemas, coreSchemas...)
	allSchemas = append(allSchemas, composableSchemas...)

	return allSchemas, nil
}

func kindRegistryRoot(opts Options) string {
	return filepath.Join(opts.KindRegistryPath, "grafana")
}

func kindRegistryKindPath(opts Options, kind string) string {
	return filepath.Join(kindRegistryRoot(opts), opts.KindRegistryVersion, kind)
}

func locateEntrypoints(opts Options, kind string) ([]string, error) {
	directory := kindRegistryKindPath(opts, kind)
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
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if !stat.IsDir() {
		return false, fmt.Errorf("'%s' is not a directory", dir)
	}

	return true, nil
}
