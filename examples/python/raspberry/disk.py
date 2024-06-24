import generated.python.grafana_foundation_sdk.models.dashboard as dashboard
from generated.python.grafana_foundation_sdk.cog import builder as cogbuilder
from generated.python.grafana_foundation_sdk.builders.table import Panel as TableBuilder
from generated.python.grafana_foundation_sdk.models.common import FieldTextAlignment, TableAutoCellOptions, TableCellHeight
from generated.python.grafana_foundation_sdk.builders.common import TableFooterOptions
from .common import default_timeseries, basic_prometheus_query, table_prometheus_query


def disk_io_timeseries() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_timeseries()
        .title("Disk I/O")
        .fill_opacity(0)
        .unit("Bps")
        .with_target(
            basic_prometheus_query('rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', "{{device}} read")
        )
        .with_target(
            basic_prometheus_query('rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', "{{device}} written")
        )
        .with_target(
            basic_prometheus_query('rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', "{{device}} IO time")
        )
        .with_override(
            dashboard.MatcherConfig(id_val="byRegexp", options="/ io time/"),
            [
                dashboard.DynamicConfigValue(id_val="unit", value="percentunit"),
            ],
        )
    )


def disk_space_usage_table() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        TableBuilder()
        .title("Disk Space Usage")
        .align(FieldTextAlignment.AUTO)
        .cell_options(TableAutoCellOptions())
        .unit("decbytes")
        .cell_height(TableCellHeight.SM)
        .footer(
            TableFooterOptions()
            .count_rows(False)
            .reducer(["sum"])
        )
        .with_target(
            table_prometheus_query('max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})', "A")
        )
        .with_target(
            table_prometheus_query('max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})', "B")
        )
        # Transformations
        .with_transformation(
            dashboard.DataTransformerConfig(
                id_val="groupBy",
                options={
                    "fields": {
                        "Value #A": {
                            "aggregations": ["lastNotNull"],
                            "operation": "aggregate"
                        },
                        "Value #B": {
                            "aggregations": ["lastNotNull"],
                            "operation": "aggregate"
                        },
                        "mountpoint": {
                            "aggregations": [],
                            "operation": "groupby"
                        }
                    }
                }
            )
        )
        .with_transformation(
            dashboard.DataTransformerConfig(id_val="merge", options={})
        )
        .with_transformation(
            dashboard.DataTransformerConfig(
                id_val="calculateField",
                options={
                    "alias": "Used",
                    "binary": {
                        "left": "Value #A (lastNotNull)",
                        "operator": "-",
                        "reducer": "sum",
                        "right": "Value #B (lastNotNull)"
                    },
                    "mode": "binary",
                    "reduce": {"reducer": "sum"}
                }
            )
        )
        .with_transformation(
            dashboard.DataTransformerConfig(
                id_val="calculateField",
                options={
                    "alias": "Used, %",
                    "binary": {
                        "left": "Used",
                        "operator": "/",
                        "reducer": "sum",
                        "right": "Value #A (lastNotNull)"
                    },
                    "mode": "binary",
                    "reduce": {"reducer": "sum"}
                }
            )
        )
        .with_transformation(
            dashboard.DataTransformerConfig(
                id_val="organize",
                options={
                    "excludeByName": {},
                    "indexByName": {},
                    "renameByName": {
                        "Value #A (lastNotNull)": "Size",
                        "Value #B (lastNotNull)": "Available",
                        "mountpoint": "Mounted on"
                    }
                }
            )
        )
        .with_transformation(
            dashboard.DataTransformerConfig(
                id_val="sortBy",
                options={
                    "fields": {},
                    "sort": [{"field": "Mounted on"}]
                }
            )
        )
        # Overrides configuration
        .with_override(
            dashboard.MatcherConfig(id_val="byName", options="Mounted on"),
            [dashboard.DynamicConfigValue(id_val="custom.width", value=260)]
        )
        .with_override(
            dashboard.MatcherConfig(id_val="byName", options="Size"),
            [dashboard.DynamicConfigValue(id_val="custom.width", value=93)]
        )
        .with_override(
            dashboard.MatcherConfig(id_val="byName", options="Used"),
            [dashboard.DynamicConfigValue(id_val="custom.width", value=72)]
        )
        .with_override(
            dashboard.MatcherConfig(id_val="byName", options="Available"),
            [dashboard.DynamicConfigValue(id_val="custom.width", value=88)]
        )
        .with_override(
            dashboard.MatcherConfig(id_val="byName", options="Used, %"),
            [
                dashboard.DynamicConfigValue(id_val="unit", value="percentunit"),
                dashboard.DynamicConfigValue(
                    id_val="custom.cellOptions",
                    value={"mode": "gradient", "type": "gauge"},
                ),
                dashboard.DynamicConfigValue(id_val="max", value="1"),
                dashboard.DynamicConfigValue(id_val="min", value="0"),
            ]
        )
    )

