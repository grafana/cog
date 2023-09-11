package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/txtartest"
	"github.com/stretchr/testify/require"
)

func TestRawTypes_Generate(t *testing.T) {
	test := txtartest.TxTarTest{
		Root: "../../../testdata/jennies/rawtypes",
		Name: "jennies/TypescriptRawTypes",
	}

	jenny := RawTypes{}

	test.Run(t, func(tc *txtartest.Test) {
		req := require.New(tc)

		files, err := jenny.Generate(tc.TypesIR())
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
