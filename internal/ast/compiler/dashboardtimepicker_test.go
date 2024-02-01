package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestDashboardTimePicker(t *testing.T) {
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
				)),
				ast.NewObject("dashboard", "Dashboard", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
					ast.NewStructField("timepicker", ast.NewStruct(
						ast.NewStructField("refresh_intervals", ast.NewArray(ast.String())),
					)),
				)),
			),
		},
	}

	// Prepare expected output
	expected := ast.Schemas{
		// Unchanged
		schemas[0],

		// The timepicker is no longer an anonymous struct
		&ast.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
				ast.NewObject("dashboard", "Panel", ast.NewStruct(
					ast.NewStructField("Title", ast.String()),
				)),
				ast.NewObject("dashboard", "Dashboard", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
					ast.NewStructField("timepicker", ast.NewRef("dashboard", "TimePicker")),
				)),
				ast.NewObject("dashboard", "TimePicker", ast.NewStruct(
					ast.NewStructField("refresh_intervals", ast.NewArray(ast.String())),
				)),
			),
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &DashboardTimePicker{}, schemas, expected)
}
