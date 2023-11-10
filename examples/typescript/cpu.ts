import {PanelBuilder as TimeseriesPanelBuilder} from "../../generated/timeseries";
import {PanelBuilder as GaugePanelBuilder} from "../../generated/gauge";
import {basicPrometheusQuery, defaultGauge, defaultTimeseries} from "./common";
import {StackingConfigBuilder, StackingMode} from "../../generated/common";
import {ThresholdsConfigBuilder, ThresholdsMode} from "../../generated/dashboard";

export const cpuUsageTimeseries = (): TimeseriesPanelBuilder => {
    const 	query = `(
  (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
/ ignoring(cpu) group_left
  count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
)`;

    return defaultTimeseries()
        .title("CPU Usage")
        .stacking(new StackingConfigBuilder().mode(StackingMode.StackingModeNormal))
        .thresholds(
            new ThresholdsConfigBuilder()
                .mode(ThresholdsMode.ThresholdsModeAbsolute)
                .steps([
                    {value: null, color: "green"},
                    {value: 80.0, color: "red"},
                ])
        )
        .min(0)
        .max(1)
        .unit("percentunit")
        .withTarget(basicPrometheusQuery(query, "{{ cpu }}"));
};

export const loadAverageTimeseries = (): TimeseriesPanelBuilder => {
    return defaultTimeseries()
        .title("Load Average")
        .stacking(new StackingConfigBuilder().mode(StackingMode.StackingModeNormal))
        .thresholds(
            new ThresholdsConfigBuilder()
                .mode(ThresholdsMode.ThresholdsModeAbsolute)
                .steps([
                    {value: null, color: "green"},
                    {value: 80.0, color: "red"},
                ])
        )
        .min(0)
        .unit("short")
        .withTarget(
            basicPrometheusQuery(`node_load1{job="integrations/raspberrypi-node", instance="$instance"}`, "1m load average"),
        )
        .withTarget(
            basicPrometheusQuery(`node_load5{job="integrations/raspberrypi-node", instance="$instance"}`, "5m load average"),
        )
        .withTarget(
            basicPrometheusQuery(`node_load15{job="integrations/raspberrypi-node", instance="$instance"}`, "15m load average"),
        )
        .withTarget(
            basicPrometheusQuery(`count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})`, "logical cores"),
        );
};

export const cpuTemperatureGauge = (): GaugePanelBuilder => {
    return defaultGauge()
        .title("CPU Temperature")
        .min(30)
        .max(100)
        .unit("celsius")
        .thresholds(
            new ThresholdsConfigBuilder()
                .mode(ThresholdsMode.ThresholdsModeAbsolute)
                .steps([
                    {value: null, color: "rgba(50, 172, 45, 0.97)"},
                    {value: 65.0, color: "rgba(237, 129, 40, 0.89)"},
                    {value: 85.0, color: "rgba(245, 54, 54, 0.9)"},
                ])
        )
        .withTarget(
            basicPrometheusQuery(`avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})`, ""),
        );
};
