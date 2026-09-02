package csharp

import (
	"testing"

	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/testutils"
	"github.com/stretchr/testify/require"
)

// TestBuilder_Generate verifies that the C# Builder jenny emits the
// expected fluent builder classes for every fixture under
// testdata/jennies/builders/. Goldens land under
// testdata/jennies/builders/<case>/CSharpBuilder/.
func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[languages.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "CSharpBuilder",
	}

	cfg := Config{GenerateBuilders: true}
	cfg.InterpolateParameters(func(input string) string { return input })

	language := New(cfg)
	jenny := Builder{
		config:          language.config,
		tmpl:            initTemplates(language.config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}

	test.Run(t, func(tc *testutils.Test[languages.Context]) {
		req := require.New(tc)

		context := tc.UnmarshalJSONInput(testutils.BuildersContextInputFile)
		context, err := languages.GenerateBuilderNilChecks(language, context)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
