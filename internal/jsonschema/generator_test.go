package jsonschema

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestGenerateAST(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[string]{
		TestDataRoot: "../../testdata/jsonschema",
		Name:         "GenerateAST",
	}

	test.Run(t, func(tc *testutils.Test[string]) {
		req := require.New(tc)

		const schemaFile = "schema.json"
		schemaAst, err := GenerateAST(tc.OpenInput(schemaFile), Config{
			Package:    "grafanatest",
			SchemaPath: filepath.Join(tc.RootDir, schemaFile),
		})
		req.NoError(err)
		req.NotNil(schemaAst)

		tc.WriteJSON(testutils.GeneratorOutputFile, schemaAst)
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

func TestGenerateAST_parseDraft2020Array(t *testing.T) {
	req := require.New(t)

	input := strings.NewReader(`{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$ref": "#/definitions/EntryPoint",
  "definitions": {
    "EntryPoint": {
      "type": "array",
      "items": {"type": "string"}
    }
  }
}`)

	schemaAst, err := GenerateAST(input, Config{Package: "grafanatest"})
	req.NoError(err)
	req.NotNil(t, schemaAst)

	objectType := schemaAst.Objects.Get("EntryPoint").Type

	req.True(objectType.IsArray())
	req.True(objectType.AsArray().ValueType.IsScalar())
	req.Equal(ast.KindString, objectType.AsArray().ValueType.AsScalar().ScalarKind)
}
