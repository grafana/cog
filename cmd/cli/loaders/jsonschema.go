package loaders

import (
	"os"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jsonschema"
)

func jsonschemaLoader(opts Options) (ast.Schemas, error) {
	allSchemas := make([]*ast.Schema, 0, len(opts.JSONSchemaEntrypoints))
	for _, entrypoint := range opts.JSONSchemaEntrypoints {
		reader, err := os.Open(entrypoint)
		if err != nil {
			return nil, err
		}

		schemaAst, err := jsonschema.GenerateAST(reader, jsonschema.Config{
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
