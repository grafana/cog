package java

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuidlers_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "JavaBuilders",
		Skip:         map[string]string{
			// "initialization_safeguards": "",
		},
	}

	jenny := RawTypes{}

	test.Run(t, func(tc *testutils.Test) {
		req := require.New(tc)

		files, err := jenny.Generate(tc.BuildersContext())
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
