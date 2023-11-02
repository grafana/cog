package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestDashboardPanelsRewrite(t *testing.T) {
	// Prepare test input
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "team",
			Objects: []ast.Object{
				ast.NewObject("team", "Team", ast.NewStruct(
					ast.NewStructField("Name", ast.String()),
				)),
			},
		},

		&ast.Schema{
			Package: "dashboard",
			Objects: []ast.Object{
				ast.NewObject("dashboard", "Panel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String()),
				)),
				ast.NewObject("dashboard", "RowPanel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String(ast.Value("row"))),
					ast.NewStructField("panels", ast.NewArray(ast.NewDisjunction(ast.Types{
						ast.NewRef("dashboard", "Panel"),
						ast.NewRef("dashboard", "GraphPanel"),
					}))),
				)),
				ast.NewObject("dashboard", "GraphPanel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String(ast.Value("graph"))),
				)),
				ast.NewObject("dashboard", "Dashboard", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
					ast.NewStructField("panels", ast.NewArray(ast.NewDisjunction(ast.Types{
						ast.NewRef("dashboard", "RowPanel"),
						ast.NewRef("dashboard", "Panel"),
						ast.NewRef("dashboard", "GraphPanel"),
					}))),
				)),
			},
		},
	}

	// Prepare expected output
	expected := ast.Schemas{
		// Unchanged
		schemas[0],

		// The panels field are rewritten for RowPanel and Dashboard
		&ast.Schema{
			Package: "dashboard",
			Objects: []ast.Object{
				ast.NewObject("dashboard", "Panel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String()),
				)),
				ast.NewObject("dashboard", "RowPanel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String(ast.Value("row"))),
					ast.NewStructField("panels", ast.NewArray(ast.NewRef("dashboard", "Panel"))),
				)),
				ast.NewObject("dashboard", "GraphPanel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String(ast.Value("graph"))),
				)),
				ast.NewObject("dashboard", "Dashboard", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
					ast.NewStructField("panels", ast.NewArray(ast.NewDisjunction(ast.Types{
						ast.NewRef("dashboard", "Panel"),
						ast.NewRef("dashboard", "RowPanel"),
					}, ast.Discriminator("type", map[string]string{
						"row":                     "RowPanel",
						ast.DiscriminatorCatchAll: "Panel",
					})))),
				)),
			},
		},
	}

	// Run the compiler pass
	runDashboardPanelsRewritePass(t, schemas, expected)
}

func runDashboardPanelsRewritePass(t *testing.T, input ast.Schemas, expectedOutput ast.Schemas) {
	t.Helper()

	req := require.New(t)

	compilerPass := &DashboardPanelsRewrite{}
	processedFiles, err := compilerPass.Process(input)
	req.NoError(err)
	req.Len(processedFiles, len(input))
	for i := range input {
		req.Empty(cmp.Diff(expectedOutput[i], processedFiles[i]))
	}
}
