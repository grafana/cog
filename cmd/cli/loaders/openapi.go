package loaders

import (
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/openapi"
)

func openapiLoader(opts Options) ([]*ast.Schema, error) {
	allSchemas := make([]*ast.Schema, 0, len(opts.OpenAPIEntrypoints))
	for _, entrypoint := range opts.OpenAPIEntrypoints {
		pkg := filepath.Base(filepath.Dir(entrypoint))
		schemaAst, err := openapi.GenerateAST(entrypoint, openapi.Config{
			Package: pkg,
		})
		if err != nil {
			return nil, err
		}

		allSchemas = append(allSchemas, schemaAst)
	}

	return allSchemas, nil
}
