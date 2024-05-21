package java

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
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

	language := New()
	jenny := RawTypes{}

	test.Run(t, func(tc *testutils.Test) {
		var err error
		req := require.New(tc)

		context := tc.BuildersContext()
		context, err = languages.GenerateBuilderNilChecks(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
