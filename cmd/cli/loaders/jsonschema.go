package loaders

import (
	"os"
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jsonschema"
)

func jsonschemaLoader(opts Options) ([]*ast.File, error) {
	allSchemas := make([]*ast.File, 0, len(opts.JSONSchemaEntrypoints))
	for _, entrypoint := range opts.JSONSchemaEntrypoints {
		pkg := filepath.Base(filepath.Dir(entrypoint))

		reader, err := os.Open(entrypoint)
		if err != nil {
			return nil, err
		}

		schemaAst, err := jsonschema.GenerateAST(reader, jsonschema.Config{
			Package: pkg, // TODO: extract from input schema/folder?
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}
