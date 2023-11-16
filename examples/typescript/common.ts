import {PanelBuilder as TimeseriesPanelBuilder} from "../../generated/timeseries";
import {PanelBuilder as LogsPanelBuilder} from "../../generated/logs";
import {PanelBuilder as GaugePanelBuilder} from "../../generated/gauge";
import {
    GraphDrawStyle,
    LegendDisplayMode,
    LegendPlacement,
    LogsSortOrder,
    ReduceDataOptionsBuilder,
    VisibilityMode,
    VizLegendOptionsBuilder,
    VizOrientation
} from "../../generated/common";
import * as prometheus from "../../generated/prometheus";
import * as loki from "../../generated/loki";
import {PromQueryFormat} from "../../generated/prometheus";

export const basicPrometheusQuery = (query: string, legend: string): prometheus.DataqueryBuilder => {
    return new prometheus.DataqueryBuilder()
        .expr(query)
        .legendFormat(legend);
};

export const basicLokiQuery = (query: string): loki.DataqueryBuilder => {
    return new loki.DataqueryBuilder()
        .expr(query);
};

export const tablePrometheusQuery = (query: string, ref: string): prometheus.DataqueryBuilder => {
    return new prometheus.DataqueryBuilder()
        .expr(query)
        .instant(true)
        .format(PromQueryFormat.Table)
        .refId(ref);
};

export const defaultTimeseries = (): TimeseriesPanelBuilder => {
    return new TimeseriesPanelBuilder()
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

export const defaultLogs = (): LogsPanelBuilder => {
    return new LogsPanelBuilder()
        .span(24)
        .datasource({
            type: "loki",
            uid:  "grafanacloud-logs",
        })
        .showTime(true)
        .enableLogDetails(true)
        .sortOrder(LogsSortOrder.Descending)
        .wrapLogMessage(true);
};

export const defaultGauge = (): GaugePanelBuilder => {
    return new GaugePanelBuilder()
        .orientation(VizOrientation.Auto)
        .reduceOptions(
            new ReduceDataOptionsBuilder()
                .calcs(["lastNotNull"])
                .values(false)
        );
};
