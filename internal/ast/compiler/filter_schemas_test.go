package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestFilterSchemas(t *testing.T) {
	// Prepare test input
	allowedObjects := []ObjectReference{
		{Package: "team", Object: "BigTeam"},
		{Package: "dashboard", Object: "Dashboard"},
	}
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "team",
			Objects: testutils.ObjectsMap(
				ast.NewObject("team", "Team", ast.NewStruct(
					ast.NewStructField("Name", ast.String()),
				)),
				ast.NewObject("team", "BigTeam", ast.NewStruct(
					ast.NewStructField("BigName", ast.String()),
				)),
			),
		},

		&ast.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
				ast.NewObject("dashboard", "Link", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Url", ast.String()),
				)),
				ast.NewObject("dashboard", "Variable", ast.NewStruct(
					ast.NewStructField("Label", ast.String()),
					ast.NewStructField("Foo", ast.String()),
				)),
				ast.NewObject("dashboard", "Panel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String()),
				)),
				ast.NewObject("dashboard", "RowPanel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String(ast.Value("row"))),
					ast.NewStructField("panels", ast.NewArray(
						ast.NewRef("dashboard", "Panel"),
					)),
					ast.NewStructField("links", ast.NewArray(
						ast.NewRef("dashboard", "Link"),
					)),
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
					}))),
				)),
			),
		},
	}

	// Prepare expected output
	expected := ast.Schemas{
		&ast.Schema{
			Package: "team",
			Objects: testutils.ObjectsMap(
				ast.NewObject("team", "BigTeam", ast.NewStruct(
					ast.NewStructField("BigName", ast.String()),
				)),
			),
		},

		&ast.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
				ast.NewObject("dashboard", "Link", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Url", ast.String()),
				)),
				ast.NewObject("dashboard", "Panel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String()),
				)),
				ast.NewObject("dashboard", "RowPanel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
					ast.NewStructField("Type", ast.String(ast.Value("row"))),
					ast.NewStructField("panels", ast.NewArray(
						ast.NewRef("dashboard", "Panel"),
					)),
					ast.NewStructField("links", ast.NewArray(
						ast.NewRef("dashboard", "Link"),
					)),
				)),
				ast.NewObject("dashboard", "Dashboard", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
					ast.NewStructField("panels", ast.NewArray(ast.NewDisjunction(ast.Types{
						ast.NewRef("dashboard", "RowPanel"),
						ast.NewRef("dashboard", "Panel"),
					}))),
				)),
			),
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &FilterSchemas{AllowedObjects: allowedObjects}, schemas, expected)
}
