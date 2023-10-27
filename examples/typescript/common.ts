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

export const basicPrometheusQuery = (query: string, legend: string): prometheus.dataquery => {
    return new prometheus.DataqueryBuilder()
        .expr(query)
        .legendFormat(legend)
        .build();
};

export const basicLokiQuery = (query: string): loki.dataquery => {
    return new loki.DataqueryBuilder()
        .expr(query)
        .build();
};

export const tablePrometheusQuery = (query: string, ref: string): prometheus.dataquery => {
    return new prometheus.DataqueryBuilder()
        .expr(query)
        .instant(true)
        .format(PromQueryFormat.PromQueryFormatTable)
        .refId(ref)
        .build();
};

export const defaultTimeseries = (): TimeseriesPanelBuilder => {
    return new TimeseriesPanelBuilder()
        .lineWidth(1)
        .fillOpacity(10)
        .drawStyle(GraphDrawStyle.GraphDrawStyleLine)
        .showPoints(VisibilityMode.VisibilityModeNever)
        .legend(
            new VizLegendOptionsBuilder()
                .showLegend(true)
                .placement(LegendPlacement.LegendPlacementBottom)
                .displayMode(LegendDisplayMode.LegendDisplayModeList)
        );
};

export const defaultLogs = (): LogsPanelBuilder => {
    return new LogsPanelBuilder()
        .datasource({
            type: "loki",
            uid:  "grafanacloud-logs",
        })
        .showTime(true)
        .enableLogDetails(true)
        .sortOrder(LogsSortOrder.LogsSortOrderDescending)
        .wrapLogMessage(true);
};

export const defaultGauge = (): GaugePanelBuilder => {
    return new GaugePanelBuilder()
        .orientation(VizOrientation.VizOrientationAuto)
        .reduceOptions(
            new ReduceDataOptionsBuilder()
                .calcs(["lastNotNull"])
                .values(false)
        );
};
