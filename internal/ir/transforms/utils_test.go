package transforms

import (
	"encoding/json"
	"testing"

	"github.com/grafana/cog/internal/ir"
	"github.com/stretchr/testify/require"
)

const testPkgName = "test"

func runPassOnObjects(t *testing.T, pass Transform, input []ir.Object, expectedOutput []ir.Object) {
	t.Helper()

	inputSchema := ir.NewSchema(testPkgName, ir.SchemaMeta{})
	inputSchema.AddObjects(input...)

	expectedOutputSchema := ir.NewSchema(testPkgName, ir.SchemaMeta{})
	expectedOutputSchema.AddObjects(expectedOutput...)

	runPassOnSchema(t, pass, inputSchema, expectedOutputSchema)
}

func runPassOnSchema(t *testing.T, pass Transform, input *ir.Schema, expectedOutput *ir.Schema) {
	t.Helper()

	runPassOnSchemas(t, pass, ir.Schemas{input}, ir.Schemas{expectedOutput})
}

func runPassOnSchemas(t *testing.T, pass Transform, input ir.Schemas, expectedOutput ir.Schemas) {
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
