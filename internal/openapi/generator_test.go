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
		schemaAst, err := GenerateAST(getFilePath(tc), Config{Package: "grafanatest"})
		if err != nil {
			writeError(tc, err)
		} else {
			require.NotNil(t, schemaAst)
			writeIR(tc, schemaAst)
		}
	})
}

func getFilePath(tc *txtartest.Test) string {
	tc.Helper()

	for _, f := range tc.Archive.Files {
		if strings.HasSuffix(f.Name, ".json") {
			file, _ := os.CreateTemp("../../testdata/openapi", "tmp.json")
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

func writeIR(tc *txtartest.Test, irFile *ast.Schema) {
	tc.Helper()

	marshaledIR, err := json.MarshalIndent(irFile, "", "  ")
	require.NoError(tc, err)

	_, err = tc.Write(marshaledIR)
	require.NoError(tc, err)
}

func writeError(tc *txtartest.Test, err error) {
	tc.Helper()

	_, err = tc.Write([]byte(err.Error()))
	require.NoError(tc, err)
}
