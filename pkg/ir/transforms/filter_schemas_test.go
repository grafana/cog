package transforms

import (
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestFilterSchemas(t *testing.T) {
	// Prepare test input
	allowedObjects := []ObjectReference{
		{Package: "team", Object: "BigTeam"},
		{Package: "dashboard", Object: "Dashboard"},
	}
	schemas := ir.Schemas{
		&ir.Schema{
			Package: "team",
			Objects: testutils.ObjectsMap(
				ir.NewObject("team", "Team", ir.NewStruct(
					ir.NewStructField("Name", ir.String()),
				)),
				ir.NewObject("team", "BigTeam", ir.NewStruct(
					ir.NewStructField("BigName", ir.String()),
				)),
			),
		},

		&ir.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
				ir.NewObject("dashboard", "Link", ir.NewStruct(
					ir.NewStructField("Title", ir.String()),
					ir.NewStructField("Url", ir.String()),
				)),
				ir.NewObject("dashboard", "Variable", ir.NewStruct(
					ir.NewStructField("Label", ir.String()),
					ir.NewStructField("Foo", ir.String()),
				)),
				ir.NewObject("dashboard", "Panel", ir.NewStruct(
					ir.NewStructField("Title", ir.String()),
					ir.NewStructField("Type", ir.String()),
				)),
				ir.NewObject("dashboard", "RowPanel", ir.NewStruct(
					ir.NewStructField("Title", ir.String()),
					ir.NewStructField("Type", ir.String(ir.Value("row"))),
					ir.NewStructField("panels", ir.NewArray(
						ir.NewRef("dashboard", "Panel"),
					)),
					ir.NewStructField("links", ir.NewArray(
						ir.NewRef("dashboard", "Link"),
					)),
				)),
				ir.NewObject("dashboard", "GraphPanel", ir.NewStruct(
					ir.NewStructField("Title", ir.String()),
					ir.NewStructField("Type", ir.String(ir.Value("graph"))),
				)),
				ir.NewObject("dashboard", "Dashboard", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
					ir.NewStructField("panels", ir.NewArray(ir.NewDisjunction(ir.Types{
						ir.NewRef("dashboard", "RowPanel"),
						ir.NewRef("dashboard", "Panel"),
					}))),
				)),
			),
		},
	}

	// Prepare expected output
	expected := ir.Schemas{
		&ir.Schema{
			Package: "team",
			Objects: testutils.ObjectsMap(
				ir.NewObject("team", "BigTeam", ir.NewStruct(
					ir.NewStructField("BigName", ir.String()),
				)),
			),
		},

		&ir.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
				ir.NewObject("dashboard", "Link", ir.NewStruct(
					ir.NewStructField("Title", ir.String()),
					ir.NewStructField("Url", ir.String()),
				)),
				ir.NewObject("dashboard", "Panel", ir.NewStruct(
					ir.NewStructField("Title", ir.String()),
					ir.NewStructField("Type", ir.String()),
				)),
				ir.NewObject("dashboard", "RowPanel", ir.NewStruct(
					ir.NewStructField("Title", ir.String()),
					ir.NewStructField("Type", ir.String(ir.Value("row"))),
					ir.NewStructField("panels", ir.NewArray(
						ir.NewRef("dashboard", "Panel"),
					)),
					ir.NewStructField("links", ir.NewArray(
						ir.NewRef("dashboard", "Link"),
					)),
				)),
				ir.NewObject("dashboard", "Dashboard", ir.NewStruct(
					ir.NewStructField("title", ir.String()),
					ir.NewStructField("panels", ir.NewArray(ir.NewDisjunction(ir.Types{
						ir.NewRef("dashboard", "RowPanel"),
						ir.NewRef("dashboard", "Panel"),
					}))),
				)),
			),
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &FilterSchemas{AllowedObjects: allowedObjects}, schemas, expected)
}
