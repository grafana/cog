package rewrite

import (
	"encoding/json"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
	"github.com/stretchr/testify/require"
)

type rewriteTestCase struct {
	description string

	inputBuilders  ast.Builders
	builderRules   []builder.RewriteRule
	optionRules    []option.RewriteRule
	outputBuilders ast.Builders
}

func testData() []rewriteTestCase {
	return []rewriteTestCase{
		{
			description:    "no rewrite rules",
			inputBuilders:  ast.Builders{dashboardBuilder(), panelBuilder()},
			builderRules:   nil,
			optionRules:    nil,
			outputBuilders: ast.Builders{dashboardBuilder(), panelBuilder()},
		},

		{
			description:   "omit an entire builder",
			inputBuilders: ast.Builders{dashboardBuilder(), panelBuilder()},
			builderRules: []builder.RewriteRule{
				builder.Omit(builder.ByObjectName("test_pkg", "Dashboard")),
			},
			optionRules:    nil,
			outputBuilders: ast.Builders{panelBuilder()},
		},

		{
			description:   "rename single option in single builder",
			inputBuilders: ast.Builders{dashboardBuilder(), panelBuilder()},
			builderRules:  nil,
			optionRules: []option.RewriteRule{
				option.Rename(
					option.ByName("test_pkg", "Panel", "type"),
					"kind",
				),
			},
			outputBuilders: ast.Builders{
				dashboardBuilder(),
				{
					Package: "test_pkg",
					For: ast.NewObject(
						"test_pkg",
						"Panel",
						ast.NewStruct(
							ast.NewStructField("id", ast.NewScalar(ast.KindInt64)),
							ast.NewStructField("type", ast.String()),
						),
					),
					Options: []ast.Option{
						{
							Name: "id",
							Args: []ast.Argument{
								{Name: "id", Type: ast.NewScalar(ast.KindInt64)},
							},
							Assignments: []ast.Assignment{
								ast.ArgumentAssignment(
									ast.Path{{Identifier: "id", Type: ast.NewScalar(ast.KindInt64)}},
									ast.Argument{Name: "id", Type: ast.NewScalar(ast.KindInt64)},
								),
							},
						},
						{
							Name: "kind",
							Args: []ast.Argument{
								{Name: "type", Type: ast.String()},
							},
							Assignments: []ast.Assignment{
								ast.ArgumentAssignment(
									ast.Path{{Identifier: "type", Type: ast.String()}},
									ast.Argument{Name: "type", Type: ast.String()},
								),
							},
							VeneerTrail: []string{"Rename[type → kind]"},
						},
					},
				},
			},
		},

		{
			description:   "omit single option in single builder",
			inputBuilders: ast.Builders{dashboardBuilder(), panelBuilder()},
			builderRules:  nil,
			optionRules: []option.RewriteRule{
				option.Omit(
					option.ByName("test_pkg", "Dashboard", "title"),
				),
			},
			outputBuilders: ast.Builders{
				{
					Package: "test_pkg",
					For: ast.NewObject(
						"test_pkg",
						"Dashboard",
						ast.NewStruct(
							ast.NewStructField("uid", ast.String()),
							ast.NewStructField("title", ast.String()),
						),
					),
					Options: []ast.Option{
						{
							Name: "uid",
							Args: []ast.Argument{
								{Name: "uid", Type: ast.String()},
							},
							Assignments: []ast.Assignment{
								ast.ArgumentAssignment(
									ast.Path{{Identifier: "uid", Type: ast.String()}},
									ast.Argument{Name: "uid", Type: ast.String()},
								),
							},
						},
					},
				},
				panelBuilder(),
			},
		},
	}
}

func dashboardBuilder() ast.Builder {
	return ast.Builder{
		Package: "test_pkg",
		For: ast.NewObject(
			"test_pkg",
			"Dashboard",
			ast.NewStruct(
				ast.NewStructField("uid", ast.String()),
				ast.NewStructField("title", ast.String()),
			),
		),
		Options: []ast.Option{
			{
				Name: "uid",
				Args: []ast.Argument{
					{Name: "uid", Type: ast.String()},
				},
				Assignments: []ast.Assignment{
					ast.ArgumentAssignment(
						ast.Path{{Identifier: "uid", Type: ast.String()}},
						ast.Argument{Name: "uid", Type: ast.String()},
					),
				},
			},
			{
				Name: "title",
				Args: []ast.Argument{
					{Name: "title", Type: ast.String()},
				},
				Assignments: []ast.Assignment{
					ast.ArgumentAssignment(
						ast.Path{{Identifier: "title", Type: ast.String()}},
						ast.Argument{Name: "title", Type: ast.String()},
					),
				},
			},
		},
	}
}

func panelBuilder() ast.Builder {
	return ast.Builder{
		Package: "test_pkg",
		For: ast.NewObject(
			"test_pkg",
			"Panel",
			ast.NewStruct(
				ast.NewStructField("id", ast.NewScalar(ast.KindInt64)),
				ast.NewStructField("type", ast.String()),
			),
		),
		Options: []ast.Option{
			{
				Name: "id",
				Args: []ast.Argument{
					{Name: "id", Type: ast.NewScalar(ast.KindInt64)},
				},
				Assignments: []ast.Assignment{
					ast.ArgumentAssignment(
						ast.Path{{Identifier: "id", Type: ast.NewScalar(ast.KindInt64)}},
						ast.Argument{Name: "id", Type: ast.NewScalar(ast.KindInt64)},
					),
				},
			},
			{
				Name: "type",
				Args: []ast.Argument{
					{Name: "type", Type: ast.String()},
				},
				Assignments: []ast.Assignment{
					ast.ArgumentAssignment(
						ast.Path{{Identifier: "type", Type: ast.String()}},
						ast.Argument{Name: "type", Type: ast.String()},
					),
				},
			},
		},
	}
}

func TestRewriter_ApplyTo(t *testing.T) {
	testCases := testData()

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.description, func(t *testing.T) {
			req := require.New(t)

			rewriter := NewRewrite([]LanguageRules{
				{
					Language:     AllLanguages,
					BuilderRules: tc.builderRules,
					OptionRules:  tc.optionRules,
				},
			}, Config{Debug: false})

			// save our original/expected states
			originalBuildersJSONBeforeApply := mustMarshalJSON(t, tc.inputBuilders)
			expectedBuildersJSON := mustMarshalJSON(t, tc.outputBuilders)

			// apply the rewrite rules
			rewrittenBuilders, err := rewriter.ApplyTo(ast.Schemas{}, tc.inputBuilders, "go")
			req.NoError(err)

			// save the output states
			originalBuildersJSONAfterApply := mustMarshalJSON(t, tc.inputBuilders)
			rewrittenBuildersJSON := mustMarshalJSON(t, rewrittenBuilders)

			// check that everything went fine
			req.JSONEq(originalBuildersJSONBeforeApply, originalBuildersJSONAfterApply, "input builders aren't modified")
			req.JSONEq(expectedBuildersJSON, rewrittenBuildersJSON, "rewrite result is what we expect")
		})
	}
}

func mustMarshalJSON(t *testing.T, input any) string {
	t.Helper()

	req := require.New(t)

	jsonPayload, err := json.Marshal(input)
	req.NoError(err)

	return string(jsonPayload)
}
