from generated.builders.dashboard import (
    Dashboard as DashboardBuilder,
    TimePicker as TimePickerBuilder,
    Row as RowBuilder
)
from generated.models.common import TimeZoneBrowser
from generated.models.dashboard import DashboardCursorSync
from generated.cog.encoder import JSONEncoder
from examples.python.raspberry.cpu import cpu_usage_timeseries, cpu_load_average_timeseries, cpu_temperature_gauge
from examples.python.raspberry.memory import memory_usage_timeseries, memory_usage_gauge


def build_dashboard() -> DashboardBuilder:
    builder = (
        DashboardBuilder("[TEST] Node Exporter / Raspberry")
        .uid("test-dashboard-raspberry")
        .tags(["generated", "raspberrypi-node-integration"])
        .refresh("30s")
        .time("now-30m", "now")
        #.timezone(TimeZoneBrowser)
        .timezone("browser")
        .timepicker(
            TimePickerBuilder()
            .refresh_intervals(["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"])
            .time_options(["5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"])
        )
        .tooltip(DashboardCursorSync.CROSSHAIR)
        # CPU
        .with_row(RowBuilder("CPU"))
        .with_panel(cpu_usage_timeseries())
        .with_panel(cpu_temperature_gauge())
        .with_panel(cpu_load_average_timeseries())
        # Memory
        .with_row(RowBuilder("Memory"))
        .with_panel(memory_usage_timeseries())
        .with_panel(memory_usage_gauge())
        # Disk
        .with_row(RowBuilder("Disk"))
        # Network
        .with_row(RowBuilder("Network"))
        # Logs
        .with_row(RowBuilder("Logs"))
    )

    return builder


if __name__ == '__main__':
    dashboard = build_dashboard().build()
    encoder = JSONEncoder(sort_keys=True, indent=2)

    print(
        encoder.encode(dashboard)
    )
