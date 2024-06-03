package openapi

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestGenerateAST(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[string]{
		TestDataRoot: "../../testdata/openapi",
		Name:         "GenerateAST",
	}

	test.Run(t, func(tc *testutils.Test[string]) {
		req := require.New(tc)
		ctx := context.TODO()

		schemaAst, err := GenerateAST(ctx, getSchemaAsReader(tc), Config{Package: "grafanatest"})
		req.NoError(err)
		require.NotNil(t, schemaAst)

		tc.WriteJSON(testutils.GeneratorOutputFile, schemaAst)
	})
}

func getSchemaAsReader(tc *testutils.Test[string]) *openapi3.T {
	tc.Helper()

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	oapi, err := loader.LoadFromFile(filepath.Join(tc.RootDir, "schema.json"))
	if err != nil {
		tc.Fatalf("could not open schema: %s", err)
	}

	return oapi
}
