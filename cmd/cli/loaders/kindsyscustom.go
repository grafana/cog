package loaders

import (
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
)

func kindsysCustomLoader(opts Options) ([]*ast.Schema, error) {
	cueFsOverlay, err := buildCueOverlay(opts)
	if err != nil {
		return nil, err
	}

	libraries, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	allSchemas := make([]*ast.Schema, 0, len(opts.KindsysCustomEntrypoints))
	for _, entrypoint := range opts.KindsysCustomEntrypoints {
		pkg := filepath.Base(entrypoint)

		// Load Cue files into Cue build.Instances slice
		// the second arg is a configuration object, we'll see this later
		bis := load.Instances([]string{entrypoint}, &load.Config{
			Overlay:    cueFsOverlay,
			ModuleRoot: "/",
		})

		values, err := cuecontext.New().BuildInstances(bis)
		if err != nil {
			return nil, err
		}

		schemaRoot := values[0]
		schemaAsCueValue := schemaRoot.LookupPath(cue.ParsePath("lineage.schemas[0].schema"))

		kindIdentifier, err := inferCoreKindIdentifier(schemaRoot) // same strategy than with core kinds
		if err != nil {
			return nil, err
		}

		schemaAst, err := simplecue.GenerateAST(schemaAsCueValue, simplecue.Config{
			Package: pkg, // TODO: extract from somewhere else?
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindCore,
				Identifier: kindIdentifier,
			},
			Libraries: libraries,
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}
