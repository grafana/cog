package terraform

import (
	"path/filepath"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
)

type RawTypes struct {
	config Config
}

func (jenny RawTypes) JennyName() string {
	return "TerraformRawTypes"
}

func (jenny RawTypes) Generate(context languages.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			formatPackageName(schema.Package),
			"types_gen.go",
		)

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}
	return files, nil
}

func (jenny RawTypes) generateSchema(_ languages.Context, _ *ast.Schema) ([]byte, error) {
	return []byte(""), nil
}
