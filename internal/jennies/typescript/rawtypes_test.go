package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[ast.Schema]{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "TypescriptRawTypes",
	}

	language := New(Config{})
	jenny := RawTypes{}

	test.Run(t, func(tc *testutils.Test[ast.Schema]) {
		req := require.New(tc)

		// We run the compiler passes defined fo Go since without them, we
		// might not be able to translate some of the IR's semantics into TS.
		// Example: disjunctions.
		schema := tc.UnmarshalJSONInput(testutils.RawTypesIRInputFile)
		processedAsts, err := language.CompilerPasses().Process(ast.Schemas{&schema})
		req.NoError(err)
		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		context := common.Context{Schemas: processedAsts}
		context, err = languages.FormatIdentifiers(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
