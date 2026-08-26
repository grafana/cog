package builder

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
	"github.com/stretchr/testify/require"
)

func TestByObjectName(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ir.Builder{
		Name: "EmptyDashboard",
		For:  ir.NewObject("dashboard", "Dashboard", ir.NewStruct()),
	}

	req.True(ByObjectName("dashboard", "Dashboard").Matches(ir.Schemas{}, dashboardBuilder))
	req.True(ByObjectName("dashboard", "dashboard").Matches(ir.Schemas{}, dashboardBuilder))
	req.False(ByObjectName("dashboard", "EmptyDashboard").Matches(ir.Schemas{}, dashboardBuilder))
}

func TestByBuilder(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ir.Builder{
		Package: "builderpkg",
		Name:    "EmptyDashboard",
		For:     ir.NewObject("dashboard", "Dashboard", ir.NewStruct()),
	}

	req.True(ByName("builderpkg", "EmptyDashboard").Matches(ir.Schemas{}, dashboardBuilder))
	req.True(ByName("builderpkg", "emptydashboard").Matches(ir.Schemas{}, dashboardBuilder))
	req.False(ByName("dashboard", "EmptyDashboard").Matches(ir.Schemas{}, dashboardBuilder))
	req.False(ByName("dashboard", "emptydashboard").Matches(ir.Schemas{}, dashboardBuilder))
	req.False(ByName("dashboard", "Dashboard").Matches(ir.Schemas{}, dashboardBuilder))
}
