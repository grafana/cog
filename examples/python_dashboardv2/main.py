import json

from grafana_foundation_sdk.builders.dashboardv2beta1 import (
    Dashboard,
    TimeSettings,
    DatasourceVariable,
    QueryVariable,
    TabsLayout,
    TabsLayoutTab,
    AutoGridLayout,
    AutoGridLayoutItem,
)
from grafana_foundation_sdk.cog.encoder import JSONEncoder
from grafana_foundation_sdk.cog.plugins import register_default_plugins
from grafana_foundation_sdk.models.common import (
    TimeZoneBrowser,
)
from grafana_foundation_sdk.models.dashboardv2beta1 import (
    Dashboard as DashboardModel,
    DashboardCursorSync,
    VariableHide,
    VariableOption,
    VariableRefresh,
    VariableSort,
)

from raspberry.common import basic_prometheus_query
from raspberry.cpu import (
    cpu_usage_timeseries,
    cpu_load_average_timeseries,
    cpu_temperature_gauge,
)
from raspberry.disk import disk_io_timeseries, disk_space_usage_table
from raspberry.logs import (
    errors_in_system_logs,
    all_system_logs,
    auth_logs,
    kernel_logs,
)
from raspberry.memory import memory_usage_timeseries, memory_usage_gauge
from raspberry.network import (
    network_received_timeseries,
    network_transmitted_timeseries,
)


def build_dashboard() -> Dashboard:
    builder = (
        Dashboard("[TEST] Node Exporter / Raspberry")
        .tags(["generated", "raspberrypi-node-integration"])
        .cursor_sync(DashboardCursorSync.CROSSHAIR)
        .time_settings(
            TimeSettings()
            .auto_refresh("30s")
            .auto_refresh_intervals(
                ["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"]
            )
            .from_val("now-30m")
            .to("now")
            .timezone(TimeZoneBrowser),
            )
        # "Data Source" variable
        .variable(
            DatasourceVariable("datasource")
            .label("Data Source")
            .hide(VariableHide.DONT_HIDE)
            .plugin_id("prometheus")
            .current(
                VariableOption(
                    selected=True,
                    text="grafanacloud-potatopi-prom",
                    value="grafanacloud-prom",
                )
            )
        )
        # "Instance" variable
        .variable(
            QueryVariable("instance")
            .label("Instance")
            .hide(VariableHide.DONT_HIDE)
            .refresh(VariableRefresh.ON_TIME_RANGE_CHANGED)
            .query(
                basic_prometheus_query('label_values(node_uname_info{job="integrations/raspberrypi-node", sysname!="Darwin"}, instance)', "")
            )
            .current(VariableOption(selected=False, text="potato", value="potato"))
            .sort(VariableSort.DISABLED)
        )
        .elements(
            {
                # CPU
                "cpu_usage": cpu_usage_timeseries(),
                "cpu_temp": cpu_temperature_gauge(),
                "load_avg": cpu_load_average_timeseries(),
                # Memory
                "mem_usage": memory_usage_timeseries(),
                "mem_usage_current": memory_usage_gauge(),
                # Disk
                "disk_io": disk_io_timeseries(),
                "disk_usage": disk_space_usage_table(),
                # Network
                "network_in": network_transmitted_timeseries(),
                "network_out": network_received_timeseries(),
                # Logs
                "sys_error_logs": errors_in_system_logs(),
                "auth_logs": auth_logs(),
                "kernel_logs": kernel_logs(),
                "all_sys_logs": all_system_logs(),
            }
        )
        .tabs_layout(TabsLayout()
                     .tab(TabsLayoutTab("CPU")
                          .auto_grid_layout(AutoGridLayout()
                                          .item(AutoGridLayoutItem("cpu_usage"))
                                          .item(AutoGridLayoutItem("cpu_temp"))
                                          .item(AutoGridLayoutItem("load_avg"))))
                     .tab(TabsLayoutTab("Memory")
                          .auto_grid_layout(AutoGridLayout()
                                            .item(AutoGridLayoutItem("mem_usage"))
                                            .item(AutoGridLayoutItem("mem_usage_current"))))
                     .tab(TabsLayoutTab("Disk")
                          .auto_grid_layout(AutoGridLayout()
                                            .item(AutoGridLayoutItem("disk_io"))
                                            .item(AutoGridLayoutItem("disk_usage"))))
                     .tab(TabsLayoutTab("Network")
                          .auto_grid_layout(AutoGridLayout()
                                            .item(AutoGridLayoutItem("network_in"))
                                            .item(AutoGridLayoutItem("network_out"))))
                     .tab(TabsLayoutTab("Logs")
                          .auto_grid_layout(AutoGridLayout()
                                            .item(AutoGridLayoutItem("sys_error_logs"))
                                            .item(AutoGridLayoutItem("auth_logs"))
                                            .item(AutoGridLayoutItem("kernel_logs"))
                                            .item(AutoGridLayoutItem("all_sys_logs"))))
        )
    )

    return builder


if __name__ == "__main__":
    dashboard = build_dashboard().build()
    encoder = JSONEncoder(sort_keys=True, indent=2)

    dashboard_json = encoder.encode(dashboard)

    print(dashboard_json)

    register_default_plugins()

    decoded_dashboard = DashboardModel.from_json(json.loads(dashboard_json))

    print(decoded_dashboard)
