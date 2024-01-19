package openapi

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestGenerateAST(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../testdata/openapi",
		Name:         "GenerateAST",
	}

	test.Run(t, func(tc *testutils.Test) {
		req := require.New(tc)

		schemaAst, err := GenerateAST(getFilePath(tc), Config{Package: "grafanatest"})
		req.NoError(err)
		require.NotNil(t, schemaAst)

		writeIR(schemaAst, tc)
	})
}

func getFilePath(tc *testutils.Test) string {
	tc.Helper()

	return filepath.Join(tc.RootDir, "schema.json")
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
