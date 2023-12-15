from generated.builders.timeseries import Panel as TimeseriesBuilder
from generated.builders.gauge import Panel as GaugeBuilder
from generated.builders.prometheus import Dataquery as PrometheusQueryBuilder
from generated.models.common import GraphDrawStyle, VisibilityMode, LegendPlacement, LegendDisplayMode, VizOrientation
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

