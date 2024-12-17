package test;

import java.util.List;
import java.util.Map;

import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.common.FieldTextAlignment;
import com.grafana.foundation.common.TableCellHeight;
import com.grafana.foundation.common.TableFooterOptionsBuilder;
import com.grafana.foundation.dashboard.DashboardFieldConfigSourceOverrides;
import com.grafana.foundation.dashboard.DashboardFieldConfigSourceOverridesBuilder;
import com.grafana.foundation.dashboard.DataTransformerConfig;
import com.grafana.foundation.dashboard.DynamicConfigValue;
import com.grafana.foundation.dashboard.MatcherConfig;
import com.grafana.foundation.timeseries.PanelBuilder;

public class Disk {
        public static PanelBuilder diskIOTimeseries() {
                return Common.defaultTimeSeries().title("Disk I/O").fillOpacity(0.0).unit("Bps")
                                .withTarget(Common.basicPrometheusQuery(
                                                "rate(node_disk_read_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                "{{device}} read"))
                                .withTarget(Common.basicPrometheusQuery(
                                                "rate(node_disk_written_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                "{{device}} written"))
                                .withTarget(Common.basicPrometheusQuery(
                                                "rate(node_disk_io_time_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                "{{device}} IO time"))
                                .withOverride(new DashboardFieldConfigSourceOverridesBuilder()
                                                .matcher(new MatcherConfig("byRegexp", "/ io time/"))
                                                .properties(List.of(new DynamicConfigValue("unit", "percentunit"))));
        }

        public static com.grafana.foundation.table.PanelBuilder diskSpaceUsageTable() {
                return new com.grafana.foundation.table.PanelBuilder().title("Disk Space Usage")
                                .align(FieldTextAlignment.AUTO).unit("decbytes").cellHeight(TableCellHeight.SM)
                                .footer(new TableFooterOptionsBuilder().countRows(false).reducer(List.of("sum")))
                                .withTarget(Common.tablePrometheusQuery(
                                                "max by (mountpoint) (node_filesystem_size_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\", fstype!=\"\"})",
                                                "A"))
                                .withTarget(Common.tablePrometheusQuery(
                                                "max by (mountpoint) (node_filesystem_avail_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\", fstype!=\"\"})",
                                                "B"))
                                .withTransformation(transformation("groupBy", transformer1Options()))
                                .withTransformation(transformation("merge", null))
                                .withTransformation(transformation("calculateField", transformer3Options()))
                                .withTransformation(transformation("calculateField", transformer4Options()))
                                .withTransformation(transformation("organize", transformer5Options()))
                                .withTransformation(transformation("sortBy", transformer6Options()))
                                .withOverride(defaultOverrides("Mounted on", 260))
                                .withOverride(defaultOverrides("Size", 93)).withOverride(defaultOverrides("Used", 72))
                                .withOverride(defaultOverrides("Available", 88))
                                .withOverride(new DashboardFieldConfigSourceOverridesBuilder()
                                                .matcher(new MatcherConfig("byName", "Used, %")).properties(List.of(
                                                                new DynamicConfigValue("unit", "percentunit"),
                                                                new DynamicConfigValue("custom.cellOptions",
                                                                                Map.of("mode", "gradient", "type",
                                                                                                "gauge")),
                                                                new DynamicConfigValue("min", 0),
                                                                new DynamicConfigValue("max", 1))));
        }

        private static DataTransformerConfig transformation(String id, Object options) {
                return new DataTransformerConfig(id, false, null, null, options);
        }

        private static Map<String, Object> transformer1Options() {
                return Map.of(
                                "fields", Map.of(
                                                "Value #A", Map.of(
                                                                "aggregations", List.of("lastNotNull"),
                                                                "operation", "aggregate"),
                                                "Value #B", Map.of(
                                                                "aggregations", List.of("lastNotNull"),
                                                                "operation", "aggregate"),
                                                "mountpoint", Map.of(
                                                                "aggregations", List.of(),
                                                                "operation", "groupby")));
        }

        private static Map<String, Object> transformer3Options() {
                return Map.of(
                                "alias", "Used",
                                "binary", Map.of(
                                                "left", "Value #A (lastNotNull)",
                                                "operator", "-",
                                                "reducer", "sum",
                                                "right", "Value #B (lastNotNull)"),
                                "mode", "binary",
                                "reduce", Map.of("reducer", "sum"));
        }

        private static Map<String, Object> transformer4Options() {
                return Map.of(
                                "alias", "Used, %",
                                "binary", Map.of(
                                                "left", "Used",
                                                "operator", "/",
                                                "reducer", "sum",
                                                "right", "Value #A (lastNotNull)"),
                                "mode", "binary",
                                "reduce", Map.of("reducer", "sum"));
        }

        private static Map<String, Object> transformer5Options() {
                return Map.of(
                                "excludeByName", List.of(),
                                "indexByName", List.of(),
                                "renameByName", Map.of(
                                                "Value #A (lastNotNull)", "Size",
                                                "Value #B (lastNotNull)", "Available",
                                                "mountpoint", "Mounted on"));
        }

        private static Map<String, Object> transformer6Options() {
                return Map.of(
                                "fields", List.of(),
                                "sort", Map.of("field", "Mounted on"));
        }

        private static Builder<DashboardFieldConfigSourceOverrides> defaultOverrides(String options, Integer value) {
                return new DashboardFieldConfigSourceOverridesBuilder().matcher(new MatcherConfig("byName", options))
                                .properties(List.of(new DynamicConfigValue("custom.width", value)));
        }
}
