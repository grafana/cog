import json
from grafana_foundation_sdk.builders.dashboard import (
    Dashboard,
    TimePicker,
    Row,
    DatasourceVariable,
    QueryVariable,
)
from grafana_foundation_sdk.models.dashboard import (
    Dashboard as DashboardModel,
    DashboardCursorSync,
    VariableHide,
    VariableOption,
    VariableRefresh,
    DataSourceRef,
    VariableSort,
)
from grafana_foundation_sdk.models.common import (
    TimeZoneBrowser,
)
from grafana_foundation_sdk.cog.encoder import JSONEncoder
from grafana_foundation_sdk.cog.plugins import register_default_plugins
from raspberry.cpu import cpu_usage_timeseries, cpu_load_average_timeseries, cpu_temperature_gauge
from raspberry.disk import disk_io_timeseries, disk_space_usage_table
from raspberry.logs import errors_in_system_logs, all_system_logs, auth_logs, kernel_logs
from raspberry.memory import memory_usage_timeseries, memory_usage_gauge
from raspberry.network import network_received_timeseries, network_transmitted_timeseries


def build_dashboard() -> Dashboard:
    builder = (
        Dashboard("[TEST] Node Exporter / Raspberry")
        .uid("test-dashboard-raspberry")
        .tags(["generated", "raspberrypi-node-integration"])
        .refresh("30s")
        .time("now-30m", "now")
        .timezone(TimeZoneBrowser)
        .timezone("browser")
        .timepicker(
            TimePicker().refresh_intervals(["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"])
        )
        .tooltip(DashboardCursorSync.CROSSHAIR)
        # "Data Source" variable
        .with_variable(
            DatasourceVariable("datasource")
            .label("Data Source")
            .hide(VariableHide.DONT_HIDE)
            .type("prometheus")
            .current(VariableOption(selected=True, text="grafanacloud-potatopi-prom", value="grafanacloud-prom"))
        )
        # "Instance" variable
        .with_variable(
            QueryVariable("instance")
            .label("Instance")
            .hide(VariableHide.DONT_HIDE)
            .refresh(VariableRefresh.ON_TIME_RANGE_CHANGED)
            .query('label_values(node_uname_info{job="integrations/raspberrypi-node", sysname!="Darwin"}, instance)')
            .datasource(DataSourceRef(type_val="prometheus", uid="$datasource"))
            .current(VariableOption(selected=False, text="potato", value="potato"))
            .sort(VariableSort.DISABLED)
        )
        # CPU
        .with_row(Row("CPU"))
        .with_panel(cpu_usage_timeseries())
        .with_panel(cpu_temperature_gauge())
        .with_panel(cpu_load_average_timeseries())
        # Memory
        .with_row(Row("Memory"))
        .with_panel(memory_usage_timeseries())
        .with_panel(memory_usage_gauge())
        # Disk
        .with_row(Row("Disk"))
        .with_panel(disk_io_timeseries())
        .with_panel(disk_space_usage_table())
        # Network
        .with_row(Row("Network"))
        .with_panel(network_transmitted_timeseries())
        .with_panel(network_received_timeseries())
        # Logs
        .with_row(Row("Logs"))
        .with_panel(errors_in_system_logs())
        .with_panel(auth_logs())
        .with_panel(kernel_logs())
        .with_panel(all_system_logs())
    )

    return builder


if __name__ == '__main__':
    dashboard = build_dashboard().build()
    encoder = JSONEncoder(sort_keys=True, indent=2)

    dashboard_json = encoder.encode(dashboard)

    print(dashboard_json)

    register_default_plugins()

    decoded_dashboard = DashboardModel.from_json(
        json.loads(dashboard_json)
    )

    print(decoded_dashboard)
