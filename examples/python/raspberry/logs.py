import grafana_foundation_sdk.models.dashboard as dashboard
from grafana_foundation_sdk.cog import builder as cogbuilder
from .common import default_logs, basic_loki_query


def errors_in_system_logs() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_logs()
        .title("Errors in system logs")
        .with_target(basic_loki_query('{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}'))
        .with_target(basic_loki_query('{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"'))
    )


def auth_logs() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_logs()
        .title("Auth logs")
        .with_target(basic_loki_query('{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}'))
        .with_target(basic_loki_query('{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}'))
    )


def kernel_logs() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_logs()
        .title("Kernel logs")
        .with_target(basic_loki_query('{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}'))
        .with_target(basic_loki_query('{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}'))
    )


def all_system_logs() -> cogbuilder.Builder[dashboard.Panel]:
    return (
        default_logs()
        .title("All system logs")
        .with_target(basic_loki_query('{transport!="", job="integrations/raspberrypi-node", instance="$instance"}'))
        .with_target(basic_loki_query('{filename!="", job="integrations/raspberrypi-node", instance="$instance"}'))
    )
