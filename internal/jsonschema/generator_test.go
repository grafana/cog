package jsonschema

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/txtartest"
	"github.com/stretchr/testify/require"
)

func TestGenerateAST(t *testing.T) {
	test := txtartest.TxTarTest{
		Root: "../../testdata/jsonschema",
		Name: "jsonschema/GenerateAST",
	}

	test.Run(t, func(tc *txtartest.Test) {
		schemaAst, err := GenerateAST(getReader(tc), Config{Package: "grafanatest"})
		if err != nil {
			writeError(tc, err)
		} else {
			require.NotNil(t, schemaAst)
			writeIR(tc, schemaAst)
		}
	})
}

func getReader(tc *txtartest.Test) io.Reader {
	tc.Helper()

	for _, a := range tc.Archive.Files {
		if strings.HasSuffix(a.Name, ".json") {
			return bytes.NewReader(a.Data)
		}
	}

	tc.Error("Cannot find test files")

	return nil
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
