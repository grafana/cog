package compiler

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
)

func TestDashboardTargetsRewrite_withNoRefToTarget(t *testing.T) {
	// Prepare test input
	schemas := ast.Schemas{
		&ast.Schema{
			Package: dashboardPackage,
			Objects: []ast.Object{
				ast.NewObject(dashboardPackage, "RefToStruct", ast.NewRef(dashboardPackage, "AStruct")),

				ast.NewObject(dashboardPackage, "AStruct", ast.NewStruct(
					ast.NewStructField("AString", ast.String()),
				)),
			},
		},

		&ast.Schema{
			Package: "not_dashboard_package",
			Objects: []ast.Object{
				ast.NewObject("not_dashboard_package", "RefToTarget", ast.NewRef("not_dashboard_package", dashboardTargetObject)),

				ast.NewObject("not_dashboard_package", "AStruct", ast.NewStruct(
					ast.NewStructField("AString", ast.String()),
				)),
			},
		},
	}

	// Run the compiler pass
	runPassOnSchemas(t, &DashboardTargetsRewrite{}, schemas, schemas)
}

func TestDashboardTargetsRewrite_withRefToTarget(t *testing.T) {
	// Prepare test input
	schema := &ast.Schema{
		Package: dashboardPackage,
		Objects: []ast.Object{
			ast.NewObject(dashboardPackage, "RefToTarget", ast.NewRef(dashboardPackage, dashboardTargetObject)),

			ast.NewObject(dashboardPackage, "AStruct", ast.NewStruct(
				ast.NewStructField("Targets", ast.NewArray(ast.NewRef(dashboardPackage, dashboardTargetObject))),
				ast.NewStructField("SingleTarget", ast.NewRef(dashboardPackage, dashboardTargetObject)),
				ast.NewStructField("MapOfTargets", ast.NewMap(ast.String(), ast.NewRef(dashboardPackage, dashboardTargetObject))),
				ast.NewStructField("Disjunction", ast.NewDisjunction(ast.Types{
					ast.String(),
					ast.NewRef(dashboardPackage, dashboardTargetObject),
				})),
			)),
		},
	}

	dataqueryComposableSlotType := ast.NewComposableSlot(ast.SchemaVariantDataQuery)

	// Prepare expected output
	expected := &ast.Schema{
		Package: dashboardPackage,
		Objects: []ast.Object{
			ast.NewObject(dashboardPackage, "RefToTarget", dataqueryComposableSlotType),

			ast.NewObject(dashboardPackage, "AStruct", ast.NewStruct(
				ast.NewStructField("Targets", ast.NewArray(dataqueryComposableSlotType)),
				ast.NewStructField("SingleTarget", dataqueryComposableSlotType),
				ast.NewStructField("MapOfTargets", ast.NewMap(ast.String(), dataqueryComposableSlotType)),
				ast.NewStructField("Disjunction", ast.NewDisjunction(ast.Types{
					ast.String(),
					dataqueryComposableSlotType,
				})),
			)),
		},
	}

	// Run the compiler pass
	runPassOnSchema(t, &DashboardTargetsRewrite{}, schema, expected)
}
