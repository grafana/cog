package openapi

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestSchema_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "OpenAPI",
	}

	config := Config{debug: true}
	jenny := Schema{Config: config}
	compilerPasses := New(config).CompilerPasses()

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined fo OpenAPI since without them, we
		// might not be able to translate some of the IR's semantics.
		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := compilerPasses.Process(ast.Schemas{&schema})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(languages.Context{
			Schemas: processedAsts,
		})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
