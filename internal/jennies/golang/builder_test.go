package golang

import (
	"testing"

	"github.com/grafana/cog/internal/builders"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/logs"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Generate(t *testing.T) {
	test := testutils.GoldenFilesTestSuite[languages.Context]{
		TestDataRoot: "../../../testdata/jennies/builders",
		Name:         "GoBuilder",
		Skip: map[string]string{
			"builder_delegation_in_disjunction": "disjunctions are eliminated with compiler passes",
			"dashboard_panel":                   "this test if for Java generics for dashboard.Panel",
		},
	}

	config := Config{
		PackageRoot: "github.com/grafana/cog/generated",
	}
	language := New(logs.NoopLogger(), config)
	jenny := Builder{
		Config:          config,
		Tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
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

func TestBuilder_emptyValueForGuard(t *testing.T) {
	config := Config{
		PackageRoot: "github.com/grafana/cog/generated",
	}
	jenny := Builder{
		Config:          config,
		Tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}

	jenny.typeImportMapper = func(pkg string) string {
		return pkg
	}
	imports := NewImportMap(jenny.Config.PackageRoot)

	testCases := []struct {
		desc     string
		context  languages.Context
		input    ir.Type
		expected string
	}{
		{
			desc:     "map",
			context:  languages.Context{},
			input:    ir.NewMap(ir.String(), ir.String()),
			expected: "map[string]string{}",
		},
		{
			desc:     "array",
			context:  languages.Context{},
			input:    ir.NewArray(ir.String()),
			expected: "[]string{}",
		},
		{
			desc: "ref",
			context: languages.Context{
				Schemas: []*ir.Schema{
					{
						Package: "somePkg",
						Objects: orderedmap.FromMap(map[string]ir.Object{
							"SomeType": ir.NewObject("somePkg", "SomeType", ir.NewStruct( /* the fields don't actually matter here */ )),
						}),
					},
				},
			},
			input:    ir.NewRef("somePkg", "SomeType"),
			expected: "somePkg.NewSomeType()",
		},
		{
			desc:    "struct",
			context: languages.Context{},
			input:   ir.NewStruct(ir.NewStructField("field", ir.String())),
			expected: `&struct {
    Field string ` + "`" + `json:"field,omitempty"` + "`" + `
}{}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			req := require.New(t)

			jenny.typeFormatter = builderTypeFormatter(jenny.Config, tc.context, imports, jenny.typeImportMapper)
			jenny.pathFormatter = makePathFormatter(jenny.typeFormatter)

			req.Equal(tc.expected, jenny.emptyValueForGuard(tc.context, tc.input))
		})
	}
}
