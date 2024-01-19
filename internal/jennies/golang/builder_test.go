package golang

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "GoBuilder",
	}

	jenny := Builder{
		Config: Config{
			PackageRoot: "github.com/grafana/cog/generated",
		},
	}

	test.Run(t, func(tc *testutils.Test) {
		req := require.New(tc)

		files, err := jenny.Generate(tc.BuildersContext())
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
