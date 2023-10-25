package main

import (
	"github.com/grafana/cog/generated/common/stackingconfig"
	"github.com/grafana/cog/generated/dashboard/thresholdsconfig"
	gauge "github.com/grafana/cog/generated/gauge/panel"
	timeseries "github.com/grafana/cog/generated/timeseries/panel"
	common "github.com/grafana/cog/generated/types/common"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func memoryUsageTimeseries() *timeseries.Builder {
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
		Stacking(stackingconfig.New().Mode(common.StackingModeNormal)). // TODO: painful, not intuitive
		Thresholds(thresholdsconfig.New().Mode(types.ThresholdsModeAbsolute).Steps([]types.Threshold{
			{Value: nil, Color: "green"},
			{Value: toPtr(80.0), Color: "red"},
		})).
		Min(0).
		Unit("bytes").
		Decimals(2).
		Targets([]types.Target{
			basicPrometheusQuery(memUsedQuery, "Used"),
			basicPrometheusQuery(`node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Buffers"),
			basicPrometheusQuery(`node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Cached"),
			basicPrometheusQuery(`node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Free"),
		})
}

func memoryUsageGauge() *gauge.Builder {
	query := `100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)`

	return defaultGauge().
		Title("Memory Usage").
		Min(30).
		Max(100).
		Unit("percent").
		Thresholds(thresholdsconfig.New().Mode(types.ThresholdsModeAbsolute).Steps([]types.Threshold{
			{Value: nil, Color: "rgba(50, 172, 45, 0.97)"},
			{Value: toPtr(80.0), Color: "rgba(237, 129, 40, 0.89)"},
			{Value: toPtr(90.0), Color: "rgba(245, 54, 54, 0.9)"},
		})).
		Targets([]types.Target{
			basicPrometheusQuery(query, ""),
		})
}
