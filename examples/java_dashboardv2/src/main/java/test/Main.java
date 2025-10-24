package test;

import java.util.List;

import com.grafana.foundation.common.Constants;
import com.grafana.foundation.dashboardv2beta1.AutoGridLayoutBuilder;
import com.grafana.foundation.dashboardv2beta1.AutoGridLayoutItemBuilder;
import com.grafana.foundation.dashboardv2beta1.Dashboard;
import com.grafana.foundation.dashboardv2beta1.DashboardBuilder;
import com.grafana.foundation.dashboardv2beta1.DashboardCursorSync;
import com.grafana.foundation.dashboardv2beta1.DatasourceVariableBuilder;
import com.grafana.foundation.dashboardv2beta1.QueryVariableBuilder;
import com.grafana.foundation.dashboardv2beta1.StringOrArrayOfString;
import com.grafana.foundation.dashboardv2beta1.TabsLayoutBuilder;
import com.grafana.foundation.dashboardv2beta1.TabsLayoutTabBuilder;
import com.grafana.foundation.dashboardv2beta1.TimeSettingsBuilder;
import com.grafana.foundation.dashboardv2beta1.VariableHide;
import com.grafana.foundation.dashboardv2beta1.VariableOption;
import com.grafana.foundation.dashboardv2beta1.VariableRefresh;
import com.grafana.foundation.dashboardv2beta1.VariableSort;
import com.grafana.foundation.prometheus.PrometheusDataQueryKindBuilder;
import com.grafana.relocated.jackson.core.JsonProcessingException;

public class Main {
        public static void main(String[] args) {
                Dashboard dashboard = new DashboardBuilder("[TEST] Node Exporter / Raspberry")
                                .tags(List.of("generated", "raspberrypi-node-integration"))
                                .timeSettings(new TimeSettingsBuilder()
                                                .autoRefresh("30s")
                                                .autoRefreshIntervals(List.of("5s", "10s", "30s", "1m", "5m", "15m",
                                                                "30m", "1h", "2h", "1d"))
                                                .from("now-30m")
                                                .to("now")
                                                .timezone(Constants.TimeZoneBrowser))
                                .cursorSync(DashboardCursorSync.CROSSHAIR)
                                .datasourceVariable(new DatasourceVariableBuilder("datasource")
                                                .label("Data Source")
                                                .hide(VariableHide.DONT_HIDE).pluginId("prometheus")
                                                .current(new VariableOption(true,
                                                                StringOrArrayOfString.createString(
                                                                                "grafanacloud-potatopi-prom"),
                                                                StringOrArrayOfString
                                                                                .createString("grafanacloud-prom"))))
                                .queryVariable(new QueryVariableBuilder("instance")
                                                .label("Instance")
                                                .hide(VariableHide.DONT_HIDE)
                                                .refresh(VariableRefresh.ON_TIME_RANGE_CHANGED)
                                                .query(new PrometheusDataQueryKindBuilder()
                                                                .expr("label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)"))
                                                .current(new VariableOption(true,
                                                                StringOrArrayOfString.createString(
                                                                                "potato"),
                                                                StringOrArrayOfString
                                                                                .createString("potato")))
                                                .sort(VariableSort.DISABLED))
                                // CPU
                                .panel("cpu_usage", CPU.cpuUsageTimeseries())
                                .panel("cpu_temp", CPU.cpuTemperatureGauge())
                                .panel("load_avg", CPU.loadAverageTimeseries())
                                // Memory
                                .panel("mem_usage", Memory.memoryUsageTimeseries())
                                .panel("mem_usage_current", Memory.memoryUsageGauge())
                                // Disk
                                .panel("disk_io", Disk.diskIOTimeseries())
                                .panel("disk_usage", Disk.diskSpaceUsageTable())
                                // Network
                                .panel("network_in", Network.networkReceivedTimeseries())
                                .panel("network_out", Network.networkReceivedTimeseries())
                                // Logs
                                .panel("sys_error_logs", Logs.errorsInSystemLogs())
                                .panel("auth_logs", Logs.authLogs())
                                .panel("kernel_logs", Logs.kernelLogs())
                                .panel("all_sys_logs", Logs.allSystemLogs())
                                // Tabs Layout
                                .tabsLayout(new TabsLayoutBuilder()
                                                .tab(new TabsLayoutTabBuilder("CPU")
                                                                .autoGridLayout(new AutoGridLayoutBuilder()
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "cpu_usage"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "cpu_temp"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "load_avg"))))
                                                .tab(new TabsLayoutTabBuilder("Memory")
                                                                .autoGridLayout(new AutoGridLayoutBuilder()
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "mem_usage"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "mem_usage_current"))))
                                                .tab(new TabsLayoutTabBuilder("Disk")
                                                                .autoGridLayout(new AutoGridLayoutBuilder()
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "disk_io"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "disk_usage"))))
                                                .tab(new TabsLayoutTabBuilder("Network")
                                                                .autoGridLayout(new AutoGridLayoutBuilder()
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "network_in"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "network_out"))))
                                                .tab(new TabsLayoutTabBuilder("Logs")
                                                                .autoGridLayout(new AutoGridLayoutBuilder()
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "sys_error_logs"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "auth_logs"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "kernel_logs"))
                                                                                .item(new AutoGridLayoutItemBuilder(
                                                                                                "all_sys_logs")))))
                                .build();
                try {
                        System.out.println(dashboard.toJSON());
                } catch (JsonProcessingException e) {
                        throw new RuntimeException(e);
                }
        }

}
