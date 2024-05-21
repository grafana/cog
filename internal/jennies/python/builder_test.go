package python

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
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

	language := New(Config{})
	jenny := Builder{}

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
