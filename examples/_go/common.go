package main

import (
	"github.com/grafana/cog/generated/common"
	"github.com/grafana/cog/generated/dashboard"
	"github.com/grafana/cog/generated/gauge"
	"github.com/grafana/cog/generated/logs"
	"github.com/grafana/cog/generated/loki"
	"github.com/grafana/cog/generated/prometheus"
	"github.com/grafana/cog/generated/timeseries"
)

func toPtr[T any](input T) *T {
	return &input
}

func basicPrometheusQuery(query string, legend string) *prometheus.Dataquery {
	queryBuilder := prometheus.NewDataqueryBuilder().
		Expr(query).
		LegendFormat(legend)

	result, err := queryBuilder.Build() // TODO
	if err != nil {
		panic(err)
	}

	return result
}

func basicLokiQuery(query string) *loki.Dataquery {
	queryBuilder := loki.NewDataqueryBuilder().
		Expr(query)

	result, err := queryBuilder.Build() // TODO
	if err != nil {
		panic(err)
	}

	return result
}

func tablePrometheusQuery(query string, ref string) *prometheus.Dataquery {
	queryBuilder := prometheus.NewDataqueryBuilder().
		Expr(query).
		Instant(true).
		Format(prometheus.PromQueryFormatTable).
		RefId(ref)

	result, err := queryBuilder.Build() // TODO
	if err != nil {
		panic(err)
	}

	return result
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
