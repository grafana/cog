package loaders

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type LoaderRef string

const (
	CUE           LoaderRef = "cue"
	KindsysCore   LoaderRef = "kindsys-core"
	KindsysCustom LoaderRef = "kindsys-custom"
	JSONSchema    LoaderRef = "jsonschema"
)

type cueIncludeImport struct {
	fsPath     string // path of the library on the filesystem
	importPath string // path used in CUE files to import that library
}

type Loader func(opts Options) ([]*ast.File, error)

type Options struct {
	CueEntrypoints           []string
	KindsysCoreEntrypoints   []string
	KindsysCustomEntrypoints []string
	JSONSchemaEntrypoints    []string

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
	all := map[LoaderRef]Loader{
		CUE:           cueLoader,
		KindsysCore:   kindsysCoreLoader,
		KindsysCustom: kindsysCustomLoader,
		JSONSchema:    jsonschemaLoader,
	}

	loader, ok := all[schemaType]
	if !ok {
		return nil, fmt.Errorf("no loader found for '%s'", schemaType)
	}

	return loader, nil
}

func LoadAll(opts Options) ([]*ast.File, error) {
	var files []*ast.File

	loaders := []LoaderRef{
		CUE,
		KindsysCore,
		KindsysCustom,
		JSONSchema,
	}

	for _, loaderRef := range loaders {
		loader, err := ForSchemaType(loaderRef)
		if err != nil {
			return nil, err
		}

		schemas, err := loader(opts)
		if err != nil {
			return nil, err
		}

		files = append(files, schemas...)
	}

	return files, nil
}
