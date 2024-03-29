package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestDashboardPanelsRewrite(t *testing.T) {
	// Prepare test input
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "team",
			Objects: testutils.ObjectsMap(
				ast.NewObject("team", "Team", ast.NewStruct(
					ast.NewStructField("Name", ast.String()),
				)),
			),
		},

		&ast.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
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
			),
		},
	}

	// Prepare expected output
	expected := ast.Schemas{
		// Unchanged
		schemas[0],

		// The panels field are rewritten for RowPanel and Dashboard
		&ast.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
				ast.NewObject("dashboard", "Panel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String()),
				)),
				ast.NewObject("dashboard", "RowPanel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String(ast.Value("row"))),
					ast.NewStructField("panels", ast.NewArray(ast.NewRef("dashboard", "Panel")), ast.PassesTrail("DashboardPanelsRewrite[changed type]")),
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
					}))), ast.PassesTrail("DashboardPanelsRewrite[changed type]")),
				)),
			),
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &DashboardPanelsRewrite{}, schemas, expected)
}
