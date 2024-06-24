package main

import (
	"github.com/grafana/cog/generated/go/common"
	"github.com/grafana/cog/generated/go/dashboard"
	"github.com/grafana/cog/generated/go/gauge"
	"github.com/grafana/cog/generated/go/logs"
	"github.com/grafana/cog/generated/go/loki"
	"github.com/grafana/cog/generated/go/prometheus"
	"github.com/grafana/cog/generated/go/timeseries"
)

func toPtr[T any](input T) *T {
	return &input
}

func basicPrometheusQuery(query string, legend string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(query).
		LegendFormat(legend)
}

func basicLokiQuery(query string) *loki.DataqueryBuilder {
	return loki.NewDataqueryBuilder().
		Expr(query)
}

func tablePrometheusQuery(query string, ref string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(query).
		Instant().
		Format(prometheus.PromQueryFormatTable).
		RefId(ref)
}

func defaultTimeseries() *timeseries.PanelBuilder {
	return timeseries.NewPanelBuilder().
		LineWidth(1).
		FillOpacity(10).
		DrawStyle(common.GraphDrawStyleLine).
		ShowPoints(common.VisibilityModeNever).
		Legend(
			common.NewVizLegendOptionsBuilder().
				ShowLegend(true).
				Placement(common.LegendPlacementBottom).
				DisplayMode(common.LegendDisplayModeList),
		)
}

func defaultLogs() *logs.PanelBuilder {
	return logs.NewPanelBuilder().
		Span(24).
		Datasource(dashboard.DataSourceRef{
			Type: toPtr("loki"),
			Uid:  toPtr("grafanacloud-logs"),
		}).
		ShowTime(true).
		EnableLogDetails(true).
		SortOrder(common.LogsSortOrderDescending).
		WrapLogMessage(true)
}

func defaultGauge() *gauge.PanelBuilder {
	return gauge.NewPanelBuilder().
		Orientation(common.VizOrientationAuto).
		// TODO: not intuitive
		ReduceOptions(
			common.NewReduceDataOptionsBuilder().
				Calcs([]string{"lastNotNull"}).Values(false),
		)
}
