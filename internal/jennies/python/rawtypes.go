package python

import (
	"path/filepath"
	"strings"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
)

type RawTypes struct {
}

func (jenny RawTypes) JennyName() string {
	return "PythonRawTypes"
}

func (jenny RawTypes) Generate(context common.Context) (codejen.Files, error) {
	files := make(codejen.Files, 0, len(context.Schemas))

	for _, schema := range context.Schemas {
		output, err := jenny.generateSchema(context, schema)
		if err != nil {
			return nil, err
		}

		filename := filepath.Join("models", schema.Package+".py")

		files = append(files, *codejen.NewFile(filename, output, jenny))
	}

	return files, nil
}

func (jenny RawTypes) generateSchema(_ common.Context, _ *ast.Schema) ([]byte, error) {
	var buffer strings.Builder

	buffer.WriteString("# TODO")

	return []byte(buffer.String()), nil
}
