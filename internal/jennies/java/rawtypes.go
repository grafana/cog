package java

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"path/filepath"
	"strings"
)

type RawTypes struct {
	config Config
}

func (jenny RawTypes) JennyName() string {
	return "JavaRawTypes"
}

func (jenny RawTypes) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, len(context.Schemas))

	for i, schema := range context.Schemas {
		output, err := jenny.generateSchema(schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join(
			strings.ToLower(schema.Package),
			"types_gen.java",
		)

		files[i] = *codejen.NewFile(filename, output, jenny)
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(schema *ast.Schema) ([]byte, error) {
	return nil, nil
}
