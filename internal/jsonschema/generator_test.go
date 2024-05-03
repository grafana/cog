package jsonschema

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
		req.NotNil(schemaAst)

		testutils.WriteIR(schemaAst, tc)
	})
}

func TestGenerateAST_parsesEnumValues(t *testing.T) {
	req := require.New(t)

	input := strings.NewReader(`{
  "$ref": "#/definitions/SortOrder",
  "definitions": {
    "SortOrder": {
      "type": "number",
      "enum": [
        1,
        2,
        3,
        4,
        5
      ]
    }
  },
  "$schema": "http://json-schema.org/draft-07/schema#"
}`)

	schemaAst, err := GenerateAST(input, Config{Package: "grafanatest"})
	req.NoError(err)
	req.NotNil(t, schemaAst)

	enumType := schemaAst.Objects.At(0).Type.Enum

	req.Equal(int64(1), enumType.Values[0].Value)
}

func getSchemaAsReader(tc *testutils.Test) io.Reader {
	tc.Helper()

	file, err := os.Open(filepath.Join(tc.RootDir, "schema.json"))
	if err != nil {
		tc.Fatalf("could not open schema: %s", err)
	}

	return file
}
