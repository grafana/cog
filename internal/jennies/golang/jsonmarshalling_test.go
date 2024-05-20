package golang

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestJSONMarshalling_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "GoJSONMarshalling",
	}

	config := Config{
		PackageRoot: "github.com/grafana/cog/generated",
	}
	language := New(config)
	jenny := JSONMarshalling{
		Config:               config,
		IdentifiersFormatter: language.IdentifiersFormatter(),
	}

	test.Run(t, func(tc *testutils.Test) {
		req := require.New(tc)

		// We run the compiler passes defined fo Go since without them, we
		// might not be able to translate some of the IR's semantics into Go.
		// Example: disjunctions.
		processedAsts, err := language.CompilerPasses().Process(ast.Schemas{tc.TypesIR()})
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
