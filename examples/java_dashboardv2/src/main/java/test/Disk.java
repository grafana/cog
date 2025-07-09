package test;

import java.util.List;
import java.util.Map;

import com.grafana.foundation.common.FieldTextAlignment;
import com.grafana.foundation.common.TableCellHeight;
import com.grafana.foundation.common.TableFooterOptionsBuilder;
import com.grafana.foundation.dashboardv2alpha1.DynamicConfigValue;
import com.grafana.foundation.dashboardv2alpha1.PanelBuilder;
import com.grafana.foundation.dashboardv2alpha1.QueryGroupBuilder;
import com.grafana.foundation.dashboardv2alpha1.TargetBuilder;
import com.grafana.foundation.table.TableVizConfigKindBuilder;
import com.grafana.foundation.units.Constants;

public class Disk {
        public static PanelBuilder diskIOTimeseries() {
                return new PanelBuilder<>()
                                .title("Disk I/O")
                                .visualization(Common.defaultTimeseries()
                                                .fillOpacity(0.0)
                                                .unit(Constants.BytesPerSecondSI)
                                                .overrideByRegexp("/ io time/",
                                                                List.of(new DynamicConfigValue("unit",
                                                                                Constants.PercentUnit))))
                                .data(new QueryGroupBuilder()
                                                .targets(List.of(
                                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                                "rate(node_disk_read_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                                "{{device}} read")),
                                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                                "rate(node_disk_written_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                                "{{device}} written")),
                                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                                "rate(node_disk_io_time_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                                "{{device}} IO time")))));

        }

        public static PanelBuilder diskSpaceUsageTable() {
                return new PanelBuilder<>()
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
                                                                new DynamicConfigValue("unit", Constants.PercentUnit),
                                                                new DynamicConfigValue("custom.cellOptions",
                                                                                Map.of("mode", "gradient", "type",
                                                                                                "gauge")),
                                                                new DynamicConfigValue("min", 0),
                                                                new DynamicConfigValue("max", 1))))
                                .data(new QueryGroupBuilder()
                                                .targets(List.of(
                                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                                "rate(node_disk_read_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                                "{{device}} read")),
                                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                                "rate(node_disk_written_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                                "{{device}} written")),
                                                                new TargetBuilder().query(Common.basicPrometheusQuery(
                                                                                "rate(node_disk_io_time_seconds_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"\"}[$__rate_interval])",
                                                                                "{{device}} read")))));

        }
}
