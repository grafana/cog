package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/cog/plugins"
	"github.com/grafana/cog/generated/go/dashboard"
	"github.com/grafana/cog/generated/go/timeseries"
)

func dashboardBuilder() []byte {
	builder := dashboard.NewDashboardBuilder("[TEST] Node Exporter / Raspberry").
		//Uid("test-dashboard-raspberry").
		Tags([]string{"generated", "raspberrypi-node-integration"}).
		TimeSettings(dashboard.NewTimeSettingsBuilder().
			AutoRefresh("30s").
			AutoRefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
			From("now-30m").
			To("now").
			Timezone("browser"),
		).
		CursorSync(dashboard.DashboardCursorSyncCrosshair).
		Layout(dashboard.NewGridLayoutKindBuilder().
			Spec(
				dashboard.NewGridLayoutSpecBuilder().
					Items([]cog.Builder[dashboard.GridLayoutItemKind]{
						dashboard.NewGridLayoutItemKindBuilder().Spec(
							dashboard.NewGridLayoutItemSpecBuilder().
								X(0).
								Y(0).
								Height(200).
								Width(200).
								Element(dashboard.NewElementReferenceBuilder().Name("somePanel")),
						),
					}),
			),
		).
		Element("somePanel", dashboard.NewPanelKindBuilder().Spec(
			dashboard.NewPanelBuilder().
				Uid("somePanel").
				Title("Some panel").
				Description("veeery descriptive").
				VizConfig(
					timeseries.NewVizConfigKindBuilder().
						Unit("s"),
				),
		))

	sampleDashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}
	dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
	if err != nil {
		panic(err)
	}

	return dashboardJson
}

func main() {
	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	dashboardAsBytes := dashboardBuilder()

	fmt.Println(string(dashboardAsBytes))
}
