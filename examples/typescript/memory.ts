import {PanelBuilder as TimeseriesPanelBuilder} from "../../generated/typescript/src/timeseries";
import {PanelBuilder as GaugePanelBuilder} from "../../generated/typescript/src/gauge";
import {basicPrometheusQuery, defaultGauge, defaultTimeseries} from "./common";
import {StackingConfigBuilder, StackingMode} from "../../generated/typescript/src/common";
import {ThresholdsConfigBuilder, ThresholdsMode} from "../../generated/typescript/src/dashboard";

export const memoryUsageTimeseries = (): TimeseriesPanelBuilder => {
    const 	memUsedQuery = `(
  node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}
)`;

    return defaultTimeseries()
        .title("Memory Usage")
        .span(18)
        .stacking(new StackingConfigBuilder().mode(StackingMode.Normal))
        .thresholds(
            new ThresholdsConfigBuilder()
                .mode(ThresholdsMode.Absolute)
                .steps([
                    {value: null, color: "green"},
                    {value: 80.0, color: "red"},
                ])
        )
        .min(0)
        .unit("bytes")
        .withTarget(basicPrometheusQuery(memUsedQuery, "Used"))
        .withTarget(
            basicPrometheusQuery(`node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Buffers"),
        )
        .withTarget(
            basicPrometheusQuery(`node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Cached"),
        )
        .withTarget(
            basicPrometheusQuery(`node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Free"),
        );
};

export const memoryUsageGauge = (): GaugePanelBuilder => {
    const query = `100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)`;

    return defaultGauge()
        .title("Memory Usage")
        .span(6)
        .min(30)
        .max(100)
        .unit("percent")
        .thresholds(
            new ThresholdsConfigBuilder()
                .mode(ThresholdsMode.Absolute)
                .steps([
                    {value: null, color: "rgba(50, 172, 45, 0.97)"},
                    {value: 80.0, color: "rgba(237, 129, 40, 0.89)"},
                    {value: 90.0, color: "rgba(245, 54, 54, 0.9)"},
                ])
        )
        .withTarget(basicPrometheusQuery(query, ""));
};
