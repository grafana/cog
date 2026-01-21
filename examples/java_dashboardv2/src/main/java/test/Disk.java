package test;

import java.util.List;
import java.util.Map;

import com.grafana.foundation.common.FieldTextAlignment;
import com.grafana.foundation.common.TableCellHeight;
import com.grafana.foundation.common.TableFooterOptionsBuilder;
import com.grafana.foundation.dashboardv2beta1.DynamicConfigValue;
import com.grafana.foundation.dashboardv2beta1.PanelBuilder;
import com.grafana.foundation.dashboardv2beta1.QueryGroupBuilder;
import com.grafana.foundation.dashboardv2beta1.TargetBuilder;
import com.grafana.foundation.dashboardv2beta1.TransformationBuilder;
import com.grafana.foundation.table.TableVizConfigKindBuilder;
import com.grafana.foundation.units.Constants;

public class Disk {
        public static PanelBuilder diskIOTimeseries() {
                return new PanelBuilder()
                                .title("Disk I/O")
                                .visualization(Common.defaultTimeSeries()
                                                .fillOpacity(0.0)
                                                .unit(Constants.BitsPerSecondSI)
                                                .overrideByRegexp("/ io time/", List.of(
                                                                new DynamicConfigValue("unit",
                                                                                Constants.PercentUnit))))
                                .data(new QueryGroupBuilder().targets(List.of(
                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                "rate(node_disk_read_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                "{{device}} read")).refId("A"),
                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                "rate(node_disk_written_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                "{{device}} written")).refId("B"),
                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                "rate(node_disk_io_time_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                "{{device}} IO time")).refId("C"))));
        }

        public static PanelBuilder diskSpaceUsageTable() {
                return new PanelBuilder()
                                .title("Disk Space Usage")
                                .visualization(new TableVizConfigKindBuilder()
                                                .align(FieldTextAlignment.AUTO)
                                                .unit(Constants.BytesSI)
                                                .cellHeight(TableCellHeight.SM)
                                                .footer(new TableFooterOptionsBuilder().countRows(false)
                                                                .reducer(List.of("sum")))
                                                .overrideByName("Mounted on",
                                                                List.of(new DynamicConfigValue("custom.width", 260)))
                                                .overrideByName("Size",
                                                                List.of(new DynamicConfigValue("custom.width", 93)))
                                                .overrideByName("Used",
                                                                List.of(new DynamicConfigValue("custom.width", 72)))
                                                .overrideByName("Available",
                                                                List.of(new DynamicConfigValue("custom.width", 88)))
                                                .overrideByName("Used, %", List.of(
                                                                new DynamicConfigValue("unit", "percentunit"),
                                                                new DynamicConfigValue("custom.cellOptions",
                                                                                Map.of("mode", "gradient", "type",
                                                                                                "gauge")),
                                                                new DynamicConfigValue("min", 0),
                                                                new DynamicConfigValue("max", 1))))
                                .data(new QueryGroupBuilder()
                                                .targets(List.of(
                                                                new TargetBuilder().query(Common.tablePrometheusQuery(
                                                                                "max by (mountpoint) (node_filesystem_size_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\", fstype!=\"\"})"))
                                                                                .refId("A"),
                                                                new TargetBuilder().query(Common.tablePrometheusQuery(
                                                                                "max by (mountpoint) (node_filesystem_avail_bytes{job=\"integrations/raspberrypi-node\", instance=\"$instance\", fstype!=\"\"})"))
                                                                                .refId("B")))
                                                .transformation(transformation("groupBy", transformer1Options()))
                                                .transformation(transformation("merge", null))
                                                .transformation(transformation("calculateField", transformer3Options()))
                                                .transformation(transformation("calculateField", transformer4Options()))
                                                .transformation(transformation("organize", transformer5Options()))
                                                .transformation(transformation("sortBy", transformer6Options())));
        }

        private static TransformationBuilder transformation(String id, Object options) {
                return new TransformationBuilder().id(id).kind(id).options(options);
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
}
