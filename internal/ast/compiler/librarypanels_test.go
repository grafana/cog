package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestLibraryPanels_withNoLibraryPanel(t *testing.T) {
	// Prepare test input
	schemas := ast.Schemas{
		&ast.Schema{
			Package: dashboardPackage,
			Objects: []ast.Object{
				ast.NewObject(dashboardPackage, dashboardPanelObject, ast.NewStruct(
					ast.NewStructField("AString", ast.String()),
				)),
			},
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &LibraryPanels{}, schemas, schemas)
}

func TestLibraryPanels_withNoPanel(t *testing.T) {
	// Prepare test input
	schemas := ast.Schemas{
		&ast.Schema{
			Package: libraryPanelPackage,
			Objects: []ast.Object{
				ast.NewObject(libraryPanelObject, libraryPanelObject, ast.NewStruct(
					ast.NewStructField(libraryPanelModelField, ast.Any()),
				)),
			},
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &LibraryPanels{}, schemas, schemas)
}

func TestLibraryPanels_rewrite(t *testing.T) {
	// Prepare test input
	schemas := ast.Schemas{
		&ast.Schema{
			Package: dashboardPackage,
			Objects: []ast.Object{
				ast.NewObject(dashboardPackage, dashboardPanelObject, ast.NewStruct(
					ast.NewStructField(dashboardPanelIDField, ast.String()),
					ast.NewStructField(dashboardPanelGridPosField, ast.Any()),
					ast.NewStructField(dashboardPanelLibraryPanelField, ast.Any()),
					ast.NewStructField(dashboardPanelTypeField, ast.String()),
					ast.NewStructField(dashboardPanelsField, ast.NewArray(ast.Any())),
				)),
			},
		},

		&ast.Schema{
			Package: libraryPanelPackage,
			Objects: []ast.Object{
				ast.NewObject(libraryPanelObject, libraryPanelObject, ast.NewStruct(
					ast.NewStructField("uid", ast.String()),
					ast.NewStructField(libraryPanelModelField, ast.Any()),
				)),
			},
		},
	}

	// Prepare expected output
	expected := ast.Schemas{
		// the dashboard schema is left untouched
		schemas[0],

		&ast.Schema{
			Package: libraryPanelPackage,
			Objects: []ast.Object{
				ast.NewObject(libraryPanelObject, libraryPanelObject, ast.NewStruct(
					ast.NewStructField("uid", ast.String()),
					ast.NewStructField(libraryPanelModelField, ast.NewStruct(
						ast.NewStructField(dashboardPanelTypeField, ast.String()),
						ast.NewStructField(dashboardPanelsField, ast.NewArray(ast.Any())),
					)),
				)),
			},
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &LibraryPanels{}, schemas, expected)
}
