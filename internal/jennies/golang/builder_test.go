package golang

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "GoBuilder",
		Skip: map[string]string{
			"builder_delegation_in_disjunction": "disjunctions are eliminated with compiler passes",
		},
	}

	language := New(Config{})
	jenny := Builder{
		Config: Config{
			PackageRoot: "github.com/grafana/cog/generated",
		},
	}

	test.Run(t, func(tc *testutils.Test) {
		var err error
		req := require.New(tc)

		context := tc.BuildersContext()
		context, err = languages.FormatIdentifiers(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
