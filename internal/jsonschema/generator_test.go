package jsonschema

import (
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
		req := require.New(tc)

		schemaAst, err := GenerateAST(getSchemaAsReader(tc), Config{Package: "grafanatest"})
		req.NoError(err)
		require.NotNil(t, schemaAst)

		writeIR(schemaAst, tc)
	})
}

func getSchemaAsReader(tc *txtartest.Test) io.Reader {
	tc.Helper()

	for _, f := range tc.Archive.Files {
		if strings.HasSuffix(f.Name, ".json") {
			return strings.NewReader(string(f.Data))
		}
	}

	tc.Fatal("could not load types IR: file '*.json' not found in test archive")

	return nil
}

func writeIR(irFile *ast.Schema, tc *txtartest.Test) {
	tc.Helper()

	marshaledIR, err := json.MarshalIndent(irFile, "", "  ")
	require.NoError(tc, err)

	_, err = tc.Write(marshaledIR)
	require.NoError(tc, err)
}
