package golang

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../../testdata/jennies/rawtypes",
		Name:         "GoRawTypes",
	}

	jenny := RawTypes{
		Config: Config{
			PackageRoot: "github.com/grafana/cog/generated",
		},
	}
	compilerPasses := New().CompilerPasses()

	test.Run(t, func(tc *testutils.Test) {
		req := require.New(tc)

		// We run the compiler passes defined fo Go since without them, we
		// might not be able to translate some of the IR's semantics into Go.
		// Example: disjunctions.
		processedAsts, err := compilerPasses.Process(ast.Schemas{tc.TypesIR()})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(common.Context{
			Schemas: processedAsts,
		})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
