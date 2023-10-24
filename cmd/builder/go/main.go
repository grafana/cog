package main

/*
import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/dashboard/dashboard"
	"github.com/grafana/cog/generated/dashboard/rowpanel"
	"github.com/grafana/cog/generated/dashboard/timepicker"
	prometheus "github.com/grafana/cog/generated/prometheus/dataquery"
	timeseries "github.com/grafana/cog/generated/timeseries/panel"
	common "github.com/grafana/cog/generated/types/common"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func main() {
	someQuery := prometheus.New().
		Expr("rate(agent_wal_samples_appended_total{}[10m])").
		LegendFormat("Samples")
	query, err := someQuery.Build()
	if err != nil {
		panic(err)
	}

	someTimeseriesPanel := timeseries.New().
		Title("Some timeseries panel").
		Transparent(true).
		Description("Let there be data").
		Decimals(2).
		AxisSoftMin(0).
		AxisSoftMax(50).
		LineWidth(5).
		DrawStyle(common.GraphDrawStylePoints).
		Targets([]types.Target{
			query,
		})

	builder := dashboard.New("Some title").
		Uid("test-dashboard-codegen").
		Description("Some description").
		Refresh("1m").
		Time("now-3h", "now").
		Timezone("utc").
		Timepicker(
			timepicker.New().RefreshIntervals([]string{"30s", "1m", "5m"}),
		).
		Tooltip(types.DashboardCursorSyncCrosshair).
		Tags([]string{"generated", "from", "go"}).
		Links([]types.DashboardLink{
			{
				Title:       "Some link",
				Url:         "http://google.com",
				AsDropdown:  false,
				TargetBlank: true,
			},
		}).
		WithRow(rowpanel.New("Overview")).
		WithPanel(someTimeseriesPanel)

	sampleDashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}
	dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dashboardJson))
}
*/
