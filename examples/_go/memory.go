package main

import (
	"github.com/grafana/cog/generated/common"
	"github.com/grafana/cog/generated/dashboard"
	"github.com/grafana/cog/generated/gauge"
	"github.com/grafana/cog/generated/timeseries"
)

func memoryUsageTimeseries() *timeseries.PanelBuilder {
	memUsedQuery := `(
  node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}
)`

	return defaultTimeseries().
		Title("Memory Usage").
		Stacking(common.NewStackingConfigBuilder().Mode(common.StackingModeNormal)). // TODO: painful, not intuitive
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{
					{Value: nil, Color: "green"},
					{Value: toPtr(80.0), Color: "red"},
				}),
		).
		Min(0).
		Unit("bytes").
		Decimals(2).
		Targets([]dashboard.Target{
			basicPrometheusQuery(memUsedQuery, "Used"),
			basicPrometheusQuery(`node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Buffers"),
			basicPrometheusQuery(`node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Cached"),
			basicPrometheusQuery(`node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Free"),
		})
}

func memoryUsageGauge() *gauge.PanelBuilder {
	query := `100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)`

	return defaultGauge().
		Title("Memory Usage").
		Min(30).
		Max(100).
		Unit("percent").
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{
					{Value: nil, Color: "rgba(50, 172, 45, 0.97)"},
					{Value: toPtr(80.0), Color: "rgba(237, 129, 40, 0.89)"},
					{Value: toPtr(90.0), Color: "rgba(245, 54, 54, 0.9)"},
				}),
		).
		Targets([]dashboard.Target{
			basicPrometheusQuery(query, ""),
		})
}
