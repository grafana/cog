package compiler

import (
	"fmt"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestSetDatasourceToDataQuery(t *testing.T) {
	schemas := ast.Schemas{
		&ast.Schema{
			Package: "dashboard",
			Objects: testutils.ObjectsMap(
				ast.NewObject(dashboardPackage, dashboardDatasource, ast.NewStruct(
					ast.NewStructField("Uid", ast.NewScalar(ast.KindString, ast.Nullable())),
					ast.NewStructField("Type", ast.NewScalar(ast.KindString, ast.Nullable())),
				)),
			),
		},
		&ast.Schema{
			Package: "no_dataquery",
			Objects: testutils.ObjectsMap(
				ast.NewObject("no_dataquery", datasourceName, ast.String()),
				ast.NewObject("no_dataquery", dashboardDatasource, ast.String()),
			),
		},
		&ast.Schema{
			Package: "dataquery",
			Metadata: ast.SchemaMeta{
				Variant: ast.SchemaVariantDataQuery,
			},
			Objects: testutils.ObjectsMap(
				ast.NewObject("dataquery", datasourceName, ast.Any()),
				ast.NewObject("dataquery", "Query", ast.NewStruct(
					ast.NewStructField("datasource", ast.NewRef("dataquery", datasourceName)),
					ast.NewStructField("Type", ast.NewScalar(ast.KindString, ast.Nullable())),
				)),
			),
		},
	}

	expected := ast.Schemas{
		schemas[0],
		schemas[1],
		&ast.Schema{
			Package: "dataquery",
			Metadata: ast.SchemaMeta{
				Variant: ast.SchemaVariantDataQuery,
			},
			Objects: testutils.ObjectsMap(
				ast.NewObject("dataquery", "Query", ast.NewStruct(
					ast.NewStructField(
						"datasource",
						ast.NewRef(dashboardPackage, dashboardDatasource),
						ast.PassesTrail(fmt.Sprintf("SetDatasourceToDataquery[%s.%s]", dashboardPackage, dashboardDatasource)),
					),
					ast.NewStructField("Type", ast.NewScalar(ast.KindString, ast.Nullable())),
				)),
			),
		},
	}

	runPassOnSchemas(t, &SetDatasourceToDataquery{}, schemas, expected)
}
