import grafana_foundation_sdk.models.dashboardv2beta1 as dashboard
from grafana_foundation_sdk.builders.dashboardv2beta1 import (
    Panel,
    QueryGroup,
    Target,
    ThresholdsConfig,
)

from .common import default_timeseries, basic_prometheus_query, default_gauge


def memory_usage_timeseries() -> Panel:
    mem_used_query = """(
  node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}
)"""

    return (
        Panel()
        .title("Memory Usage")
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
            .unit("bytes")
            .decimals(2)
        )
        .data(
            QueryGroup()
            .target(Target().query(basic_prometheus_query(mem_used_query, "Used")))
            .target(
                Target().query(
                    basic_prometheus_query(
                        'node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}',
                        "Buffers",
                    )
                )
            )
            .target(
                Target().query(
                    basic_prometheus_query(
                        'node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}',
                        "Cached",
                    )
                )
            )
            .target(
                Target().query(
                    basic_prometheus_query(
                        'node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}',
                        "Free",
                    )
                )
            )
        )
    )


def memory_usage_gauge() -> Panel:
    query = """100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)"""

    return (
        Panel()
        .title("Memory Usage")
        .visualization(
            default_gauge()
            .min(30)
            .max(100)
            .unit("percent")
            .thresholds(
                ThresholdsConfig()
                .mode(dashboard.ThresholdsMode.ABSOLUTE)
                .steps(
                    [
                        dashboard.Threshold(color="rgba(50, 172, 45, 0.97)"),
                        dashboard.Threshold(
                            color="rgba(237, 129, 40, 0.89)", value=80.0
                        ),
                        dashboard.Threshold(color="rgba(245, 54, 54, 0.9)", value=90.0),
                    ]
                )
            )
        )
        .data(QueryGroup().target(Target().query(basic_prometheus_query(query, "")).ref_id("A")))
    )
