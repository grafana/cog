package golang

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[common.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "GoBuilder",
		Skip: map[string]string{
			"builder_delegation_in_disjunction": "disjunctions are eliminated with compiler passes",
		},
	}

	config := Config{
		PackageRoot: "github.com/grafana/cog/generated",
	}
	language := New(config)
	jenny := Builder{Config: config}

	test.Run(t, func(tc *testutils.Test[common.Context]) {
		var err error
		req := require.New(tc)

		context := tc.UnmarshalJSONInput(testutils.BuildersContextInputFile)
		context, err = languages.GenerateBuilderNilChecks(language, context)
		req.NoError(err)
		context, err = languages.FormatIdentifiers(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
