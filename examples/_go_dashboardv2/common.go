package main

import (
	"github.com/grafana/cog/generated/go/common"
	"github.com/grafana/cog/generated/go/gauge"
	"github.com/grafana/cog/generated/go/logs"
	"github.com/grafana/cog/generated/go/loki"
	"github.com/grafana/cog/generated/go/prometheus"
	"github.com/grafana/cog/generated/go/timeseries"
)

func basicPrometheusQuery(query string, legend string) *prometheus.QueryBuilder {
	return prometheus.NewQueryBuilder().
		Expr(query).
		LegendFormat(legend)
}

func basicLokiQuery(query string) *loki.QueryBuilder {
	return loki.NewQueryBuilder().Expr(query)
}

func tablePrometheusQuery(query string) *prometheus.QueryBuilder {
	return prometheus.NewQueryBuilder().
		Expr(query).
		Instant().
		Format(prometheus.PromQueryFormatTable)
}

func defaultTimeseries() *timeseries.VisualizationBuilder {
	return timeseries.NewVisualizationBuilder().
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

func defaultLogs() *logs.VisualizationBuilder {
	return logs.NewVisualizationBuilder().
		ShowTime(true).
		EnableLogDetails(true).
		SortOrder(common.LogsSortOrderDescending).
		WrapLogMessage(true)
}

func defaultGauge() *gauge.VisualizationBuilder {
	return gauge.NewVisualizationBuilder().
		Orientation(common.VizOrientationAuto).
		// TODO: not intuitive
		ReduceOptions(
			common.NewReduceDataOptionsBuilder().
				Calcs([]string{"lastNotNull"}).
				Values(false),
		)
}
