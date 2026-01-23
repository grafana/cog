import grafana_foundation_sdk.models.dashboardv2beta1 as dashboard
from grafana_foundation_sdk.builders.common import TableFooterOptions
from grafana_foundation_sdk.builders.dashboardv2beta1 import (
    Panel,
    QueryGroup,
    Target,
    Transformation,
)
from grafana_foundation_sdk.builders.table import Visualization as TableVisualization
from grafana_foundation_sdk.models.common import (
    FieldTextAlignment,
    TableAutoCellOptions,
    TableCellHeight,
)

from .common import default_timeseries, basic_prometheus_query, table_prometheus_query


def disk_io_timeseries() -> Panel:
    return (
        Panel()
        .title("Disk I/O")
        .visualization(
            default_timeseries()
            .fill_opacity(0)
            .unit("Bps")
            .override_by_regexp(
                "/ io time/",
                [
                    dashboard.DynamicConfigValue(id_val="unit", value="percentunit"),
                ],
            )
        )
        .data(
            QueryGroup()
            .target(
                Target().query(
                    basic_prometheus_query(
                        'rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])',
                        "{{device}} read",
                    )
                ).ref_id("A")
            )
            .target(
                Target().query(
                    basic_prometheus_query(
                        'rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])',
                        "{{device}} written",
                    )
                ).ref_id("B")
            )
            .target(
                Target().query(
                    basic_prometheus_query(
                        'rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])',
                        "{{device}} IO time",
                    )
                ).ref_id("C")
            )
        )
    )


def disk_space_usage_table() -> Panel:
    return (
        Panel()
        .title("Disk Space Usage")
        .visualization(
            TableVisualization()
            .align(FieldTextAlignment.AUTO)
            .cell_options(TableAutoCellOptions())
            .unit("decbytes")
            .cell_height(TableCellHeight.SM)
            .footer(TableFooterOptions().count_rows(False).reducer(["sum"]))
            # Overrides configuration
            .override_by_name(
                "Mounted on",
                [
                    dashboard.DynamicConfigValue(id_val="custom.width", value=260),
                ],
            )
            .override_by_name(
                "Size",
                [
                    dashboard.DynamicConfigValue(id_val="custom.width", value=93),
                ],
            )
            .override_by_name(
                "Used",
                [
                    dashboard.DynamicConfigValue(id_val="custom.width", value=72),
                ],
            )
            .override_by_name(
                "Available",
                [
                    dashboard.DynamicConfigValue(id_val="custom.width", value=88),
                ],
            )
            .override_by_name(
                "Used, %",
                [
                    dashboard.DynamicConfigValue(id_val="unit", value="percentunit"),
                    dashboard.DynamicConfigValue(
                        id_val="custom.cellOptions",
                        value={"mode": "gradient", "type": "gauge"},
                    ),
                    dashboard.DynamicConfigValue(id_val="max", value="1"),
                    dashboard.DynamicConfigValue(id_val="min", value="0"),
                ],
            )
        )
        .data(
            QueryGroup()
            .target(
                Target().query(
                    table_prometheus_query(
                        'max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})'
                    )
                ).ref_id("A")
            )
            .target(
                Target().query(
                    table_prometheus_query(
                        'max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})',
                    )
                ).ref_id("B")
            )
            # Transformations
            .transformation(
                Transformation()
                .id("groupBy")
                .kind("groupBy")
                .options(
                    {
                        "fields": {
                            "Value #A": {
                                "aggregations": ["lastNotNull"],
                                "operation": "aggregate",
                            },
                            "Value #B": {
                                "aggregations": ["lastNotNull"],
                                "operation": "aggregate",
                            },
                            "mountpoint": {"aggregations": [], "operation": "groupby"},
                        }
                    }
                )
            )
            .transformation(Transformation().id("merge").kind("merge").options({}))
            .transformation(
                Transformation()
                .id("calculateField")
                .kind("calculateField")
                .options(
                    {
                        "alias": "Used",
                        "binary": {
                            "left": "Value #A (lastNotNull)",
                            "operator": "-",
                            "reducer": "sum",
                            "right": "Value #B (lastNotNull)",
                        },
                        "mode": "binary",
                        "reduce": {"reducer": "sum"},
                    }
                )
            )
            .transformation(
                Transformation()
                .id("calculateField")
                .kind("calculateField")
                .options(
                    {
                        "alias": "Used, %",
                        "binary": {
                            "left": "Used",
                            "operator": "/",
                            "reducer": "sum",
                            "right": "Value #A (lastNotNull)",
                        },
                        "mode": "binary",
                        "reduce": {"reducer": "sum"},
                    }
                )
            )
            .transformation(
                Transformation()
                .id("organize")
                .kind("organize")
                .options(
                    {
                        "excludeByName": {},
                        "indexByName": {},
                        "renameByName": {
                            "Value #A (lastNotNull)": "Size",
                            "Value #B (lastNotNull)": "Available",
                            "mountpoint": "Mounted on",
                        },
                    }
                )
            )
            .transformation(
                Transformation()
                .id("sortBy")
                .kind("sortBy")
                .options({"fields": {}, "sort": [{"field": "Mounted on"}]})
            )
        )
    )
