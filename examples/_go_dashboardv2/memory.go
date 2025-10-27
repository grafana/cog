package main

import (
	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/common"
	dashboard "github.com/grafana/cog/generated/go/dashboardv2beta1"
	"github.com/grafana/cog/generated/go/units"
)

func memoryUsageTimeseries() *dashboard.PanelBuilder {
	memUsedQuery := `(
  node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}
)`

	return dashboard.NewPanelBuilder().
		Title("Memory Usage").
		Visualization(
			defaultTimeseries().
				Stacking(common.NewStackingConfigBuilder().Mode(common.StackingModeNormal)). // TODO: painful, not intuitive
				Thresholds(
					dashboard.NewThresholdsConfigBuilder().
						Mode(dashboard.ThresholdsModeAbsolute).
						Steps([]dashboard.Threshold{
							// TODO: shouldn't value be nullable?
							{Value: 0, Color: "green"},
							{Value: 80.0, Color: "red"},
						}),
				).
				Min(0).
				Unit(units.BytesIEC).
				Decimals(2),
		).
		Data(
			dashboard.NewQueryGroupBuilder().
				Targets([]cog.Builder[dashboard.PanelQueryKind]{
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(memUsedQuery, "Used")).RefId("A"),
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Buffers")).RefId("B"),
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Cached")).RefId("C"),
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Free")).RefId("D"),
				}),
		)
}

func memoryUsageGauge() *dashboard.PanelBuilder {
	query := `100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)`

	return dashboard.NewPanelBuilder().
		Title("Memory Usage").
		Visualization(
			defaultGauge().
				Min(30).
				Max(100).
				Unit(units.Percent).
				Thresholds(
					dashboard.NewThresholdsConfigBuilder().
						Mode(dashboard.ThresholdsModeAbsolute).
						Steps([]dashboard.Threshold{
							// TODO: shouldn't value be nullable?
							{Value: 0, Color: "rgba(50, 172, 45, 0.97)"},
							{Value: 65.0, Color: "rgba(237, 129, 40, 0.89)"},
							{Value: 85.0, Color: "rgba(245, 54, 54, 0.9)"},
						}),
				),
		).
		Data(
			dashboard.NewQueryGroupBuilder().
				Target(dashboard.NewTargetBuilder().Query(basicPrometheusQuery(query, "")).RefId("A")),
		)
}
