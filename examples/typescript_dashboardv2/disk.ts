import {VisualizationBuilder as TableVizBuilder} from "../../generated/typescript/src/table";
import {basicPrometheusQuery, defaultTimeseries, tablePrometheusQuery} from "./common";
import {
    FieldTextAlignment,
    TableCellDisplayMode,
    TableCellHeight,
    TableFooterOptionsBuilder
} from "../../generated/typescript/src/common";
import {
    PanelBuilder,
    QueryGroupBuilder,
    TargetBuilder,
    ThresholdsConfigBuilder,
    ThresholdsMode,
    TransformationBuilder,
} from "../../generated/typescript/src/dashboardv2beta1";

export const diskIOTimeseries = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Disk I/O")
        .visualization(defaultTimeseries()
            .fillOpacity(0)
            .unit("Bps")
            .thresholds(
                new ThresholdsConfigBuilder()
                    .mode(ThresholdsMode.Absolute)
                    .steps([
                        {value: null, color: "green"},
                        {value: 80.0, color: "red"},
                    ])
            )
            .overrideByRegexp("/ io time/", [
                {id: "unit", value: "percentunit"},
            ])
        )
        .data(new QueryGroupBuilder()
            .targets([
                new TargetBuilder().query(basicPrometheusQuery(`rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} read").refId("A")),
                new TargetBuilder().query(basicPrometheusQuery(`rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} written").refId("B")),
                new TargetBuilder().query(basicPrometheusQuery(`rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} IO time").refId("C")),
            ])
        );
};

export const diskSpaceUsageTable = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Disk Space Usage")
        .visualization(new TableVizBuilder()
            .align(FieldTextAlignment.Auto)
            .cellOptions({type: TableCellDisplayMode.Auto})
            .cellHeight(TableCellHeight.Sm)
            .footer(new TableFooterOptionsBuilder().countRows(false).reducer(["sum"]))
            .unit("decbytes")

            // Overrides configuration
            .overrideByName("Mounted on",[
                {id: "custom.width", value: 260},
            ])
            .overrideByName("Size", [
                {id: "custom.width", value: 93},
            ])
            .overrideByName("Used", [
                {id: "custom.width", value: 72},
            ])
            .overrideByName("Available", [
                {id: "custom.width", value: 88},
            ])
            .overrideByName("Used, %", [
                {id: "unit", value: "percentunit"},
                {
                    id: "custom.cellOptions",
                    value: {mode: "gradient", type: "gauge"}
                },
                {id: "max", value: 1},
                {id: "min", value: 0}
            ])
        )
        .data(new QueryGroupBuilder()
            .targets([
                new TargetBuilder().query(tablePrometheusQuery(`max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`).refId("A")),
                new TargetBuilder().query(tablePrometheusQuery(`max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`).refId("B")),
            ])
            // TODO: transformations are very clunky
            .transformation(new TransformationBuilder()
                .kind("groupBy")
                .id("groupBy")
                .options({
                    fields: {
                        "Value #A": {
                            aggregations: ["lastNotNull"],
                            operation: "aggregate"
                        },
                        "Value #B": {
                            aggregations: ["lastNotNull"],
                            operation: "aggregate"
                        },
                        mountpoint: {
                            aggregations: [],
                            operation: "groupby"
                        }
                    }
                })
            )
            .transformation(new TransformationBuilder()
                .kind("merge")
                .id("merge")
                .options({})
            )
            .transformation(new TransformationBuilder()
                .kind("calculateField") // is `kind` the type of transformation?
                .id("calculateField") // should this be different from `kind`?
                .options({
                    alias: "Used",
                    binary: {
                        left: "Value #A (lastNotNull)",
                        operator: "-",
                        reducer: "sum",
                        right: "Value #B (lastNotNull)"
                    },
                    mode: "binary",
                    reduce: {reducer: "sum"}
                })
            )
            .transformation(new TransformationBuilder()
                .kind("calculateField")
                .id("calculateField")
                .options({
                    alias: "Used, %",
                    binary: {
                        left: "Used",
                        operator: "/",
                        reducer: "sum",
                        right: "Value #A (lastNotNull)"
                    },
                    mode: "binary",
                    reduce: {"reducer": "sum"}
                })
            )
            .transformation(new TransformationBuilder()
                .kind("organize")
                .id("organize")
                .options({
                    excludeByName: {},
                    indexByName: {},
                    renameByName: {
                        "Value #A (lastNotNull)": "Size",
                        "Value #B (lastNotNull)": "Available",
                        mountpoint: "Mounted on"
                    }
                })
            )
            .transformation(new TransformationBuilder()
                .kind("sortBy")
                .id("sortBy")
                .options({
                    fields: {},
                    sort: [
                        {field: "Mounted on"}
                    ]
                })
            )
        );
};
