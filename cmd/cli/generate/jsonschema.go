package generate

import (
	"os"
	"path/filepath"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jsonschema"
)

func jsonschemaLoader(opts options) ([]*ast.File, error) {
	allSchemas := make([]*ast.File, 0, len(opts.entrypoints))
	for _, entrypoint := range opts.entrypoints {
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
