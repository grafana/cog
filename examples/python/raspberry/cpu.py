import grafana_foundation_sdk.models.dashboard as dashboard
from grafana_foundation_sdk.builders.dashboard import ThresholdsConfig as ThresholdsConfigBuilder
from grafana_foundation_sdk.cog import builder as cogbuilder
from .common import default_timeseries, basic_prometheus_query, default_gauge
from grafana_foundation_sdk.builders.common import StackingConfig as StackingConfigBuilder
from grafana_foundation_sdk.models.common import StackingMode


def cpu_usage_timeseries() -> cogbuilder.Builder[dashboard.Panel]:
    query = '''(
        (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
    / ignoring(cpu) group_left
    count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
    )'''

    return (
        default_timeseries()
        .title("CPU Usage")
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
        .max(1)
        .unit("percentunit")
        .with_target(basic_prometheus_query(query, "{{ cpu }}"))
    )


def cpu_temperature_gauge() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_gauge()
        .title("CPU Temperatore")
        .span(6)
        .min(0)
        .max(100)
        .unit("celsius")
        .thresholds(
            ThresholdsConfigBuilder()
            .mode(dashboard.ThresholdsMode.ABSOLUTE)
            .steps([
                dashboard.Threshold(color="rgba(50, 172, 45, 0.97)"),
                dashboard.Threshold(color="rgba(237, 129, 40, 0.89)", value=65.0),
                dashboard.Threshold(color="rgba(245, 54, 54, 0.9)", value=85.0),
            ])
        )
        .with_target(basic_prometheus_query("avg(node_hwmon_temp_celsius{job=\"integrations/raspberrypi-node\", instance=\"$instance\"})", ""))
    )


def cpu_load_average_timeseries() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_timeseries()
        .title("Load Average")
        .span(18)
        .thresholds(
            ThresholdsConfigBuilder()
            .mode(dashboard.ThresholdsMode.ABSOLUTE)
            .steps([
                dashboard.Threshold(color="green"),
                dashboard.Threshold(color="red", value=80.0),
            ])
        )
        .min(0)
        .unit("short")
        .with_target(basic_prometheus_query("node_load1{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}", "1m load average"))
        .with_target(basic_prometheus_query("node_load5{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}", "5m load average"))
        .with_target(basic_prometheus_query("node_load15{job=\"integrations/raspberrypi-node\", instance=\"$instance\"}", "15m load average"))
        .with_target(basic_prometheus_query("count(node_cpu_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", mode=\"idle\"})", "logical cores"))
    )
