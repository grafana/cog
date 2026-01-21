import {basicPrometheusQuery, defaultGauge, defaultTimeseries} from "./common";
import {StackingConfigBuilder, StackingMode} from "../../generated/typescript/src/common";
import {
    PanelBuilder,
    QueryGroupBuilder,
    TargetBuilder,
    ThresholdsConfigBuilder,
    ThresholdsMode
} from "../../generated/typescript/src/dashboardv2beta1";

export const cpuUsageTimeseries = (): PanelBuilder => {
    const query = `(
  (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
/ ignoring(cpu) group_left
  count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
)`;

    return new PanelBuilder()
        .title("CPU Usage")
        .visualization(defaultTimeseries()
            .stacking(new StackingConfigBuilder().mode(StackingMode.Normal))
            .min(0)
            .max(1)
            .unit("percentunit")
            .thresholds(
                new ThresholdsConfigBuilder()
                    .mode(ThresholdsMode.Absolute)
                    .steps([
                        {value: null, color: "green"},
                        {value: 80.0, color: "red"},
                    ])
            )
        )
        .data(
            new QueryGroupBuilder()
                .target(new TargetBuilder().query(basicPrometheusQuery(query, "{{ cpu }}")).refId("A"))
        );
};

export const loadAverageTimeseries = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Load Average")
        .visualization(defaultTimeseries()
            .stacking(new StackingConfigBuilder().mode(StackingMode.Normal))
            .min(0)
            .unit("short")
            .thresholds(
                new ThresholdsConfigBuilder()
                    .mode(ThresholdsMode.Absolute)
                    .steps([
                        {value: null, color: "green"},
                        {value: 80.0, color: "red"},
                    ])
            )
        )
        .data(new QueryGroupBuilder().targets([
            new TargetBuilder().query(basicPrometheusQuery(`node_load1{job="integrations/raspberrypi-node", instance="$instance"}`, "1m load average")).refId("A"),
            new TargetBuilder().query(basicPrometheusQuery(`node_load5{job="integrations/raspberrypi-node", instance="$instance"}`, "5m load average")).refId("B"),
            new TargetBuilder().query(basicPrometheusQuery(`node_load15{job="integrations/raspberrypi-node", instance="$instance"}`, "15m load average")).refId("C"),
            new TargetBuilder().query(basicPrometheusQuery(`count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})`, "logical cores")).refId("D"),
        ]));
};

export const cpuTemperatureGauge = (): PanelBuilder => {
    return new PanelBuilder()
        .title("CPU Temperature")
        .visualization(defaultGauge()
            .min(30)
            .max(100)
            .unit("celsius")
            .thresholds(
                new ThresholdsConfigBuilder()
                    .mode(ThresholdsMode.Absolute)
                    .steps([
                        {value: null, color: "rgba(50, 172, 45, 0.97)"},
                        {value: 65.0, color: "rgba(237, 129, 40, 0.89)"},
                        {value: 85.0, color: "rgba(245, 54, 54, 0.9)"},
                    ])
            )
        )
        .data(
            new QueryGroupBuilder()
                .target(new TargetBuilder().query(basicPrometheusQuery(`avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})`, "")).refId("A"))
        );
};
