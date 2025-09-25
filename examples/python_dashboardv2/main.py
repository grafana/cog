import json
from grafana_foundation_sdk.builders.dashboardv2alpha1 import (
    Dashboard,
    TimeSettings,
    DatasourceVariable,
    QueryVariable,
    GridLayout,
    GridLayoutItem,
    GridLayoutRow,
    DataSourceRef,
)
from grafana_foundation_sdk.models.dashboardv2alpha1 import (
    DashboardSpec as DashboardModel,
    DashboardCursorSync,
    VariableHide,
    VariableOption,
    VariableRefresh,
    VariableSort,
)
from grafana_foundation_sdk.models.common import (
    TimeZoneBrowser,
)
from grafana_foundation_sdk.cog.encoder import JSONEncoder
from grafana_foundation_sdk.cog.plugins import register_default_plugins
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
from raspberry.common import basic_prometheus_query


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
            .datasource(DataSourceRef(type_val="prometheus", uid="$datasource"))
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
        .grid_layout(
            GridLayout()
            .row(GridLayoutRow("CPU"))
            .item(GridLayoutItem("cpu_usage"))
            .row(GridLayoutRow("Memory"))
            .row(GridLayoutRow("Disk"))
            .row(GridLayoutRow("Network"))
            .row(GridLayoutRow("Logs"))
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
