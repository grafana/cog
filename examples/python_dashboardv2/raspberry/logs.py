from grafana_foundation_sdk.builders.dashboardv2alpha0 import Panel, QueryGroup, Target
from .common import default_logs, basic_loki_query


def errors_in_system_logs() -> Panel:
    return (
        Panel()
        .title("Errors in system logs")
        .visualization(default_logs())
        .data(
            QueryGroup()
                .target(Target().query(basic_loki_query('{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}')))
                .target(Target().query(basic_loki_query('{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"')))
        )
    )


def auth_logs() -> Panel:
    return (
        Panel()
        .title("Auth logs")
        .visualization(default_logs())
        .data(
            QueryGroup()
                .target(Target().query(basic_loki_query('{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}')))
                .target(Target().query(basic_loki_query('{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}')))
        )
    )


def kernel_logs() -> Panel:
    return (
        Panel()
        .title("Kernel logs")
        .visualization(default_logs())
        .data(
            QueryGroup()
                .target(Target().query(basic_loki_query('{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}')))
                .target(Target().query(basic_loki_query('{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}')))
        )
    )


def all_system_logs() -> Panel:
    return (
        Panel()
        .title("All system logs")
        .visualization(default_logs())
        .data(
            QueryGroup()
                .target(Target().query(basic_loki_query('{transport!="", job="integrations/raspberrypi-node", instance="$instance"}')))
                .target(Target().query(basic_loki_query('{filename!="", job="integrations/raspberrypi-node", instance="$instance"}')))
        )
    )
