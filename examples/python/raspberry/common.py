from grafana_foundation_sdk.builders.timeseries import Panel as TimeseriesBuilder
from grafana_foundation_sdk.builders.gauge import Panel as GaugeBuilder
from grafana_foundation_sdk.builders.logs import Panel as LogsBuilder
from grafana_foundation_sdk.builders.prometheus import Dataquery as PrometheusQueryBuilder
from grafana_foundation_sdk.builders.loki import Dataquery as LokiQueryBuilder
from grafana_foundation_sdk.models.dashboard import DataSourceRef
from grafana_foundation_sdk.models.prometheus import PromQueryFormat
from grafana_foundation_sdk.models.common import GraphDrawStyle, VisibilityMode, LegendPlacement, LegendDisplayMode, VizOrientation, LogsSortOrder
from grafana_foundation_sdk.builders.common import VizLegendOptions as VizLegendOptionsBuilder, ReduceDataOptions as ReduceDataOptionsBuilder


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


def table_prometheus_query(query: str, ref_id: str) -> PrometheusQueryBuilder:
    return (
        PrometheusQueryBuilder()
        .expr(query)
        .instant()
        .format(PromQueryFormat.TABLE)
        .ref_id(ref_id)
    )


def basic_loki_query(query: str) -> LokiQueryBuilder:
    return LokiQueryBuilder().expr(query)

