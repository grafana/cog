import {PanelBuilder as TimeseriesPanelBuilder} from "../../generated/timeseries";
import {PanelBuilder as GaugePanelBuilder} from "../../generated/gauge";
import {basicPrometheusQuery, defaultGauge, defaultTimeseries} from "./common";
import {StackingConfigBuilder, StackingMode} from "../../generated/common";
import {ThresholdsConfigBuilder, ThresholdsMode} from "../../generated/dashboard";

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
        .unit("bytes")
        .targets([
            basicPrometheusQuery(memUsedQuery, "Used"),
            basicPrometheusQuery(`node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Buffers"),
            basicPrometheusQuery(`node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Cached"),
            basicPrometheusQuery(`node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Free"),
        ]);
};

export const memoryUsageGauge = (): GaugePanelBuilder => {
    const query = `100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)`;

    return defaultGauge()
        .title("Memory Usage")
        .min(30)
        .max(100)
        .unit("percent")
        .thresholds(
            new ThresholdsConfigBuilder()
                .mode(ThresholdsMode.ThresholdsModeAbsolute)
                .steps([
                    {value: null, color: "rgba(50, 172, 45, 0.97)"},
                    {value: 80.0, color: "rgba(237, 129, 40, 0.89)"},
                    {value: 90.0, color: "rgba(245, 54, 54, 0.9)"},
                ])
        )
        .targets([
            basicPrometheusQuery(query, ""),
        ]);
};
