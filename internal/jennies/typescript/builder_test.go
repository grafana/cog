package typescript

import (
	"testing"

	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[languages.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "TypescriptBuilder",
	}

	language := New(Config{})
	jenny := Builder{tmpl: initTemplates([]string{})}

	test.Run(t, func(tc *testutils.Test[languages.Context]) {
		var err error
		req := require.New(tc)

		context := tc.UnmarshalJSONInput(testutils.BuildersContextInputFile)
		context, err = languages.GenerateBuilderNilChecks(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
