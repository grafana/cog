package loaders

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
)

type cueIncludeImport struct {
	fsPath     string // path of the library on the filesystem
	importPath string // path used in CUE files to import that library
}

type Loader func(opts Options) ([]*ast.File, error)

type Options struct {
	Entrypoints []string
	SchemasType string

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

func ForSchemaType(schemaType string) (Loader, error) {
	all := map[string]Loader{
		"cue":            cueLoader,
		"kindsys-core":   kindsysCoreLoader,
		"kindsys-custom": kindsysCustomLoader,
		"jsonschema":     jsonschemaLoader,
	}

	loader, ok := all[schemaType]
	if !ok {
		return nil, fmt.Errorf("no loader found for '%s'", schemaType)
	}

	return loader, nil
}
