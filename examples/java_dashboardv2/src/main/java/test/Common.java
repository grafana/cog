package test;

import java.util.List;

import com.grafana.foundation.common.GraphDrawStyle;
import com.grafana.foundation.common.LegendDisplayMode;
import com.grafana.foundation.common.LegendPlacement;
import com.grafana.foundation.common.LogsSortOrder;
import com.grafana.foundation.common.ReduceDataOptionsBuilder;
import com.grafana.foundation.common.VisibilityMode;
import com.grafana.foundation.common.VizLegendOptionsBuilder;
import com.grafana.foundation.common.VizOrientation;
import com.grafana.foundation.gauge.GaugeVizConfigKindBuilder;
import com.grafana.foundation.logs.LogsVizConfigKindBuilder;
import com.grafana.foundation.loki.LokiDataQueryKindBuilder;
import com.grafana.foundation.prometheus.PromQueryFormat;
import com.grafana.foundation.prometheus.PrometheusDataQueryKindBuilder;
import com.grafana.foundation.timeseries.TimeseriesVizConfigKindBuilder;

public class Common {

    // ??
    public static PrometheusDataQueryKindBuilder basicPrometheusQuery(String query, String legend) {
        return new PrometheusDataQueryKindBuilder().expr(query).legendFormat(legend);
    }

    public static LokiDataQueryKindBuilder basicLokiQuery(String query) {
        return new LokiDataQueryKindBuilder().expr(query);
    }

    // ??
    public static PrometheusDataQueryKindBuilder tablePrometheusQuery(String query) {
        return new PrometheusDataQueryKindBuilder()
                .expr(query)
                .instant()
                .legendFormat(PromQueryFormat.TABLE.Value());
    }

    public static TimeseriesVizConfigKindBuilder defaultTimeSeries() {
        return new TimeseriesVizConfigKindBuilder()
                .lineWidth(1.0)
                .fillOpacity(10.0)
                .drawStyle(GraphDrawStyle.LINE)
                .showPoints(VisibilityMode.NEVER)
                .legend(new VizLegendOptionsBuilder()
                        .showLegend(true)
                        .placement(LegendPlacement.BOTTOM)
                        .displayMode(LegendDisplayMode.LIST));
    }

    public static LogsVizConfigKindBuilder defaultLogs() {
        return new LogsVizConfigKindBuilder()
                .showTime(true)
                .enableLogDetails(true)
                .sortOrder(LogsSortOrder.DESCENDING)
                .wrapLogMessage(true);
    }

    public static GaugeVizConfigKindBuilder defaultGauge() {
        return new GaugeVizConfigKindBuilder()
                .orientation(VizOrientation.AUTO)
                .reduceOptions(new ReduceDataOptionsBuilder()
                        .calcs(List.of("lastNotNull"))
                        .values(false));
    }
}
