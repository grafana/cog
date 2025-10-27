package main

import (
	"github.com/grafana/cog/generated/go/common"
	dashboard "github.com/grafana/cog/generated/go/dashboardv2beta1"
	"github.com/grafana/cog/generated/go/units"
)

func cpuUsageTimeseries() *dashboard.PanelBuilder {
	query := `(
  (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
/ ignoring(cpu) group_left
  count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
)`

	return dashboard.NewPanelBuilder().
		Title("CPU Usage").
		Visualization(
			defaultTimeseries().
				Stacking(common.NewStackingConfigBuilder().Mode(common.StackingModeNormal)). // TODO: painful, not intuitive
				Min(0).
				Max(1).
				Unit(units.PercentUnit).
				Thresholds(
					dashboard.NewThresholdsConfigBuilder().
						Mode(dashboard.ThresholdsModeAbsolute).
						Steps([]dashboard.Threshold{
							// TODO: shouldn't value be nullable?
							{Value: 0, Color: "green"},
							{Value: 80.0, Color: "red"},
						}),
				),
		).
		Data(
			dashboard.NewQueryGroupBuilder().
				Target(
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(query, "{{ cpu }}")).RefId("A"),
				),
		)
}

func loadAverageTimeseries() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Load Average").
		Visualization(
			defaultTimeseries().
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
				Unit(units.Short),
		).
		Data(
			dashboard.NewQueryGroupBuilder().
				Target(
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`node_load1{job="integrations/raspberrypi-node", instance="$instance"}`, "1m load average")).RefId("A"),
				).
				Target(
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`node_load5{job="integrations/raspberrypi-node", instance="$instance"}`, "5m load average")).RefId("B"),
				).
				Target(
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`node_load15{job="integrations/raspberrypi-node", instance="$instance"}`, "15m load average")).RefId("C"),
				).
				Target(
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})`, "logical cores")).RefId("D"),
				),
		)
}

func cpuTemperatureGauge() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("CPU Temperature").
		Visualization(
			defaultGauge().
				Min(30).
				Max(100).
				Unit(units.Celsius).
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
				Target(
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})`, "")).RefId("A"),
				),
		)
}
