<?php

use App\Monitoring\CPU;
use App\Monitoring\Disk;
use App\Monitoring\Logs;
use App\Monitoring\Memory;
use App\Monitoring\Network;
use Grafana\Foundation\Common;
use Grafana\Foundation\Dashboardv2\DashboardBuilder;
use Grafana\Foundation\Dashboardv2\DashboardCursorSync;
use Grafana\Foundation\Dashboardv2\DashboardV2Spec;
use Grafana\Foundation\Dashboardv2\DatasourceVariableBuilder;
use Grafana\Foundation\Dashboardv2\ElementReferenceBuilder;
use Grafana\Foundation\Dashboardv2\GridLayoutBuilder;
use Grafana\Foundation\Dashboardv2\GridLayoutItemBuilder;
use Grafana\Foundation\Dashboardv2\QueryVariableBuilder;
use Grafana\Foundation\Dashboardv2\VariableHide;
use Grafana\Foundation\Dashboardv2\VariableOption;
use Grafana\Foundation\Dashboardv2\VariableRefresh;
use Grafana\Foundation\Dashboardv2\VariableSort;
use Grafana\Foundation\Dashboardv2\TimeSettingsBuilder;

require_once __DIR__.'/vendor/autoload.php';

$builder = (new DashboardBuilder(title: '[TEST] Node Exporter / Raspberry'))
    ->tags(['generated', 'raspberrypi-node-integration'])
    ->cursorSync(DashboardCursorSync::crosshair())
    ->timeSettings(
        (new TimeSettingsBuilder())
        ->autoRefresh('30s')
        ->autoRefreshIntervals(["5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"])
        ->from('now-30m')
        ->to('now')
        ->timezone(Common\Constants::TIME_ZONE_BROWSER)
    )
    // 'Data Source' variable
    ->variable(
        (new DatasourceVariableBuilder('datasource'))
            ->label('Data Source')
            ->hide(VariableHide::dontHide())
            ->pluginId('prometheus')
            ->current(new VariableOption(
                selected: true,
                text: 'grafanacloud-potatopi-prom',
                value: 'grafanacloud-prom',
            ))
    )
    // "Instance" variable
    ->variable(
        (new QueryVariableBuilder('instance'))
            ->label('Instance')
            ->hide(VariableHide::dontHide())
            ->refresh(VariableRefresh::onTimeRangeChanged())
            ->query('label_values(node_uname_info{job="integrations/raspberrypi-node", sysname!="Darwin"}, instance)')
            ->datasource(new Common\DataSourceRef(uid: '$datasource', type: 'prometheus'))
            ->current(new VariableOption(
                selected: true,
                text: 'potato',
                value: 'potato',
            ))
            ->sort(VariableSort::disabled())
    )
    ->elements([
        // CPU
        "cpu_usage" => CPU::usageTimeseries(),
        "cpu_temp" => CPU::temperatureGauge(),
        "load_avg" => CPU::loadAverageTimeseries(),
        // Memory
        "mem_usage" => Memory::usageTimeseries(),
        "mem_usage_current" => Memory::usageGauge(),
        // Disk
        "disk_io" => Disk::ioTimeseries(),
        "disk_usage" => Disk::spaceUsageTable(),
        // Network
        "network_in" => Network::receivedTimeseries(),
        "network_out" => Network::transmittedTimeseries(),
        // Logs
        "sys_error_logs" => Logs::errorsInSystemLogs(),
        "auth_logs" => Logs::authLogs(),
        "kernel_logs" => Logs::kernelLogs(),
        "all_sys_logs" => Logs::allSystemLogs(),
    ])
    ->layout(
        // TODO: very clunky
        // TODO: rows? size?
        // TODO: automatic calculation of grid positions
        (new GridLayoutBuilder())
            ->item((new GridLayoutItemBuilder())->element((new ElementReferenceBuilder())->name("cpu_usage")))
    )
;

$jsonEncodedDashboard = json_encode($builder->build(), JSON_PRETTY_PRINT);

echo($jsonEncodedDashboard.PHP_EOL);

// Try decoding it.
$jsonDecodedAsArray = json_decode($jsonEncodedDashboard, true);
$dashboard = DashboardV2Spec::fromArray($jsonDecodedAsArray);
var_dump($dashboard->elements['cpu_usage']->spec->vizConfig);