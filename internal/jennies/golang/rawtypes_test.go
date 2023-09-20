package golang

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/txtartest"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := txtartest.TxTarTest{
		Root: "../../../testdata/jennies/rawtypes",
		Name: "jennies/GoRawTypes",
	}

	jenny := RawTypes{}
	compilerPasses := CompilerPasses()

	test.Run(t, func(tc *txtartest.Test) {
		req := require.New(tc)

		var err error
		processedAsts := []*ast.File{tc.TypesIR()}

		// We run the compiler passes defined fo Go since without them, we
		// might not be able to translate some of the IR's semantics into Go.
		// Example: disjunctions.
		for _, compilerPass := range compilerPasses {
			processedAsts, err = compilerPass.Process(processedAsts)
			req.NoError(err)
		}

		req.Len(processedAsts, 1, "we somehow got more ast.File than we put in")

		files, err := jenny.Generate(processedAsts[0])
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
