package main

import (
	"github.com/grafana/cog/generated/common/stackingconfig"
	"github.com/grafana/cog/generated/dashboard/thresholdsconfig"
	gauge "github.com/grafana/cog/generated/gauge/panel"
	timeseries "github.com/grafana/cog/generated/timeseries/panel"
	common "github.com/grafana/cog/generated/types/common"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func cpuUsageTimeseries() *timeseries.Builder {
	query := `(
  (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
/ ignoring(cpu) group_left
  count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
)`

	return defaultTimeseries().
		Title("CPU Usage").
		Stacking(stackingconfig.New().Mode(common.StackingModeNormal)). // TODO: painful, not intuitive
		Thresholds(thresholdsconfig.New().Mode(types.ThresholdsModeAbsolute).Steps([]types.Threshold{
			{Value: nil, Color: "green"},
			{Value: toPtr(80.0), Color: "red"},
		})).
		Min(0).
		Max(1).
		Unit("percentunit").
		Targets([]types.Target{
			basicPrometheusQuery(query, "{{ cpu }}"),
		})
}

func loadAverageTimeseries() *timeseries.Builder {
	return defaultTimeseries().
		Title("Load Average").
		Thresholds(thresholdsconfig.New().Mode(types.ThresholdsModeAbsolute).Steps([]types.Threshold{
			{Value: nil, Color: "green"},
			{Value: toPtr(80.0), Color: "red"},
		})).
		Min(0).
		Unit("short").
		Targets([]types.Target{
			basicPrometheusQuery(`node_load1{job="integrations/raspberrypi-node", instance="$instance"}`, "1m load average"),
			basicPrometheusQuery(`node_load5{job="integrations/raspberrypi-node", instance="$instance"}`, "5m load average"),
			basicPrometheusQuery(`node_load15{job="integrations/raspberrypi-node", instance="$instance"}`, "15m load average"),
			basicPrometheusQuery(`count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})`, "logical cores"),
		})
}

func cpuTemperatureGauge() *gauge.Builder {
	return defaultGauge().
		Title("CPU Temperature").
		Min(30).
		Max(100).
		Unit("celsius").
		Thresholds(thresholdsconfig.New().Mode(types.ThresholdsModeAbsolute).Steps([]types.Threshold{
			{Value: nil, Color: "rgba(50, 172, 45, 0.97)"},
			{Value: toPtr(65.0), Color: "rgba(237, 129, 40, 0.89)"},
			{Value: toPtr(85.0), Color: "rgba(245, 54, 54, 0.9)"},
		})).
		Targets([]types.Target{
			basicPrometheusQuery(`avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})`, ""),
		})
}
