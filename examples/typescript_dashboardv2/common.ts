import {VisualizationBuilder as TimeseriesVizBuilder} from "../../generated/typescript/src/timeseries";
import {VisualizationBuilder as LogsVizBuilder} from "../../generated/typescript/src/logs";
import {VisualizationBuilder as GaugeVizBuilder} from "../../generated/typescript/src/gauge";
import {
    GraphDrawStyle,
    LegendDisplayMode,
    LegendPlacement,
    LogsSortOrder,
    ReduceDataOptionsBuilder,
    VisibilityMode,
    VizLegendOptionsBuilder,
    VizOrientation
} from "../../generated/typescript/src/common";
import * as prometheus from "../../generated/typescript/src/prometheus";
import {PromQueryFormat} from "../../generated/typescript/src/prometheus";
import * as loki from "../../generated/typescript/src/loki";

export const basicPrometheusQuery = (query: string, legend: string): prometheus.QueryBuilder => {
    return new prometheus.QueryBuilder()
        .expr(query)
        .legendFormat(legend);
};

export const basicLokiQuery = (query: string): loki.QueryBuilder => {
    return new loki.QueryBuilder()
        .expr(query);
};

export const tablePrometheusQuery = (query: string): prometheus.QueryBuilder => {
    return new prometheus.QueryBuilder()
        .expr(query)
        .instant()
        .format(PromQueryFormat.Table);
};

export const defaultTimeseries = (): TimeseriesVizBuilder => {
    return new TimeseriesVizBuilder()
        .lineWidth(1)
        .fillOpacity(10)
        .drawStyle(GraphDrawStyle.Line)
        .showPoints(VisibilityMode.Never)
        .legend(
            new VizLegendOptionsBuilder()
                .showLegend(true)
                .placement(LegendPlacement.Bottom)
                .displayMode(LegendDisplayMode.List)
        );
};

export const defaultLogs = (): LogsVizBuilder => {
    return new LogsVizBuilder()
        .showTime(true)
        .enableLogDetails(true)
        .sortOrder(LogsSortOrder.Descending)
        .wrapLogMessage(true);
};

export const defaultGauge = (): GaugeVizBuilder => {
    return new GaugeVizBuilder()
        .orientation(VizOrientation.Auto)
        .reduceOptions(
            new ReduceDataOptionsBuilder()
                .calcs(["lastNotNull"])
                .values(false)
        );
};
