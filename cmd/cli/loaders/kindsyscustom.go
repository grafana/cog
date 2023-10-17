package loaders

import (
	"fmt"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
	"github.com/grafana/kindsys"
	"github.com/grafana/thema"
)

func kindsysCustomLoader(opts Options) ([]*ast.Schema, error) {
	themaRuntime := thema.NewRuntime(cuecontext.New())

	libraries, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	allSchemas := make([]*ast.Schema, 0, len(opts.KindsysCustomEntrypoints))
	for _, entrypoint := range opts.KindsysCustomEntrypoints {
		pkg := filepath.Base(entrypoint)

		overlayFS, err := buildKindsysEntrypointFS(opts, entrypoint)
		if err != nil {
			return nil, err
		}

		cueInstance, err := kindsys.BuildInstance(themaRuntime.Context(), ".", pkg, overlayFS)
		if err != nil {
			return nil, fmt.Errorf("could not load kindsys instance: %w", err)
		}

		props, err := kindsys.ToKindProps[kindsys.CustomProperties](cueInstance)
		if err != nil {
			return nil, fmt.Errorf("could not convert cue value to kindsys props: %w", err)
		}

		kindDefinition := kindsys.Def[kindsys.CustomProperties]{
			V:          cueInstance,
			Properties: props,
		}

		boundKind, err := kindsys.BindCustom(themaRuntime, kindDefinition)
		if err != nil {
			return nil, fmt.Errorf("could not bind kind definition to kind: %w", err)
		}

		schemaAst, err := simplecue.GenerateAST(kindToLatestSchema(boundKind), simplecue.Config{
			Package: pkg, // TODO: extract from input schema/folder?
			SchemaMetadata: ast.SchemaMeta{
				Kind:       ast.SchemaKindCore, // TODO: is there any need for a "SchemaKindCustom"?
				Identifier: pkg,                // TODO: maybe even core kinds could have one explicitly set in their schema?
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
