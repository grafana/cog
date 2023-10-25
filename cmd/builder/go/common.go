package main

import (
	"github.com/grafana/cog/generated/common/reducedataoptions"
	"github.com/grafana/cog/generated/common/vizlegendoptions"
	gauge "github.com/grafana/cog/generated/gauge/panel"
	logs "github.com/grafana/cog/generated/logs/panel"
	loki "github.com/grafana/cog/generated/loki/dataquery"
	prometheus "github.com/grafana/cog/generated/prometheus/dataquery"
	timeseries "github.com/grafana/cog/generated/timeseries/panel"
	common "github.com/grafana/cog/generated/types/common"
	types "github.com/grafana/cog/generated/types/dashboard"
	lokitypes "github.com/grafana/cog/generated/types/loki"
	promtypes "github.com/grafana/cog/generated/types/prometheus"
)

func toPtr[T any](input T) *T {
	return &input
}

func basicPrometheusQuery(query string, legend string) *promtypes.Dataquery {
	queryBuilder := prometheus.New().
		Expr(query).
		LegendFormat(legend)

	result, err := queryBuilder.Build() // TODO
	if err != nil {
		panic(err)
	}

	return result
}

func basicLokiQuery(query string) *lokitypes.Dataquery {
	queryBuilder := loki.New().
		Expr(query)

	result, err := queryBuilder.Build() // TODO
	if err != nil {
		panic(err)
	}

	return result
}

func tablePrometheusQuery(query string, ref string) *promtypes.Dataquery {
	queryBuilder := prometheus.New().
		Expr(query).
		Instant(true).
		Format(promtypes.PromQueryFormatTable).
		RefId(ref)

	result, err := queryBuilder.Build() // TODO
	if err != nil {
		panic(err)
	}

	return result
}

func defaultTimeseries() *timeseries.Builder {
	return timeseries.New().
		LineWidth(1).
		FillOpacity(10).
		DrawStyle(common.GraphDrawStyleLine).
		ShowPoints(common.VisibilityModeNever).
		Legend(
			vizlegendoptions.New().
				ShowLegend(true).
				Placement(common.LegendPlacementBottom).
				DisplayMode(common.LegendDisplayModeList),
		)
}

func defaultLogs() *logs.Builder {
	return logs.New().
		Datasource(types.DataSourceRef{
			Type: toPtr("loki"),
			Uid:  toPtr("grafanacloud-logs"),
		}).
		ShowTime(true).
		EnableLogDetails(true).
		SortOrder(common.LogsSortOrderDescending).
		WrapLogMessage(true)
}

func defaultGauge() *gauge.Builder {
	return gauge.New().
		Orientation(common.VizOrientationAuto).
		ReduceOptions(reducedataoptions.New().Calcs([]string{"lastNotNull"}).Values(false)) // TODO: not intuitive
}
