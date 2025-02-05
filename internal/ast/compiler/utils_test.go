package compiler

import (
	"encoding/json"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

const testPkgName = "test"

func runPassOnObjects(t *testing.T, pass Pass, input []ast.Object, expectedOutput []ast.Object) {
	t.Helper()

	inputSchema := ast.NewSchema(testPkgName, ast.SchemaMeta{})
	inputSchema.AddObjects(input...)

	expectedOutputSchema := ast.NewSchema(testPkgName, ast.SchemaMeta{})
	expectedOutputSchema.AddObjects(expectedOutput...)

	runPassOnSchema(t, pass, inputSchema, expectedOutputSchema)
}

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
		expectedJSON, err := json.MarshalIndent(expectedOutput[i], "", "  ")
		req.NoError(err)
		gotJSON, err := json.MarshalIndent(processedSchemas[i], "", "  ")
		req.NoError(err)

		req.JSONEq(string(expectedJSON), string(gotJSON))
	}
}
