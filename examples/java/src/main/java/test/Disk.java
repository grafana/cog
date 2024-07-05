package test;

import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.common.FieldTextAlignment;
import com.grafana.foundation.common.TableCellHeight;
import com.grafana.foundation.common.TableFooterOptions;
import com.grafana.foundation.dashboard.DashboardFieldConfigSourceOverrides;
import com.grafana.foundation.dashboard.DataTransformerConfig;
import com.grafana.foundation.dashboard.DynamicConfigValue;
import com.grafana.foundation.dashboard.MatcherConfig;
import com.grafana.foundation.timeseries.PanelBuilder;

import java.util.List;
import java.util.Map;

public class Disk {
    public static PanelBuilder diskIOTimeseries() {
        MatcherConfig matcher = new MatcherConfig();
        matcher.id = "byRegexp";
        matcher.options = "/ io time/";

        DynamicConfigValue dcv = new DynamicConfigValue();
        dcv.id = "unit";
        dcv.value = "percentunit";

        return Common.defaultTimeSeries().
                Title("Disk I/O").
                FillOpacity(0.0).
                Unit("Bps").
                WithTarget(Common.basicPrometheusQuery("rate(node_disk_read_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])", "{{device}} read")).
                WithTarget(Common.basicPrometheusQuery("rate(node_disk_written_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])", "{{device}} written")).
                WithTarget(Common.basicPrometheusQuery("rate(node_disk_io_time_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])", "{{device}} IO time")).
                WithOverride(new DashboardFieldConfigSourceOverrides.Builder().
                        Matcher(matcher).
                        Properties(List.of(dcv))
                );
    }

    public static com.grafana.foundation.table.PanelBuilder diskSpaceUsageTable() {
        return new com.grafana.foundation.table.PanelBuilder().
                Title("Disk Space Usage").
                Align(FieldTextAlignment.AUTO).
                Unit("decbytes").
                CellHeight(TableCellHeight.SM).
                Footer(new TableFooterOptions.Builder().
                        CountRows(false).
                        Reducer(List.of("sum"))
                ).
                WithTarget(Common.tablePrometheusQuery("max by (mountpoint) (node_filesystem_size_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\", fstype!=\"\"})", "A")).
                WithTarget(Common.tablePrometheusQuery("max by (mountpoint) (node_filesystem_avail_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\", fstype!=\"\"})", "B")).
                WithTransformation(transformer1()).
                WithTransformation(transformer2()).
                WithTransformation(transformer3()).
                WithTransformation(transformer4()).
                WithTransformation(transformer5()).
                WithTransformation(transformer6()).
                WithOverride(defaultOverrides("Mounted on", 260)).
                WithOverride(defaultOverrides("Size", 93)).
                WithOverride(defaultOverrides("Used", 72)).
                WithOverride(defaultOverrides("Available", 88)).
                WithOverride(complexOverrides());
    }

    private static DataTransformerConfig transformer1() {
        DataTransformerConfig dataTransformerConfig = new DataTransformerConfig();
        dataTransformerConfig.id = "groupBy";
        dataTransformerConfig.options = Map.of(
                "fields", Map.of(
                        "Value #A", Map.of(
                                "aggregations", List.of("lastNotNull"),
                                "operation", "aggregate"
                        ),
                        "Value #B", Map.of(
                                "aggregations", List.of("lastNotNull"),
                                "operation", "aggregate"
                        ),
                        "mountpoint", Map.of(
                                "aggregations", List.of(),
                                "operation", "groupby"
                        )
                )
        );

        return dataTransformerConfig;
    }

    private static DataTransformerConfig transformer2() {
        DataTransformerConfig dataTransformerConfig = new DataTransformerConfig();
        dataTransformerConfig.id = "merge";

        return dataTransformerConfig;
    }

    private static DataTransformerConfig transformer3() {
        DataTransformerConfig dataTransformerConfig = new DataTransformerConfig();
        dataTransformerConfig.id = "calculateField";
        dataTransformerConfig.options = Map.of(
                "alias", "Used",
                "binary", Map.of(
                        "left", "Value #A (lastNotNull)",
                        "operator", "-",
                        "reducer", "sum",
                        "right", "Value #B (lastNotNull)"
                ), "mode", "binary",
                "reduce", Map.of("reducer", "sum"));

        return dataTransformerConfig;
    }

    private static DataTransformerConfig transformer4() {
        DataTransformerConfig dataTransformerConfig = new DataTransformerConfig();
        dataTransformerConfig.id = "calculateField";
        dataTransformerConfig.options = Map.of(
                "alias", "Used, %",
                "binary", Map.of(
                        "left", "Used",
                        "operator", "/",
                        "reducer", "sum",
                        "right", "Value #A (lastNotNull)"
                ),
                "mode", "binary",
                "reduce", Map.of("reducer", "sum"));

        return dataTransformerConfig;
    }

    private static DataTransformerConfig transformer5() {
        DataTransformerConfig dataTransformerConfig = new DataTransformerConfig();
        dataTransformerConfig.id = "organize";
        dataTransformerConfig.options = Map.of(
                "excludeByName", List.of(),
                "indexByName", List.of(),
                "renameByName", Map.of(
                        "Value #A (lastNotNull)", "Size",
                        "Value #B (lastNotNull)", "Available",
                        "mountpoint", "Mounted on"
                )
        );

        return dataTransformerConfig;
    }

    private static DataTransformerConfig transformer6() {
        DataTransformerConfig dataTransformerConfig = new DataTransformerConfig();
        dataTransformerConfig.id = "sortBy";
        dataTransformerConfig.options = Map.of(
                "fields", List.of(),
                "sort", Map.of("field", "Mounted on")
        );

        return dataTransformerConfig;
    }

    private static Builder<DashboardFieldConfigSourceOverrides> defaultOverrides(String options, Integer value) {
        MatcherConfig matcherConfig = new MatcherConfig();
        matcherConfig.id = "byName";
        matcherConfig.options = options;

        DynamicConfigValue dynamicConfigValue = new DynamicConfigValue();
        dynamicConfigValue.id = "custom.width";
        dynamicConfigValue.value = value;
        return new DashboardFieldConfigSourceOverrides.Builder().Matcher(matcherConfig).Properties(List.of(dynamicConfigValue));
    }

    private static Builder<DashboardFieldConfigSourceOverrides> complexOverrides() {
        MatcherConfig matcherConfig = new MatcherConfig();
        matcherConfig.id = "byName";
        matcherConfig.options = "Used, %";

        DynamicConfigValue dynamicConfigValue1 = new DynamicConfigValue();
        dynamicConfigValue1.id = "unit";
        dynamicConfigValue1.value = "percentunit";

        DynamicConfigValue dynamicConfigValue2 = new DynamicConfigValue();
        dynamicConfigValue2.id = "custom.cellOptions";
        dynamicConfigValue2.value = Map.of("mode", "gradient", "type", "gauge");

        DynamicConfigValue dynamicConfigValue3 = new DynamicConfigValue();
        dynamicConfigValue3.id = "min";
        dynamicConfigValue3.value = 0;

        DynamicConfigValue dynamicConfigValue4 = new DynamicConfigValue();
        dynamicConfigValue4.id = "max";
        dynamicConfigValue4.value = 1;

        return new DashboardFieldConfigSourceOverrides.Builder().
                Matcher(matcherConfig).
                Properties(List.of(dynamicConfigValue1, dynamicConfigValue2, dynamicConfigValue3, dynamicConfigValue4));
    }
}
