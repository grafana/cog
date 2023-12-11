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

	req.True(ByObjectName("dashboard", "Dashboard")(dashboardBuilder))
	req.True(ByObjectName("dashboard", "dashboard")(dashboardBuilder))
	req.False(ByObjectName("dashboard", "EmptyDashboard")(dashboardBuilder))
}

func TestByBuilder(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ast.Builder{
		Name: "EmptyDashboard",
		For:  ast.NewObject("dashboard", "Dashboard", ast.NewStruct()),
	}

	req.True(ByName("dashboard", "EmptyDashboard")(dashboardBuilder))
	req.True(ByName("dashboard", "emptydashboard")(dashboardBuilder))
	req.False(ByName("dashboard", "Dashboard")(dashboardBuilder))
}
