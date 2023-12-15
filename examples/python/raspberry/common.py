from generated.builders.timeseries import Panel as TimeseriesBuilder
from generated.builders.gauge import Panel as GaugeBuilder
from generated.builders.logs import Panel as LogsBuilder
from generated.builders.prometheus import Dataquery as PrometheusQueryBuilder
from generated.builders.loki import Dataquery as LokiQueryBuilder
from generated.models.dashboard import DataSourceRef
from generated.models.common import GraphDrawStyle, VisibilityMode, LegendPlacement, LegendDisplayMode, VizOrientation, LogsSortOrder
from generated.builders.common import VizLegendOptions as VizLegendOptionsBuilder, ReduceDataOptions as ReduceDataOptionsBuilder


def default_timeseries() -> TimeseriesBuilder:
    return (
        TimeseriesBuilder()
        .line_width(1)
        .fill_opacity(10)
        .draw_style(GraphDrawStyle.LINE)
        .show_points(VisibilityMode.NEVER)
        .legend(
            VizLegendOptionsBuilder()
            .show_legend(True)
            .placement(LegendPlacement.BOTTOM)
            .display_mode(LegendDisplayMode.LIST)
        )
    )


def default_logs() -> LogsBuilder:
    return (
        LogsBuilder()
        .span(24)
        .datasource(DataSourceRef(type_val="loki", uid="grafanacloud-logs"))
        .show_time(True)
        .enable_log_details(True)
        .sort_order(LogsSortOrder.DESCENDING)
        .wrap_log_message(True)
    )


def default_gauge() -> GaugeBuilder:
    return (
        GaugeBuilder()
        .orientation(VizOrientation.AUTO)
        .reduce_options(
            ReduceDataOptionsBuilder()
            .calcs(["lastNotNull"])
            .values(False)
        )
    )


def basic_prometheus_query(query: str, legend: str) -> PrometheusQueryBuilder:
    return (
        PrometheusQueryBuilder()
        .expr(query)
        .legend_format(legend)
    )


def basic_loki_query(query: str) -> LokiQueryBuilder:
    return LokiQueryBuilder().expr(query)

