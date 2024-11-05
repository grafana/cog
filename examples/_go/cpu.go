package main

import (
	"github.com/grafana/cog/generated/go/common"
	"github.com/grafana/cog/generated/go/dashboard"
	"github.com/grafana/cog/generated/go/gauge"
	"github.com/grafana/cog/generated/go/timeseries"
)

func cpuUsageTimeseries() *timeseries.PanelBuilder {
	query := `(
  (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
/ ignoring(cpu) group_left
  count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
)`

	return defaultTimeseries().
		Title("CPU Usage").
		Span(18).
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
		Max(1).
		Unit("percentunit").
		WithTarget(
			basicPrometheusQuery(query, "{{ cpu }}"),
		)
}

func loadAverageTimeseries() *timeseries.PanelBuilder {
	return defaultTimeseries().
		Title("Load Average").
		Span(18).
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{
					{Value: nil, Color: "green"},
					{Value: toPtr(80.0), Color: "red"},
				}),
		).
		Min(0).
		Unit("short").
		WithTarget(
			basicPrometheusQuery(`node_load1{job="integrations/raspberrypi-node", instance="$instance"}`, "1m load average"),
		).
		WithTarget(
			basicPrometheusQuery(`node_load5{job="integrations/raspberrypi-node", instance="$instance"}`, "5m load average"),
		).
		WithTarget(
			basicPrometheusQuery(`node_load15{job="integrations/raspberrypi-node", instance="$instance"}`, "15m load average"),
		).
		WithTarget(
			basicPrometheusQuery(`count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})`, "logical cores"),
		)
}

func cpuTemperatureGauge() *gauge.PanelBuilder {
	return defaultGauge().
		Title("CPU Temperature").
		Span(6).
		Min(30).
		Max(100).
		Unit("celsius").
		Thresholds(
			dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{
					{Value: nil, Color: "rgba(50, 172, 45, 0.97)"},
					{Value: toPtr(65.0), Color: "rgba(237, 129, 40, 0.89)"},
					{Value: toPtr(85.0), Color: "rgba(245, 54, 54, 0.9)"},
				}),
		).
		WithTarget(
			basicPrometheusQuery(`avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})`, ""),
		)
}
