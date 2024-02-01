package loaders

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/openapi"
)

func openapiLoader(opts Options) (ast.Schemas, error) {
	allSchemas := make([]*ast.Schema, 0, len(opts.OpenAPIEntrypoints))
	for _, entrypoint := range opts.OpenAPIEntrypoints {
		schemaAst, err := openapi.GenerateAST(entrypoint, openapi.Config{
			Package:        guessPackageFromFilename(entrypoint),
			SchemaMetadata: ast.SchemaMeta{
				// TODO: extract these from somewhere
			},
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}
