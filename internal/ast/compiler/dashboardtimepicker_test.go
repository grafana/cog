package compiler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestDashboardTimePicker(t *testing.T) {
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
				)),
				ast.NewObject("dashboard", "Dashboard", ast.NewStruct(
					ast.NewStructField("title", ast.String()),
					ast.NewStructField("timepicker", ast.NewStruct(
						ast.NewStructField("refresh_intervals", ast.NewArray(ast.String())),
					)),
				)),
			},
		},
	}

	// Prepare expected output
	expected := ast.Schemas{
		// Unchanged
		schemas[0],

		// The timepicker is no longer an anonymous struct
		&ast.Schema{
			Package: "dashboard",
			Objects: []ast.Object{
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
			},
		},
	}

	// Run the compiler pass
	runDashboardTimePickerPass(t, schemas, expected)
}

func runDashboardTimePickerPass(t *testing.T, input ast.Schemas, expectedOutput ast.Schemas) {
	t.Helper()

	req := require.New(t)

	compilerPass := &DashboardTimePicker{}
	processedFiles, err := compilerPass.Process(input)
	req.NoError(err)
	req.Len(processedFiles, len(input))
	for i := range input {
		req.Empty(cmp.Diff(expectedOutput[i], processedFiles[i]))
	}
}
