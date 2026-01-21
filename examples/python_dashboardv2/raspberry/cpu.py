import grafana_foundation_sdk.models.dashboardv2beta1 as dashboard
from grafana_foundation_sdk.builders.dashboardv2beta1 import (
    Panel,
    ThresholdsConfig,
    QueryGroup,
    Target,
)

from .common import default_timeseries, basic_prometheus_query, default_gauge


def cpu_usage_timeseries() -> Panel:
    query = """(
        (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
    / ignoring(cpu) group_left
    count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
    )"""

    return (
        Panel()
        .title("CPU Usage")
        .visualization(
            default_timeseries()
            .thresholds(
                ThresholdsConfig()
                .mode(dashboard.ThresholdsMode.ABSOLUTE)
                .steps(
                    [
                        dashboard.Threshold(color="green"),
                        dashboard.Threshold(color="red", value=80.0),
                    ]
                )
            )
            .min(0)
            .max(1)
            .unit("percentunit")
        )
        .data(
            QueryGroup().target(
                Target().query(basic_prometheus_query(query, "{{ cpu }}")).ref_id("A")
            )
        )
    )


def cpu_temperature_gauge() -> Panel:
    return (
        Panel()
        .title("CPU Temperature")
        .visualization(
            default_gauge()
            .min(0)
            .max(100)
            .unit("celsius")
            .thresholds(
                ThresholdsConfig()
                .mode(dashboard.ThresholdsMode.ABSOLUTE)
                .steps(
                    [
                        dashboard.Threshold(color="rgba(50, 172, 45, 0.97)"),
                        dashboard.Threshold(
                            color="rgba(237, 129, 40, 0.89)", value=65.0
                        ),
                        dashboard.Threshold(color="rgba(245, 54, 54, 0.9)", value=85.0),
                    ]
                )
            )
        )
        .data(
            QueryGroup().target(
                Target().query(
                    basic_prometheus_query(
                        'avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})',
                        "",
                    )
                ).ref_id("A")
            )
        )
    )


def cpu_load_average_timeseries() -> Panel:
    return (
        Panel()
        .title("Load Average")
        .visualization(
            default_timeseries()
            .thresholds(
                ThresholdsConfig()
                .mode(dashboard.ThresholdsMode.ABSOLUTE)
                .steps(
                    [
                        dashboard.Threshold(color="green"),
                        dashboard.Threshold(color="red", value=80.0),
                    ]
                )
            )
            .min(0)
            .unit("short")
        )
        .data(
            QueryGroup()
            .target(
                Target().query(
                    basic_prometheus_query(
                        'node_load1{job="integrations/raspberrypi-node", instance="$instance"}',
                        "1m load average",
                    )
                ).ref_id("A")
            )
            .target(
                Target().query(
                    basic_prometheus_query(
                        'node_load5{job="integrations/raspberrypi-node", instance="$instance"}',
                        "5m load average",
                    )
                ).ref_id("B")
            )
            .target(
                Target().query(
                    basic_prometheus_query(
                        'node_load15{job="integrations/raspberrypi-node", instance="$instance"}',
                        "15m load average",
                    )
                ).ref_id("C")
            )
            .target(
                Target().query(
                    basic_prometheus_query(
                        'count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})',
                        "logical cores",
                    )
                ).ref_id("D")
            )
        )
    )
