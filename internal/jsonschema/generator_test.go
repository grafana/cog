package jsonschema

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestGenerateAST(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../testdata/jsonschema",
		Name:         "GenerateAST",
	}

	test.Run(t, func(tc *testutils.Test) {
		req := require.New(tc)

		schemaAst, err := GenerateAST(getSchemaAsReader(tc), Config{Package: "grafanatest"})
		req.NoError(err)
		require.NotNil(t, schemaAst)

		writeIR(schemaAst, tc)
	})
}

func getSchemaAsReader(tc *testutils.Test) io.Reader {
	tc.Helper()

	file, err := os.Open(filepath.Join(tc.RootDir, "schema.json"))
	if err != nil {
		tc.Fatalf("could not open schema: %s", err)
	}

	return file
}

func writeIR(irFile *ast.Schema, tc *testutils.Test) {
	tc.Helper()

	marshaledIR, err := json.MarshalIndent(irFile, "", "  ")
	require.NoError(tc, err)

	tc.WriteFile(&codejen.File{
		RelativePath: "ir.json",
		Data:         marshaledIR,
	})
}
