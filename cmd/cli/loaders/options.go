package loaders

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
)

type LoaderRef string

const (
	CUE               LoaderRef = "cue"
	KindsysCore       LoaderRef = "kindsys-core"
	KindsysComposable LoaderRef = "kindsys-composable"
	KindsysCustom     LoaderRef = "kindsys-custom"
	JSONSchema        LoaderRef = "jsonschema"
	OpenAPI           LoaderRef = "openapi"
	KindRegistry      LoaderRef = "kind-registry"
)

func loadersMap() map[LoaderRef]Loader {
	return map[LoaderRef]Loader{
		CUE:               cueLoader,
		KindsysCore:       kindsysCoreLoader,
		KindsysComposable: kindsysComposableLoader,
		KindsysCustom:     kindsysCustomLoader,
		JSONSchema:        jsonschemaLoader,
		OpenAPI:           openapiLoader,
		KindRegistry:      kindRegistryLoader,
	}
}

type Loader func(opts Options) ([]*ast.Schema, error)

type Options struct {
	CueEntrypoints               []string
	KindsysCoreEntrypoints       []string
	KindsysComposableEntrypoints []string
	KindsysCustomEntrypoints     []string
	JSONSchemaEntrypoints        []string
	OpenAPIEntrypoints           []string
	KindRegistryPath             string

	// Cue-specific options
	CueImports []string

	// Kind registry-specific options
	KindRegistryVersion string
}

func (opts Options) cueIncludeImports() ([]simplecue.LibraryInclude, error) {
	if len(opts.CueImports) == 0 {
		return nil, nil
	}

	imports := make([]simplecue.LibraryInclude, len(opts.CueImports))
	for i, importDefinition := range opts.CueImports {
		parts := strings.Split(importDefinition, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("'%s' is not a valid import definition", importDefinition)
		}

		imports[i] = simplecue.LibraryInclude{
			FSPath:     parts[0],
			ImportPath: parts[1],
		}
	}

	return imports, nil
}

func ForSchemaType(schemaType LoaderRef) (Loader, error) {
	all := loadersMap()

	loader, ok := all[schemaType]
	if !ok {
		return nil, fmt.Errorf("no loader found for '%s'", schemaType)
	}

	return loader, nil
}

func LoadAll(opts Options) ([]*ast.Schema, error) {
	var allSchemas []*ast.Schema

	for loaderRef := range loadersMap() {
		loader, err := ForSchemaType(loaderRef)
		if err != nil {
			return nil, err
		}

		schemas, err := loader(opts)
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemas...)
	}

	return allSchemas, nil
}

func guessPackageFromFilename(filename string) string {
	pkg := filepath.Base(filepath.Dir(filename))
	if pkg != "." {
		return pkg
	}

	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}
