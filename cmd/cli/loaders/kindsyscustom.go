package loaders

import (
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/simplecue"
)

func kindsysCustomLoader(opts Options) (ast.Schemas, error) {
	libraries, err := opts.cueIncludeImports()
	if err != nil {
		return nil, err
	}

	allSchemas := make([]*ast.Schema, 0, len(opts.KindsysCustomEntrypoints))
	for _, entrypoint := range opts.KindsysCustomEntrypoints {
		pkg := filepath.Base(entrypoint)

		schemaRootValue, err := parseCueEntrypoint(opts, entrypoint)
		if err != nil {
			return nil, err
		}

		kindIdentifier, err := inferCoreKindIdentifier(schemaRootValue) // same strategy than with core kinds
		if err != nil {
			return nil, err
		}

		schemaAst, err := simplecue.GenerateAST(schemaFromThemaLineage(schemaRootValue), simplecue.Config{
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
