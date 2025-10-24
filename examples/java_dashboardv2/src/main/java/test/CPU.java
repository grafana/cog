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

public class CPU {
        public static PanelBuilder cpuUsageTimeseries() {
                String query = "((1 - sum without (mode) (rate(node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", mode=~\"idle|iowait|steal\", instance=\"$instance\"}[$__rate_interval]))) "
                                +
                                "/ ignoring(cpu) group_left " +
                                "count without (cpu, mode) (node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", mode=\"idle\", instance=\"$instance\"}))";

                return new PanelBuilder()
                                .title("CPU Usage")
                                .visualization(Common.defaultTimeSeries()
                                                .stacking(new StackingConfigBuilder()
                                                                .mode(StackingMode.NORMAL))
                                                .min(0.0)
                                                .max(1.0)
                                                .unit(Constants.PercentUnit)
                                                .thresholds(new ThresholdsConfigBuilder()
                                                                .mode(ThresholdsMode.ABSOLUTE)
                                                                .steps(List.of(new Threshold(0.0, "green"),
                                                                                new Threshold(80.0, "red")))))
                                .data(new QueryGroupBuilder().target(new TargetBuilder()
                                                .query(Common.basicPrometheusQuery(query, "{{ cpu }}")).refId("A")));
        }

        public static PanelBuilder loadAverageTimeseries() {
                return new PanelBuilder()
                                .title("Load Average")
                                .visualization(Common.defaultTimeSeries()
                                                .stacking(new StackingConfigBuilder()
                                                                .mode(StackingMode.NORMAL))
                                                .min(0.0)
                                                .unit(Constants.Short)
                                                .thresholds(new ThresholdsConfigBuilder()
                                                                .mode(ThresholdsMode.ABSOLUTE)
                                                                .steps(List.of(new Threshold(0.0, "green"),
                                                                                new Threshold(80.0, "red")))))
                                .data(new QueryGroupBuilder()
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "node_load1{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                                                "1m load average"))
                                                                .refId("A"))
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "node_load5{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                                                "5m load average"))
                                                                .refId("B"))
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "node_load15{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                                                "15m load average"))
                                                                .refId("C")));
        }

        public static PanelBuilder cpuTemperatureGauge() {
                return new PanelBuilder()
                                .title("CPU Temperature")
                                .visualization(Common.defaultGauge()
                                                .min(30.0)
                                                .max(100.0)
                                                .unit(Constants.Celsius)
                                                .thresholds(new ThresholdsConfigBuilder()
                                                                .mode(ThresholdsMode.ABSOLUTE)
                                                                .steps(List.of(new Threshold(0.0,
                                                                                "rgba(50, 172, 45, 0.97)"),
                                                                                new Threshold(65.0,
                                                                                                "rgba(237, 129, 40, 0.89)"),
                                                                                new Threshold(85.0,
                                                                                                "rgba(245, 54, 54, 0.9)")))))
                                .data(new QueryGroupBuilder()
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "avg(node_hwmon_temp_celsius{job=\"integrations/raspberrypi-node\", instance=\"$instance\"})",
                                                                                ""))
                                                                .refId("A")));
        }
}
