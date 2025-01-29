package golang

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/languages"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/testutils"
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
	language := New(config)
	jenny := Builder{
		Config:          config,
		Tmpl:            initTemplates(config, common.NewAPIReferenceCollector()),
		apiRefCollector: common.NewAPIReferenceCollector(),
	}

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
		input    ast.Type
		expected string
	}{
		{
			desc:     "map",
			context:  languages.Context{},
			input:    ast.NewMap(ast.String(), ast.String()),
			expected: "map[string]string{}",
		},
		{
			desc:     "array",
			context:  languages.Context{},
			input:    ast.NewArray(ast.String()),
			expected: "[]string{}",
		},
		{
			desc: "ref",
			context: languages.Context{
				Schemas: []*ast.Schema{
					{
						Package: "somePkg",
						Objects: orderedmap.FromMap(map[string]ast.Object{
							"SomeType": ast.NewObject("somePkg", "SomeType", ast.NewStruct( /* the fields don't actually matter here */ )),
						}),
					},
				},
			},
			input:    ast.NewRef("somePkg", "SomeType"),
			expected: "somePkg.NewSomeType()",
		},
		{
			desc:    "struct",
			context: languages.Context{},
			input:   ast.NewStruct(ast.NewStructField("field", ast.String())),
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
