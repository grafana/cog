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

	req.True(ByObjectName("dashboard", "Dashboard")(ast.Schemas{}, dashboardBuilder))
	req.True(ByObjectName("dashboard", "dashboard")(ast.Schemas{}, dashboardBuilder))
	req.False(ByObjectName("dashboard", "EmptyDashboard")(ast.Schemas{}, dashboardBuilder))
}

func TestByBuilder(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ast.Builder{
		Name: "EmptyDashboard",
		For:  ast.NewObject("dashboard", "Dashboard", ast.NewStruct()),
	}

	req.True(ByName("dashboard", "EmptyDashboard")(ast.Schemas{}, dashboardBuilder))
	req.True(ByName("dashboard", "emptydashboard")(ast.Schemas{}, dashboardBuilder))
	req.False(ByName("dashboard", "Dashboard")(ast.Schemas{}, dashboardBuilder))
}
