package rewrite

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/grafana/cog/internal/logs"
	"github.com/grafana/cog/internal/veneers/builder"
	"github.com/grafana/cog/internal/veneers/option"
	"github.com/grafana/cog/pkg/ir"
	"github.com/stretchr/testify/require"
)

type rewriteTestCase struct {
	description string

	inputBuilders  ir.Builders
	builderRules   []*builder.Rule
	optionRules    []option.Rule
	outputBuilders ir.Builders
}

func testData() []rewriteTestCase {
	return []rewriteTestCase{
		{
			description:    "no rewrite rules",
			inputBuilders:  ir.Builders{dashboardBuilder(), panelBuilder()},
			builderRules:   nil,
			optionRules:    nil,
			outputBuilders: ir.Builders{dashboardBuilder(), panelBuilder()},
		},

		{
			description:   "omit an entire builder",
			inputBuilders: ir.Builders{dashboardBuilder(), panelBuilder()},
			builderRules: []*builder.Rule{
				builder.Omit(builder.ByObjectName("test_pkg", "Dashboard")),
			},
			optionRules:    nil,
			outputBuilders: ir.Builders{panelBuilder()},
		},

		{
			description:   "rename single option in single builder",
			inputBuilders: ir.Builders{dashboardBuilder(), panelBuilder()},
			builderRules:  nil,
			optionRules: []option.Rule{
				option.Rename(
					option.ByName("test_pkg", "Panel", "type"),
					"kind",
				),
			},
			outputBuilders: ir.Builders{
				dashboardBuilder(),
				{
					Package: "test_pkg",
					For: ir.NewObject(
						"test_pkg",
						"Panel",
						ir.NewStruct(
							ir.NewStructField("id", ir.NewScalar(ir.KindInt64)),
							ir.NewStructField("type", ir.String()),
						),
					),
					Options: []ir.Option{
						{
							Name: "id",
							Args: []ir.Argument{
								{Name: "id", Type: ir.NewScalar(ir.KindInt64)},
							},
							Assignments: []ir.Assignment{
								ir.ArgumentAssignment(
									ir.Path{{Identifier: "id", Type: ir.NewScalar(ir.KindInt64)}},
									ir.Argument{Name: "id", Type: ir.NewScalar(ir.KindInt64)},
								),
							},
						},
						{
							Name: "kind",
							Args: []ir.Argument{
								{Name: "type", Type: ir.String()},
							},
							Assignments: []ir.Assignment{
								ir.ArgumentAssignment(
									ir.Path{{Identifier: "type", Type: ir.String()}},
									ir.Argument{Name: "type", Type: ir.String()},
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
			inputBuilders: ir.Builders{dashboardBuilder(), panelBuilder()},
			builderRules:  nil,
			optionRules: []option.Rule{
				option.Omit(
					option.ByName("test_pkg", "Dashboard", "title"),
				),
			},
			outputBuilders: ir.Builders{
				{
					Package: "test_pkg",
					For: ir.NewObject(
						"test_pkg",
						"Dashboard",
						ir.NewStruct(
							ir.NewStructField("uid", ir.String()),
							ir.NewStructField("title", ir.String()),
						),
					),
					Options: []ir.Option{
						{
							Name: "uid",
							Args: []ir.Argument{
								{Name: "uid", Type: ir.String()},
							},
							Assignments: []ir.Assignment{
								ir.ArgumentAssignment(
									ir.Path{{Identifier: "uid", Type: ir.String()}},
									ir.Argument{Name: "uid", Type: ir.String()},
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

func dashboardBuilder() ir.Builder {
	return ir.Builder{
		Package: "test_pkg",
		For: ir.NewObject(
			"test_pkg",
			"Dashboard",
			ir.NewStruct(
				ir.NewStructField("uid", ir.String()),
				ir.NewStructField("title", ir.String()),
			),
		),
		Options: []ir.Option{
			{
				Name: "uid",
				Args: []ir.Argument{
					{Name: "uid", Type: ir.String()},
				},
				Assignments: []ir.Assignment{
					ir.ArgumentAssignment(
						ir.Path{{Identifier: "uid", Type: ir.String()}},
						ir.Argument{Name: "uid", Type: ir.String()},
					),
				},
			},
			{
				Name: "title",
				Args: []ir.Argument{
					{Name: "title", Type: ir.String()},
				},
				Assignments: []ir.Assignment{
					ir.ArgumentAssignment(
						ir.Path{{Identifier: "title", Type: ir.String()}},
						ir.Argument{Name: "title", Type: ir.String()},
					),
				},
			},
		},
	}
}

func panelBuilder() ir.Builder {
	return ir.Builder{
		Package: "test_pkg",
		For: ir.NewObject(
			"test_pkg",
			"Panel",
			ir.NewStruct(
				ir.NewStructField("id", ir.NewScalar(ir.KindInt64)),
				ir.NewStructField("type", ir.String()),
			),
		),
		Options: []ir.Option{
			{
				Name: "id",
				Args: []ir.Argument{
					{Name: "id", Type: ir.NewScalar(ir.KindInt64)},
				},
				Assignments: []ir.Assignment{
					ir.ArgumentAssignment(
						ir.Path{{Identifier: "id", Type: ir.NewScalar(ir.KindInt64)}},
						ir.Argument{Name: "id", Type: ir.NewScalar(ir.KindInt64)},
					),
				},
			},
			{
				Name: "type",
				Args: []ir.Argument{
					{Name: "type", Type: ir.String()},
				},
				Assignments: []ir.Assignment{
					ir.ArgumentAssignment(
						ir.Path{{Identifier: "type", Type: ir.String()}},
						ir.Argument{Name: "type", Type: ir.String()},
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

			rewriter := NewRewrite(slog.New(logs.NoopHandler()), []RuleSet{
				{
					Languages:    []string{AllLanguages},
					BuilderRules: tc.builderRules,
					OptionRules:  tc.optionRules,
				},
			}, Config{Debug: false})

			// save our original/expected states
			originalBuildersJSONBeforeApply := mustMarshalJSON(t, tc.inputBuilders)
			expectedBuildersJSON := mustMarshalJSON(t, tc.outputBuilders)

			// apply the rewrite rules
			rewrittenBuilders, err := rewriter.ApplyTo(ir.Schemas{}, tc.inputBuilders, "go")
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
