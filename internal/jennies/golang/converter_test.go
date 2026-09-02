package golang

import (
	"testing"

	"github.com/grafana/cog/internal/builders"
	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/apiref"
	"github.com/grafana/cog/pkg/languages"
	"github.com/grafana/cog/pkg/logs"
	"github.com/stretchr/testify/require"
)

func TestConverter_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[languages.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "GoConverter",
		Skip: map[string]string{
			"anonymous_struct":                  "anonymous structs are eliminated with compiler passes",
			"builder_delegation_in_disjunction": "disjunctions are eliminated with compiler passes",
			"dashboard_panel":                   "this test if for Java generics for dashboard.Panel",
		},
	}

	config := Config{
		PackageRoot: "github.com/grafana/cog/generated",
	}
	language := New(logs.NoopLogger(), config)
	jenny := Converter{
		Config:          config,
		NullableConfig:  language.NullableKinds(),
		Tmpl:            initTemplates(config, apiref.NewAPIReferenceCollector()),
		apiRefCollector: apiref.NewAPIReferenceCollector(),
	}

	test.Run(t, func(tc *testutils.Test[languages.Context]) {
		var err error
		req := require.New(tc)

		context := tc.UnmarshalJSONInput(testutils.BuildersContextInputFile)
		context.Builders, err = builders.GenerateNilChecks(language.NullableKinds(), context.Schemas, context.Builders)
		req.NoError(err)

		files, err := jenny.Generate(context)
		req.NoError(err)

		tc.WriteFiles(files)
	})
}
