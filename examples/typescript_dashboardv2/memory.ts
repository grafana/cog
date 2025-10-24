import {basicPrometheusQuery, defaultGauge, defaultTimeseries} from "./common";
import {StackingConfigBuilder, StackingMode} from "../../generated/typescript/src/common";
import {
    PanelBuilder,
    QueryGroupBuilder,
    TargetBuilder,
    ThresholdsConfigBuilder,
    ThresholdsMode
} from "../../generated/typescript/src/dashboardv2beta1";

export const memoryUsageTimeseries = (): PanelBuilder => {
    const memUsedQuery = `(
  node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}
)`;

    return new PanelBuilder()
        .title("Memory Usage")
        .visualization(defaultTimeseries()
            .stacking(new StackingConfigBuilder().mode(StackingMode.Normal))
            .min(0)
            .unit("bytes")
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
            new TargetBuilder().query(basicPrometheusQuery(memUsedQuery, "Used")),
            new TargetBuilder().query(basicPrometheusQuery(`node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Buffers").refId("A")),
            new TargetBuilder().query(basicPrometheusQuery(`node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Cached").refId("B")),
            new TargetBuilder().query(basicPrometheusQuery(`node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}`, "Free").refId("C")),
        ]));
};

export const memoryUsageGauge = (): PanelBuilder => {
    const query = `100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)`;

    return new PanelBuilder()
        .title("Memory Usage")
        .visualization(defaultGauge()
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
        )
        .data(
            new QueryGroupBuilder()
                .target(new TargetBuilder().query(basicPrometheusQuery(query, "")).refId("A"))
        );
};
