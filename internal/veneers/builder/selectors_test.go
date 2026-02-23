package builder

import (
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/stretchr/testify/require"
)

func TestByObjectName(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ast.Builder{
		Name: "EmptyDashboard",
		For:  ast.NewObject("dashboard", "Dashboard", ast.NewStruct()),
	}

	req.True(ByObjectName("dashboard", "Dashboard").Matches(ast.Schemas{}, dashboardBuilder))
	req.True(ByObjectName("dashboard", "dashboard").Matches(ast.Schemas{}, dashboardBuilder))
	req.False(ByObjectName("dashboard", "EmptyDashboard").Matches(ast.Schemas{}, dashboardBuilder))
}

func TestByBuilder(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ast.Builder{
		Package: "builderpkg",
		Name:    "EmptyDashboard",
		For:     ast.NewObject("dashboard", "Dashboard", ast.NewStruct()),
	}

	req.True(ByName("builderpkg", "EmptyDashboard").Matches(ast.Schemas{}, dashboardBuilder))
	req.True(ByName("builderpkg", "emptydashboard").Matches(ast.Schemas{}, dashboardBuilder))
	req.False(ByName("dashboard", "EmptyDashboard").Matches(ast.Schemas{}, dashboardBuilder))
	req.False(ByName("dashboard", "emptydashboard").Matches(ast.Schemas{}, dashboardBuilder))
	req.False(ByName("dashboard", "Dashboard").Matches(ast.Schemas{}, dashboardBuilder))
}
