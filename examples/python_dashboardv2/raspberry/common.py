import grafana_foundation_sdk.builders.gauge as gauge
import grafana_foundation_sdk.builders.logs as logs
import grafana_foundation_sdk.builders.loki as loki
import grafana_foundation_sdk.builders.prometheus as prometheus
import grafana_foundation_sdk.builders.timeseries as timeseries
from grafana_foundation_sdk.builders.common import (
    VizLegendOptions,
    ReduceDataOptions,
    StackingConfig,
)
from grafana_foundation_sdk.models.common import (
    GraphDrawStyle,
    VisibilityMode,
    LegendPlacement,
    LegendDisplayMode,
    VizOrientation,
    LogsSortOrder,
    StackingMode,
)
from grafana_foundation_sdk.models.prometheus import PromQueryFormat


def default_timeseries() -> timeseries.Visualization:
    return (
        timeseries.Visualization()
        .line_width(1)
        .fill_opacity(10)
        .draw_style(GraphDrawStyle.LINE)
        .show_points(VisibilityMode.NEVER)
        .stacking(StackingConfig().mode(StackingMode.NORMAL))
        .legend(
            VizLegendOptions()
            .show_legend(True)
            .placement(LegendPlacement.BOTTOM)
            .display_mode(LegendDisplayMode.LIST)
        )
    )


def default_logs() -> logs.Visualization:
    return (
        logs.Visualization()
        .show_time(True)
        .enable_log_details(True)
        .sort_order(LogsSortOrder.DESCENDING)
        .wrap_log_message(True)
    )


def default_gauge() -> gauge.Visualization:
    return (
        gauge.Visualization()
        .orientation(VizOrientation.AUTO)
        .reduce_options(ReduceDataOptions().calcs(["lastNotNull"]).values(False))
    )


def basic_prometheus_query(query: str, legend: str) -> prometheus.Query:
    return prometheus.Query().expr(query).legend_format(legend)


def table_prometheus_query(query: str) -> prometheus.Query:
    return (
        prometheus.Query()
        .expr(query)
        .instant()
        .format(PromQueryFormat.TABLE)
    )


def basic_loki_query(query: str) -> loki.Query:
    return loki.Query().expr(query)
