import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.cog.variants.Dataquery;
import com.grafana.foundation.common.*;
import com.grafana.foundation.dashboard.DataSourceRef;
import com.grafana.foundation.prometheus.PromQueryFormat;
import com.grafana.foundation.timeseries.PanelBuilder;

import java.util.List;

public class Common {

    // ??
    public static Builder<Dataquery> basicPrometheusQuery(String query, String legend) {
        return new com.grafana.foundation.prometheus.Dataquery.Builder().
                Expr(query).
                LegendFormat(legend);
    }

    public static Builder<Dataquery> basicLokiQuery(String query) {
        return new com.grafana.foundation.prometheus.Dataquery.Builder().Expr(query);
    }

    // ??
    public static Builder<Dataquery> tablePrometheusQuery(String query, String ref) {
        return new com.grafana.foundation.prometheus.Dataquery.Builder().
                Expr(query).
                Instant().
                LegendFormat(PromQueryFormat.TABLE.Value()).
                RefId(ref);
    }

    public static PanelBuilder defaultTimeSeries() {
        return new PanelBuilder().
                LineWidth(1.0).
                FillOpacity(10.0).
                DrawStyle(GraphDrawStyle.LINE).
                ShowPoints(VisibilityMode.NEVER).
                Legend(new VizLegendOptions.Builder().
                        ShowLegend(true).
                        Placement(LegendPlacement.BOTTOM).
                        DisplayMode(LegendDisplayMode.LIST)
                );
    }

    public static com.grafana.foundation.logs.PanelBuilder defaultLogs() {
        DataSourceRef ref = new DataSourceRef();
        ref.type = "loki";
        ref.uid = "grafana-cloud-logs";
        return new com.grafana.foundation.logs.PanelBuilder().
                Span(24).
                Datasource(ref).
                ShowTime(true).
                EnableLogDetails(true).
                SortOrder(LogsSortOrder.DESCENDING).
                WrapLogMessage(true);
    }

    public static com.grafana.foundation.gauge.PanelBuilder defaultGauge() {
        return new com.grafana.foundation.gauge.PanelBuilder().
                Orientation(VizOrientation.AUTO).
                ReduceOptions(new ReduceDataOptions.Builder().
                        Calcs(List.of("lastNotNull")).
                        Values(false)
                );
    }
}
