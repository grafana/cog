package test;

import java.util.List;

import com.grafana.foundation.common.StackingConfigBuilder;
import com.grafana.foundation.common.StackingMode;
import com.grafana.foundation.dashboardv2beta1.PanelBuilder;
import com.grafana.foundation.dashboardv2beta1.QueryGroupBuilder;
import com.grafana.foundation.dashboardv2beta1.TargetBuilder;
import com.grafana.foundation.dashboardv2beta1.Threshold;
import com.grafana.foundation.dashboardv2beta1.ThresholdsConfigBuilder;
import com.grafana.foundation.dashboardv2beta1.ThresholdsMode;
import com.grafana.foundation.units.Constants;

public class Memory {

        public static PanelBuilder memoryUsageTimeseries() {
                String memUsedQuery = "(" +
                                "  node_memory_MemTotal_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"
                                +
                                "-" +
                                "  node_memory_MemFree_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"
                                +
                                "-" +
                                "  node_memory_Buffers_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"
                                +
                                "-" +
                                "  node_memory_Cached_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"
                                +
                                ")";

                return new PanelBuilder()
                                .title("Memory Usage")
                                .visualization(Common.defaultTimeSeries()
                                                .stacking(new StackingConfigBuilder()
                                                                .mode(StackingMode.NORMAL))
                                                .thresholds(new ThresholdsConfigBuilder()
                                                                .mode(ThresholdsMode.ABSOLUTE)
                                                                .steps(List.of(
                                                                                new Threshold(0.0, "green"),
                                                                                new Threshold(80.0, "red"))))
                                                .min(0.0)
                                                .unit(Constants.BytesIEC)
                                                .decimals(2.0))
                                .data(new QueryGroupBuilder()
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(memUsedQuery,
                                                                                "Used"))
                                                                .refId("A"))
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "node_memory_Buffers_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                                                "Buffers"))
                                                                .refId("B"))
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "node_memory_Cached_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                                                "Cached"))
                                                                .refId("C"))
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "node_memory_MemFree_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                                                "Free"))
                                                                .refId("D")));
        }

        public static PanelBuilder memoryUsageGauge() {
                String query = "100 - (" +
                                "  avg(node_memory_MemAvailable_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}) /"
                                +
                                "  avg(node_memory_MemTotal_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\"})"
                                +
                                "* 100)";

                return new PanelBuilder()
                                .title("Memory Usage")
                                .visualization(Common.defaultGauge()
                                                .min(30.0)
                                                .max(100.0)
                                                .unit(Constants.Percent)
                                                .thresholds(new ThresholdsConfigBuilder()
                                                                .mode(ThresholdsMode.ABSOLUTE).steps(List.of(
                                                                                new Threshold(0.0,
                                                                                                "rgba(50, 172, 45, 0.97)"),
                                                                                new Threshold(80.0,
                                                                                                "rgba(237, 129, 40, 0.89)"),
                                                                                new Threshold(90.0,
                                                                                                "rgba(245, 54, 54, 0.9)")))))
                                .data(new QueryGroupBuilder().target(new TargetBuilder()
                                                .query(Common.basicPrometheusQuery(query, "")).refId("A")));
        }
}
