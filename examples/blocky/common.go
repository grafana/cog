package main

import (
	"github.com/grafana/cog/generated/common"
	"github.com/grafana/cog/generated/prometheus"
	"github.com/grafana/cog/generated/timeseries"
)

func basicPrometheusQuery(query string, legend string) *prometheus.DataqueryBuilder {
	return prometheus.NewDataqueryBuilder().
		Expr(query).
		LegendFormat(legend)
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
