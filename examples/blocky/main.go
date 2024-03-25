package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/cog"
	"github.com/grafana/cog/generated/cog/plugins"
	"github.com/grafana/cog/generated/common"
	"github.com/grafana/cog/generated/dashboard"
	"github.com/grafana/cog/generated/prometheus"
	"github.com/grafana/cog/generated/stat"
	"github.com/grafana/cog/generated/timeseries"
)

func reqPerClientTimeseries() *timeseries.PanelBuilder {
	query := `sum by (client) (rate(blocky_query_total[$__rate_interval])) * 60`

	return defaultTimeseries().
		Title("Request rate per client").
		Span(12).
		Unit("reqpm").
		WithTarget(basicPrometheusQuery(query, "{{ client }}"))
}

func reqPerTypeTimeseries() *timeseries.PanelBuilder {
	query := `sum by (type) (rate(blocky_query_total[$__rate_interval])) * 60`

	return defaultTimeseries().
		Title("Request rate per type").
		Span(12).
		Unit("reqpm").
		WithTarget(basicPrometheusQuery(query, "{{ type }}"))
}

func versionStat() *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title("Version").
		Span(6).
		ReduceOptions(
			common.NewReduceDataOptionsBuilder().
				Calcs([]string{"lastNotNull"}).
				Fields("/^version$/"),
		).
		WithTarget(
			prometheus.NewDataqueryBuilder().
				Expr("blocky_build_info").
				Format(prometheus.PromQueryFormatTable).
				Instant(true),
		)
}

func stateStat() *stat.PanelBuilder {
	return stat.NewPanelBuilder().
		Title("State").
		Span(6).
		ReduceOptions(
			common.NewReduceDataOptionsBuilder().
				Calcs([]string{"lastNotNull"}),
		).
		WithTarget(
			prometheus.NewDataqueryBuilder().
				Expr("sum(up{job=~\"$job\"})").
				Format(prometheus.PromQueryFormatTable).
				Instant(true),
		)
}

func dashboardBuilder() []byte {
	builder := dashboard.NewDashboardBuilder("[TEST] blocky").
		Uid("test-dashboard-blocky").
		Refresh("30s").
		Time("now-30m", "now").
		Timezone(common.TimeZoneBrowser).
		Timepicker(
			dashboard.NewTimePickerBuilder().
				RefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
				TimeOptions([]string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"}),
		).
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		WithRow(dashboard.NewRowBuilder("Overview")).
		WithPanel(versionStat()).
		WithPanel(stateStat()).
		WithRow(dashboard.NewRowBuilder("Activity")).
		WithPanel(reqPerClientTimeseries()).
		WithPanel(reqPerTypeTimeseries()).
		WithRow(
			dashboard.NewRowBuilder("Hidden").
				Collapsed(true).
				Panels([]cog.Builder[dashboard.Panel]{
					versionStat(),
					stateStat(),
				}),
		).
		WithRow(dashboard.NewRowBuilder("More Overview")).
		WithPanel(versionStat()).
		WithPanel(stateStat())

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
