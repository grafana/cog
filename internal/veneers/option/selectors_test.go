package option

import (
	"testing"

	"github.com/grafana/cog/internal/ir"
	"github.com/stretchr/testify/require"
)

func TestByName(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ir.Builder{
		For: ir.NewObject("dashboard", "Dashboard", ir.NewStruct()),
	}
	panelBuilder := ir.Builder{
		For: ir.NewObject("dashboard", "Panel", ir.NewStruct()),
	}
	options := []ir.Option{
		{Name: "Editable"},
		{Name: "Refresh"},
		{Name: "TimePicker"},
	}

	singleSelector := ByName("dashboard", "Dashboard", "Refresh")
	notFoundSelector := ByName("dashboard", "Dashboard", "notFound")

	selectedForDashboard := filter(singleSelector, dashboardBuilder, options)
	selectedForPanel := filter(singleSelector, panelBuilder, options)

	req.Len(selectedForDashboard, 1)
	req.Equal("Refresh", selectedForDashboard[0].Name)

	req.Len(selectedForPanel, 0)

	req.Len(filter(notFoundSelector, dashboardBuilder, options), 0)
}

func TestByName_withSeveralOptions(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ir.Builder{
		For: ir.NewObject("dashboard", "Dashboard", ir.NewStruct()),
	}
	panelBuilder := ir.Builder{
		For: ir.NewObject("dashboard", "Panel", ir.NewStruct()),
	}
	options := []ir.Option{
		{Name: "Editable"},
		{Name: "Refresh"},
		{Name: "TimePicker"},
	}

	multiSelector := ByName("dashboard", "Dashboard", "Refresh", "timepicker")
	notFoundSelector := ByName("dashboard", "Dashboard", "NotFound", "Editable")

	selectedForDashboard := filter(multiSelector, dashboardBuilder, options)
	selectedForPanel := filter(notFoundSelector, panelBuilder, options)

	req.Len(selectedForDashboard, 2)
	req.Equal("Refresh", selectedForDashboard[0].Name)
	req.Equal("TimePicker", selectedForDashboard[1].Name)

	req.Len(selectedForPanel, 0)

	selectedForDashboardWithNotFound := filter(notFoundSelector, dashboardBuilder, options)
	req.Len(selectedForDashboardWithNotFound, 1)
	req.Equal("Editable", selectedForDashboardWithNotFound[0].Name)
}

func TestByBuilder(t *testing.T) {
	req := require.New(t)

	dashboardBuilder := ir.Builder{
		Name:    "Panel",
		Package: "heatmap",
		For:     ir.NewObject("dashboard", "Panel", ir.NewStruct()),
	}
	options := []ir.Option{
		{Name: "Editable"},
		{Name: "Refresh"},
		{Name: "TimePicker"},
	}

	singleSelector := ByBuilder("heatmap", "Panel", "Refresh")
	notFoundSelector := ByBuilder("dashboard", "Panel", "Refresh")

	selectedForDashboard := filter(singleSelector, dashboardBuilder, options)
	req.Len(selectedForDashboard, 1)
	req.Equal("Refresh", selectedForDashboard[0].Name)

	req.Len(filter(notFoundSelector, dashboardBuilder, options), 0)
}

func filter(selector *Selector, builder ir.Builder, opts []ir.Option) []ir.Option {
	var selected []ir.Option

	for _, opt := range opts {
		if selector.Matches(builder, opt) {
			selected = append(selected, opt)
		}
	}

	return selected
}
