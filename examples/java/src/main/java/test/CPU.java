package test;

import java.util.List;

import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.common.StackingConfigBuilder;
import com.grafana.foundation.common.StackingMode;
import com.grafana.foundation.dashboard.Panel;
import com.grafana.foundation.dashboard.Threshold;
import com.grafana.foundation.dashboard.ThresholdsConfigBuilder;
import com.grafana.foundation.dashboard.ThresholdsMode;

public class CPU {
        public static Builder<Panel> cpuUsageTimeseries() {
                String query = "((1 - sum without (mode) (rate(node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", mode=~\"idle|iowait|steal\", instance=\"$instance\"}[$__rate_interval]))) "
                                +
                                "/ ignoring(cpu) group_left " +
                                "count without (cpu, mode) (node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", mode=\"idle\", instance=\"$instance\"}))";

                return Common.defaultTimeSeries().title("CPU Usage").span(18)
                                .stacking(new StackingConfigBuilder().mode(StackingMode.NORMAL))
                                .thresholds(new ThresholdsConfigBuilder().mode(ThresholdsMode.ABSOLUTE).steps(List.of(
                                                new Threshold(0.0, "green"),
                                                new Threshold(80.0, "red"))))
                                .max(1.0).min(0.0).withTarget(Common.basicPrometheusQuery(query, "{{ cpu }}"));
        }

        public static Builder<Panel> loadAverageTimeseries() {
                return Common.defaultTimeSeries().title("Load Average").span(18).thresholds(
                                new ThresholdsConfigBuilder().mode(ThresholdsMode.ABSOLUTE).steps(List.of(
                                                new Threshold(0.0, "green"),
                                                new Threshold(80.0, "red"))))
                                .min(0.0).unit("short")
                                .withTarget(Common.basicPrometheusQuery(
                                                "node_load1{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                "1m load average"))
                                .withTarget(Common.basicPrometheusQuery(
                                                "node_load5{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                "5m load average"))
                                .withTarget(Common.basicPrometheusQuery(
                                                "node_load15{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                "15m load average"))
                                .withTarget(Common.basicPrometheusQuery(
                                                "count(node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", mode=\"idle\"})",
                                                "logical cores"));
        }

        public static Builder<Panel> cpuTemperatureGauge() {
                return Common.defaultGauge().title("CPU Temperature").span(6).min(30.0).max(100.0).unit("celsius")
                                .thresholds(new ThresholdsConfigBuilder().mode(ThresholdsMode.ABSOLUTE).steps(List.of(
                                                new Threshold(0.0, "rgba(50, 172, 45, 0.97)"),
                                                new Threshold(65.0, "rgba(237, 129, 40, 0.89)"),
                                                new Threshold(85.0, "rgba(245, 54, 54, 0.9)"))))
                                .withTarget(Common.basicPrometheusQuery(
                                                "avg(node_hwmon_temp_celsius{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}",
                                                ""));
        }
}
