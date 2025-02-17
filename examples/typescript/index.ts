import {
    DashboardBuilder,
    DashboardCursorSync,
    DatasourceVariableBuilder,
    QueryVariableBuilder,
    RowBuilder,
    TimePickerBuilder,
    VariableHide,
    VariableRefresh,
    VariableSort,
} from "../../generated/typescript/src/dashboard";
import {cpuTemperatureGauge, cpuUsageTimeseries, loadAverageTimeseries} from "./cpu";
import {memoryUsageGauge, memoryUsageTimeseries} from "./memory";
import {diskIOTimeseries, diskSpaceUsageTable} from "./disk";
import {networkReceivedTimeseries, networkTransmittedTimeseries} from "./network";
import {allSystemLogs, authLogs, errorsInSystemLogs, kernelLogs} from "./logs";

const builder = new DashboardBuilder("[TEST] Node Exporter / Raspberry")
    .uid("test-dashboard-raspberry")
    .tags(["generated", "raspberrypi-node-integration"])

    .refresh("30s")
    .time({from: "now-30m", to: "now"})
    .timezone("browser")

    .timepicker(
        new TimePickerBuilder().refreshIntervals(["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"]),
    )

    .tooltip(DashboardCursorSync.Crosshair)

    // "Data Source" variable
    .withVariable(
        new DatasourceVariableBuilder("datasource")
            .label("Data Source")
            .hide(VariableHide.DontHide)
            .type("prometheus")
            .current({
                selected: true,
                text: "grafanacloud-potatopi-prom",
                value: "grafanacloud-prom",
            })
    )
    // "Instance" variable
    .withVariable(
        new QueryVariableBuilder("instance")
            .name("instance")
            .label("Instance")
            .hide(VariableHide.DontHide)
            .refresh(VariableRefresh.OnTimeRangeChanged)
            .query('label_values(node_uname_info{job="integrations/raspberrypi-node", sysname!="Darwin"}, instance)')
            .datasource({
                "type": "prometheus",
                "uid": "$datasource"
            })
            .current({
                selected: false,
                text: "potato",
                value: "potato"
            })
            .sort(VariableSort.Disabled)
    )

    .withRow(new RowBuilder("CPU"))
    .withPanel(cpuUsageTimeseries())
    .withPanel(cpuTemperatureGauge())
    .withPanel(loadAverageTimeseries())

    .withRow(new RowBuilder("Memory"))
    .withPanel(memoryUsageTimeseries())
    .withPanel(memoryUsageGauge())

    .withRow(new RowBuilder("Disk"))
    .withPanel(diskIOTimeseries())
    .withPanel(diskSpaceUsageTable())

    .withRow(new RowBuilder("Network"))
    .withPanel(networkReceivedTimeseries())
    .withPanel(networkTransmittedTimeseries())

    .withRow(new RowBuilder("Logs"))
    .withPanel(errorsInSystemLogs())
    .withPanel(authLogs())
    .withPanel(kernelLogs())
    .withPanel(allSystemLogs())
;

console.log(JSON.stringify(builder.build(), null, 2));

