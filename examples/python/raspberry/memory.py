import grafana_foundation_sdk.models.dashboard as dashboard
from grafana_foundation_sdk.builders.common import StackingConfig as StackingConfigBuilder
from grafana_foundation_sdk.builders.dashboard import ThresholdsConfig as ThresholdsConfigBuilder
from grafana_foundation_sdk.cog import builder as cogbuilder
from grafana_foundation_sdk.models.common import StackingMode

from .common import default_timeseries, basic_prometheus_query, default_gauge


def memory_usage_timeseries() -> cogbuilder.Builder[dashboard.Panel]:
    mem_used_query = '''(
  node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}
-
  node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}
)'''

    return (
        default_timeseries()
        .title("Memory Usage")
        .span(18)
        .stacking(
            StackingConfigBuilder()
            .mode(StackingMode.NORMAL)
        )
        .thresholds(
            ThresholdsConfigBuilder()
            .mode(dashboard.ThresholdsMode.ABSOLUTE)
            .steps([
                dashboard.Threshold(color="green"),
                dashboard.Threshold(color="red", value=80.0),
            ])
        )
        .min(0)
        .unit("bytes")
        .decimals(2)
        .with_target(basic_prometheus_query(mem_used_query, "Used"))
        .with_target(basic_prometheus_query('node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}', "Buffers"))
        .with_target(basic_prometheus_query('node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}', "Cached"))
        .with_target(basic_prometheus_query('node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}', "Free"))
    )


def memory_usage_gauge() -> cogbuilder.Builder[dashboard.Panel]:
    query = '''100 - (
  avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
  avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
* 100)'''

    return (
        default_gauge()
        .title("Memory Usage")
        .span(6)
        .min(30)
        .max(100)
        .unit("percent")
        .thresholds(
            ThresholdsConfigBuilder()
            .mode(dashboard.ThresholdsMode.ABSOLUTE)
            .steps([
                dashboard.Threshold(color="rgba(50, 172, 45, 0.97)"),
                dashboard.Threshold(color="rgba(237, 129, 40, 0.89)", value=80.0),
                dashboard.Threshold(color="rgba(245, 54, 54, 0.9)", value=90.0),
            ])
        )
        .with_target(basic_prometheus_query(query, ""))
    )

