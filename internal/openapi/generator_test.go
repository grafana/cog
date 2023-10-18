package openapi

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/txtartest"
	"github.com/stretchr/testify/require"
)

func TestGenerateAST(t *testing.T) {
	test := txtartest.TxTarTest{
		Root: "../../testdata/openapi",
		Name: "openapi/GenerateAST",
	}

	test.Run(t, func(tc *txtartest.Test) {
		req := require.New(tc)

		schemaAst, err := GenerateAST(getFilePath(tc), Config{Package: "grafanatest"})
		req.NoError(err)
		require.NotNil(t, schemaAst)

		writeIR(schemaAst, tc)
	})
}

func getFilePath(tc *txtartest.Test) string {
	tc.Helper()

	for _, f := range tc.Archive.Files {
		if strings.HasSuffix(f.Name, ".json") {
			file, _ := os.CreateTemp(".", "tmp.json")
			_, _ = file.Write(f.Data)
			tc.Cleanup(func() {
				_ = os.Remove(file.Name())
			})
			return file.Name()
		}
	}

	tc.Fatal("could not load types IR: file '*.json' not found in test archive")

	return ""
}

func writeIR(irFile *ast.Schema, tc *txtartest.Test) {
	tc.Helper()

	marshaledIR, err := json.MarshalIndent(irFile, "", "  ")
	require.NoError(tc, err)

	_, err = tc.Write(marshaledIR)
	require.NoError(tc, err)
}
