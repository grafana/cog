from grafana_foundation_sdk.builders.dashboardv2beta1 import Panel, QueryGroup, Target

from .common import default_timeseries, basic_prometheus_query


def network_received_timeseries() -> Panel:
    return (
        Panel()
        .title("Network Received")
        .description("Network received (bits/s)")
        .visualization(default_timeseries().min(0).unit("bps").fill_opacity(0))
        .data(
            QueryGroup().target(
                Target().query(
                    basic_prometheus_query(
                        'rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8',
                        "{{ device }}",
                    )
                ).ref_id("A")
            )
        )
    )


def network_transmitted_timeseries() -> Panel:
    return (
        Panel()
        .title("Network Transmitted")
        .description("Network transmitted (bits/s)")
        .visualization(default_timeseries().min(0).unit("bps").fill_opacity(0))
        .data(
            QueryGroup().target(
                Target().query(
                    basic_prometheus_query(
                        'rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8',
                        "{{ device }}",
                    )
                ).ref_id("A")
            )
        )
    )
