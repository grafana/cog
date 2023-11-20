package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func runPassOnSchema(t *testing.T, pass Pass, input *ast.Schema, expectedOutput *ast.Schema) {
	t.Helper()

	runPassOnSchemas(t, pass, ast.Schemas{input}, ast.Schemas{expectedOutput})
}

func runPassOnSchemas(t *testing.T, pass Pass, input ast.Schemas, expectedOutput ast.Schemas) {
	t.Helper()

	req := require.New(t)

	processedSchemas, err := pass.Process(input)
	req.NoError(err)
	req.Len(processedSchemas, len(input))
	for i := range input {
		req.Empty(cmp.Diff(expectedOutput[i], processedSchemas[i]))
	}
}
