package python

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/txtartest"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := txtartest.TxTarTest{
		Root: "../../../testdata/jennies/rawtypes",
		Name: "jennies/PythonRawTypes",
		Skip: map[string]string{
			"jennies/rawtypes/intersections": "Intersections are not implemented",
		},
	}

	jenny := RawTypes{}
	compilerPasses := New().CompilerPasses()

	test.Run(t, func(tc *txtartest.Test) {
		req := require.New(tc)

		// We run the compiler passes defined fo Python since without them, we
		// might not be able to translate some of the IR's semantics into Python.
		// Example: anonymous objects.
		processedAsts, err := compilerPasses.Process(ast.Schemas{tc.TypesIR()})
		req.NoError(err)

		req.Len(processedAsts, 1, "we somehow got more ast.Schema than we put in")

		files, err := jenny.Generate(common.Context{Schemas: processedAsts})
		req.NoError(err)

		tc.WriteFiles(files)
	})
}