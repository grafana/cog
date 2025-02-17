import grafana_foundation_sdk.models.dashboard as dashboard
from grafana_foundation_sdk.cog import builder as cogbuilder
from .common import default_timeseries, basic_prometheus_query


def network_received_timeseries() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_timeseries()
        .title("Network Received")
        .description("Network received (bits/s)")
        .min(0)
        .unit("bps")
        .fill_opacity(0)
        .with_target(
            basic_prometheus_query('rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8', "{{ device }}")
        )
    )


def network_transmitted_timeseries() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_timeseries()
        .title("Network Transmitted")
        .description("Network transmitted (bits/s)")
        .min(0)
        .unit("bps")
        .fill_opacity(0)
        .with_target(
            basic_prometheus_query('rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8', "{{ device }}")
        )
    )

