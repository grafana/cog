package python

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "PythonBuilder",
		Skip: map[string]string{
			"anonymous_struct": "Anonymous structs are not supported in Python",
		},
	}

	jenny := Builder{}

	test.Run(t, func(tc *testutils.Test) {
		req := require.New(tc)

		files, err := jenny.Generate(tc.BuildersContext())
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
