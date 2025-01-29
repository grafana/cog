import {PanelBuilder as TimeseriesPanelBuilder} from "../../generated/typescript/src/timeseries";
import {PanelBuilder as TablePanelBuilder} from "../../generated/typescript/src/table";
import {basicPrometheusQuery, defaultTimeseries, tablePrometheusQuery} from "./common";
import {
    FieldTextAlignment,
    TableCellDisplayMode,
    TableCellHeight,
    TableFooterOptionsBuilder
} from "../../generated/typescript/src/common";

export const diskIOTimeseries = (): TimeseriesPanelBuilder => {
    return defaultTimeseries()
        .title("Disk I/O")
        .fillOpacity(0)
        .unit("Bps")
        .withTarget(
            basicPrometheusQuery(`rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} read"),
        )
        .withTarget(
            basicPrometheusQuery(`rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} written"),
        )
        .withTarget(
            basicPrometheusQuery(`rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} IO time"),
        )
        .overrideByRegexp("/ io time/", [
            {id: "unit", value: "percentunit"},
        ]);
};

export const diskSpaceUsageTable = (): TablePanelBuilder => {
    return new TablePanelBuilder()
        .title("Disk Space Usage")
        .align(FieldTextAlignment.Auto)
        .cellOptions({type: TableCellDisplayMode.Auto})
        .cellHeight(TableCellHeight.Sm)
        .footer(new TableFooterOptionsBuilder().countRows(false).reducer(["sum"]))
        .unit("decbytes")
        .withTarget(
            tablePrometheusQuery(`max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "A"),
        )
        .withTarget(
            tablePrometheusQuery(`max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "B"),
        )

        // Transformations
        .withTransformation({
            id: "groupBy",
            options: {
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
            }
        })
        .withTransformation({
            id: "merge",
            options: {}
        })
        .withTransformation({
            id: "calculateField",
            options: {
                alias: "Used",
                binary: {
                    left: "Value #A (lastNotNull)",
                    operator: "-",
                    reducer: "sum",
                    right: "Value #B (lastNotNull)"
                },
                mode: "binary",
                reduce: {reducer: "sum"}
            }
        })
        .withTransformation({
            id: "calculateField",
            options: {
                alias: "Used, %",
                binary: {
                    left: "Used",
                    operator: "/",
                    reducer: "sum",
                    right: "Value #A (lastNotNull)"
                },
                mode: "binary",
                reduce: {"reducer": "sum"}
            }
        })
        .withTransformation({
            id: "organize",
            options: {
                excludeByName: {},
                indexByName: {},
                renameByName: {
                    "Value #A (lastNotNull)": "Size",
                    "Value #B (lastNotNull)": "Available",
                    mountpoint: "Mounted on"
                }
            }
        })
        .withTransformation({
            id: "sortBy",
            options: {
                fields: {},
                sort: [
                    {field: "Mounted on"}
                ]
            }
        })

        // Overrides configuration
        .overrideByName("Mounted on", [{id: "custom.width", value: 260}])
        .overrideByName("Size", [{id: "custom.width", value: 93}])
        .overrideByName("Used", [{id: "custom.width", value: 72}])
        .overrideByName("Available", [{id: "custom.width", value: 88}])
        .overrideByName("Used, %", [
            {id: "unit", value: "percentunit"},
            {
                id: "custom.cellOptions",
                value: {mode: "gradient", type: "gauge"}
            },
            {id: "max", value: 1},
            {id: "min", value: 0}
        ])
    ;
};
