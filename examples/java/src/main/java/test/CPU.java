package test;

import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.common.StackingConfig;
import com.grafana.foundation.common.StackingMode;
import com.grafana.foundation.dashboard.Panel;
import com.grafana.foundation.dashboard.Threshold;
import com.grafana.foundation.dashboard.ThresholdsConfig;
import com.grafana.foundation.dashboard.ThresholdsMode;

import java.util.List;

public class CPU {
    public static Builder<Panel> cpuUsageTimeseries() {
        String query = "((1 - sum without (mode) (rate(node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", mode=~\"idle|iowait|steal\", instance=\"$instance\"}[$__rate_interval]))) " +
                "/ ignoring(cpu) group_left " +
                "count without (cpu, mode) (node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", mode=\"idle\", instance=\"$instance\"}))";

        Threshold th1 = new Threshold();
        th1.color = "green";

        Threshold th2 = new Threshold();
        th2.color = "red";
        th2.value = 80.0;

        return Common.defaultTimeSeries().
                Title("CPU Usage").
                Span(18).
                Stacking(new StackingConfig.Builder().Mode(StackingMode.NORMAL)).
                Thresholds(new ThresholdsConfig.Builder().
                        Mode(ThresholdsMode.ABSOLUTE).
                        Steps(List.of(th1, th2))).
                Max(1.0).
                Min(0.0).
                WithTarget(Common.basicPrometheusQuery(query, "{{ cpu }}"));
    }

    public static Builder<Panel> loadAverageTimeseries() {
        Threshold th1 = new Threshold();
        th1.color = "green";

        Threshold th2 = new Threshold();
        th2.value = 80.0;
        th2.color = "red";
        return Common.defaultTimeSeries().Title("Load Average").
                Span(18).
                Thresholds(
                        new ThresholdsConfig.Builder().
                                Mode(ThresholdsMode.ABSOLUTE).
                                Steps(List.of(th1, th2))
                ).
                Min(0.0).
                Unit("short").
                WithTarget(Common.basicPrometheusQuery("node_load1{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}", "1m load average")).
                WithTarget(Common.basicPrometheusQuery("node_load5{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}", "5m load average")).
                WithTarget(Common.basicPrometheusQuery("node_load15{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}", "15m load average")).
                WithTarget(Common.basicPrometheusQuery("count(node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", mode=\"idle\"})", "logical cores"));
    }

    public static Builder<Panel> cpuTemperatureGauge() {
        Threshold th1 = new Threshold();
        th1.color = "rgba(50, 172, 45, 0.97)";

        Threshold th2 = new Threshold();
        th2.value = 65.0;
        th2.color = "rgba(237, 129, 40, 0.89)";

        Threshold th3 = new Threshold();
        th3.value = 85.0;
        th3.color = "rgba(245, 54, 54, 0.9)";
        return Common.defaultGauge().
                Title("CPU Temperature").
                Span(6).
                Min(30.0).
                Max(100.0).
                Unit("celsius").
                Thresholds(new ThresholdsConfig.Builder().
                        Mode(ThresholdsMode.ABSOLUTE).
                        Steps(List.of(th1, th2, th3))).
                WithTarget(Common.basicPrometheusQuery("avg(node_hwmon_temp_celsius{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}", ""));
    }
}
