package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/txtartest"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Generate(t *testing.T) {
	test := txtartest.TxTarTest{
		Root: "../../../testdata/jennies/builders",
		Name: "jennies/TypescriptBuilder",
	}

	jenny := Builder{}

	test.Run(t, func(tc *txtartest.Test) {
		req := require.New(tc)

		files, err := jenny.Generate(tc.BuildersContext())
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
