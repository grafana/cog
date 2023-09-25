package loaders

import (
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/openapi"
)

func openapiLoader(opts Options) ([]*ast.File, error) {
	allSchemas := make([]*ast.File, 0, len(opts.Entrypoints))
	for _, entrypoint := range opts.Entrypoints {
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
