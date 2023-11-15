import {
    DashboardBuilder,
    DashboardCursorSync,
    RowBuilder,
    TimePickerBuilder,
    VariableHide,
    VariableModelBuilder,
    VariableRefresh,
    VariableSort,
    VariableType
} from "../../generated/dashboard";
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
        new TimePickerBuilder()
            .refresh_intervals(["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"])
            .time_options(["5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"]),
    )

    .tooltip(DashboardCursorSync.Crosshair)

    // "Data Source" variable
    .withVariable(
        new VariableModelBuilder()
            .type(VariableType.Datasource)
            .name("datasource")
            .label("Data Source")
            .hide(VariableHide.DontHide)
            .refresh(VariableRefresh.OnDashboardLoad)
            .query("prometheus")
            .current({
                selected: true,
                text: "grafanacloud-potatopi-prom",
                value: "grafanacloud-prom",
            })
            .sort(VariableSort.Disabled)
    )
    // "Instance" variable
    .withVariable(
        new VariableModelBuilder()
            .type(VariableType.Query)
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

    .withRow(new RowBuilder("CPU").gridPos({h: 1, w: 24, x: 0, y: 0}))
    .withPanel(cpuUsageTimeseries().gridPos({h: 7, w: 18, x: 0, y: 0}))
    .withPanel(cpuTemperatureGauge().gridPos({h: 7, w: 6, x: 0, y: 0}))
    .withPanel(loadAverageTimeseries().gridPos({h: 7, w: 18, x: 0, y: 0}))

    .withRow(new RowBuilder("Memory").gridPos({h: 1, w: 24, x: 0, y: 0}))
    .withPanel(memoryUsageTimeseries().gridPos({h: 7, w: 18, x: 0, y: 0}))
    .withPanel(memoryUsageGauge().gridPos({h: 7, w: 6, x: 0, y: 0}))

    .withRow(new RowBuilder("Disk").gridPos({h: 1, w: 24, x: 0, y: 0}))
    .withPanel(diskIOTimeseries().gridPos({h: 7, w: 12, x: 0, y: 0}))
    .withPanel(diskSpaceUsageTable().gridPos({h: 7, w: 12, x: 0, y: 0}))

    .withRow(new RowBuilder("Network").gridPos({h: 1, w: 24, x: 0, y: 0}))
    .withPanel(networkReceivedTimeseries().gridPos({h: 7, w: 12, x: 0, y: 0}))
    .withPanel(networkTransmittedTimeseries().gridPos({h: 7, w: 12, x: 0, y: 0}))

    .withRow(new RowBuilder("Logs").gridPos({h: 1, w: 24, x: 0, y: 0}))
    .withPanel(errorsInSystemLogs().gridPos({h: 7, w: 24, x: 0, y: 0}))
    .withPanel(authLogs().gridPos({h: 7, w: 24, x: 0, y: 0}))
    .withPanel(kernelLogs().gridPos({h: 7, w: 24, x: 0, y: 0}))
    .withPanel(allSystemLogs().gridPos({h: 7, w: 24, x: 0, y: 0}))
;

console.log(JSON.stringify(builder.build(), null, 2));

