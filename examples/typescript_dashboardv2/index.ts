import {
    DashboardBuilder,
    DashboardCursorSync,
    DatasourceVariableBuilder,
    ElementReferenceBuilder,
    GridLayoutBuilder,
    GridLayoutItemBuilder,
    QueryVariableBuilder,
    TimeSettingsBuilder,
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
    .tags(["generated", "raspberrypi-node-integration"])
    .cursorSync(DashboardCursorSync.Crosshair)

    .timeSettings(
        new TimeSettingsBuilder()
            .autoRefresh("30s")
            .autoRefreshIntervals(["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"])
            .from("now-30m")
            .to("now")
            .timezone("browser")
    )

    // "Data Source" variable
    .withVariable(
        new DatasourceVariableBuilder("datasource")
            .label("Data Source")
            .hide(VariableHide.DontHide)
            //.type("prometheus")
            .current({
                selected: true,
                text: "grafanacloud-potatopi-prom",
                value: "grafanacloud-prom",
            })
    )
    // "Instance" variable
    .withVariable(
        new QueryVariableBuilder("instance")
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

    .elements({
        // CPU
        "cpu_usage": cpuUsageTimeseries(),
        "cpu_temp": cpuTemperatureGauge(),
        "load_avg": loadAverageTimeseries(),
        // Memory
        "mem_usage": memoryUsageTimeseries(),
        "mem_usage_current": memoryUsageGauge(),
        // Disk
        "disk_io": diskIOTimeseries(),
        "disk_usage": diskSpaceUsageTable(),
        // Network
        "network_in": networkReceivedTimeseries(),
        "network_out": networkTransmittedTimeseries(),
        // Logs
        "sys_error_logs": errorsInSystemLogs(),
        "auth_logs": authLogs(),
        "kernel_logs": kernelLogs(),
        "all_sys_logs": allSystemLogs(),
    })

    // TODO build layout
    // TODO: rows?
    .layout(new GridLayoutBuilder()
        .item(new GridLayoutItemBuilder()
            .element(new ElementReferenceBuilder().name("cpu_usage"))
        )
    )
;

console.log(JSON.stringify(builder.build(), null, 2));

