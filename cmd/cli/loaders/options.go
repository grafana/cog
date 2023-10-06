package loaders

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type LoaderRef string

const (
	CUE               LoaderRef = "cue"
	KindsysCore       LoaderRef = "kindsys-core"
	KindsysComposable LoaderRef = "kindsys-composable"
	KindsysCustom     LoaderRef = "kindsys-custom"
	JSONSchema        LoaderRef = "jsonschema"
	OpenAPI           LoaderRef = "openapi"
)

func loadersMap() map[LoaderRef]Loader {
	return map[LoaderRef]Loader{
		CUE:               cueLoader,
		KindsysCore:       kindsysCoreLoader,
		KindsysComposable: kindsysCompopsableLoader,
		KindsysCustom:     kindsysCustomLoader,
		JSONSchema:        jsonschemaLoader,
		OpenAPI:           openapiLoader,
	}
}

type cueIncludeImport struct {
	fsPath     string // path of the library on the filesystem
	importPath string // path used in CUE files to import that library
}

type Loader func(opts Options) ([]*ast.Schema, error)

type Options struct {
	CueEntrypoints               []string
	KindsysCoreEntrypoints       []string
	KindsysComposableEntrypoints []string
	KindsysCustomEntrypoints     []string
	JSONSchemaEntrypoints        []string
	OpenAPIEntrypoints           []string

	// Cue-specific options
	CueImports []string
}

func (opts Options) cueIncludeImports() ([]cueIncludeImport, error) {
	if len(opts.CueImports) == 0 {
		return nil, nil
	}

	imports := make([]cueIncludeImport, len(opts.CueImports))
	for i, importDefinition := range opts.CueImports {
		parts := strings.Split(importDefinition, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("'%s' is not a valid import definition", importDefinition)
		}

		imports[i].fsPath = parts[0]
		imports[i].importPath = parts[1]
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
