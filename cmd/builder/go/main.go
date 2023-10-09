package main

/*
import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/dashboard/dashboard"
	"github.com/grafana/cog/generated/dashboard/rowpanel"
	prometheus "github.com/grafana/cog/generated/prometheus/dataquery"
	timeseries "github.com/grafana/cog/generated/timeseries/panel"
	common "github.com/grafana/cog/generated/types/common"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func main() {
	someQuery, err := prometheus.New(
		prometheus.Expr("rate(agent_wal_samples_appended_total{}[10m])"),
	)
	if err != nil {
		panic(err)
	}

	someLogsPanel, err := timeseries.New(
		timeseries.Title("Some timeseries panel"),
		timeseries.Transparent(true),
		timeseries.Description("Let there be logs"),
		timeseries.Decimals(2),
		timeseries.Min(0),
		timeseries.Max(0),
		timeseries.LineWidth(5),
		timeseries.DrawStyle(common.GraphDrawStyleBars),
		timeseries.Targets([]types.Target{
			someQuery.Internal(),
		}),
	)
	if err != nil {
		panic(err)
	}

	overviewPanel, err := rowpanel.New(
		"Overview",
		rowpanel.Panels([]types.Panel{
			*someLogsPanel.Internal(),
		}),
	)
	if err != nil {
		panic(err)
	}

	builder, err := dashboard.New(
		"Some title",
		dashboard.Uid("test-dashboard-codegen"),
		dashboard.Description("Some description"),

		dashboard.Refresh("1m"),
		dashboard.Time("now-3h", "now"),
		dashboard.Timepicker(
			timepicker.RefreshIntervals([]string{"30s", "1m", "5m"}),
		),

		dashboard.Style(types.StyleEnumStyleDark),
		dashboard.Timezone("utc"),
		dashboard.Tooltip(types.DashboardCursorSyncCrosshair),
		dashboard.Tags([]string{"generated", "from", "cue"}),
		dashboard.Links([]types.DashboardLink{
			{
				Title:       "Some link",
				Url:         "http://google.com",
				AsDropdown:  false,
				TargetBlank: true,
			},
		}),

		dashboard.Rows([]types.RowPanel{
			*overviewPanel.Internal(),
		}),
	)
	if err != nil {
		panic(err)
	}

	dashboardJson, err := json.MarshalIndent(builder.Internal(), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dashboardJson))
}
*/
