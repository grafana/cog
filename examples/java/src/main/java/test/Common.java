package test;

import java.util.List;

import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.cog.variants.Dataquery;
import com.grafana.foundation.common.GraphDrawStyle;
import com.grafana.foundation.common.LegendDisplayMode;
import com.grafana.foundation.common.LegendPlacement;
import com.grafana.foundation.common.LogsSortOrder;
import com.grafana.foundation.common.ReduceDataOptionsBuilder;
import com.grafana.foundation.common.VisibilityMode;
import com.grafana.foundation.common.VizLegendOptionsBuilder;
import com.grafana.foundation.common.VizOrientation;
import com.grafana.foundation.dashboard.DataSourceRef;
import com.grafana.foundation.prometheus.PromQueryFormat;
import com.grafana.foundation.timeseries.PanelBuilder;

public class Common {

    // ??
    public static Builder<Dataquery> basicPrometheusQuery(String query, String legend) {
        return new com.grafana.foundation.prometheus.DataqueryBuilder().expr(query).legendFormat(legend);
    }

    public static Builder<Dataquery> basicLokiQuery(String query) {
        return new com.grafana.foundation.loki.DataqueryBuilder().expr(query);
    }

    // ??
    public static Builder<Dataquery> tablePrometheusQuery(String query, String ref) {
        return new com.grafana.foundation.prometheus.DataqueryBuilder().expr(query).instant()
                .legendFormat(PromQueryFormat.TABLE.Value()).refId(ref);
    }

    public static PanelBuilder defaultTimeSeries() {
        return new PanelBuilder().lineWidth(1.0).fillOpacity(10.0).drawStyle(GraphDrawStyle.LINE)
                .showPoints(VisibilityMode.NEVER).legend(new VizLegendOptionsBuilder().showLegend(true)
                        .placement(LegendPlacement.BOTTOM).displayMode(LegendDisplayMode.LIST));
    }

    public static com.grafana.foundation.logs.PanelBuilder defaultLogs() {
        return new com.grafana.foundation.logs.PanelBuilder().span(24)
                .datasource(new DataSourceRef("loki", "grafana-cloud-logs")).showTime(true).enableLogDetails(true)
                .sortOrder(LogsSortOrder.DESCENDING).wrapLogMessage(true);
    }

    public static com.grafana.foundation.gauge.PanelBuilder defaultGauge() {
        return new com.grafana.foundation.gauge.PanelBuilder().orientation(VizOrientation.AUTO)
                .reduceOptions(new ReduceDataOptionsBuilder().calcs(List.of("lastNotNull")).values(false));
    }
}
