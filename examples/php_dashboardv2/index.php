<?php

use App\Monitoring\CPU;
use App\Monitoring\Disk;
use App\Monitoring\Logs;
use App\Monitoring\Memory;
use App\Monitoring\Network;
use App\Monitoring\Common as MonitoringCommon;
use Grafana\Foundation\Common;
use Grafana\Foundation\Dashboardv2beta1\AutoGridLayoutBuilder;
use Grafana\Foundation\Dashboardv2beta1\AutoGridLayoutItemBuilder;
use Grafana\Foundation\Dashboardv2beta1\DashboardBuilder;
use Grafana\Foundation\Dashboardv2beta1\DashboardCursorSync;
use Grafana\Foundation\Dashboardv2beta1\Dashboard;
use Grafana\Foundation\Dashboardv2beta1\DatasourceVariableBuilder;
use Grafana\Foundation\Dashboardv2beta1\DataSourceRefBuilder;
use Grafana\Foundation\Dashboardv2beta1\GridLayoutBuilder;
use Grafana\Foundation\Dashboardv2beta1\GridLayoutItemBuilder;
use Grafana\Foundation\Dashboardv2beta1\GridLayoutRowBuilder;
use Grafana\Foundation\Dashboardv2beta1\QueryVariableBuilder;
use Grafana\Foundation\Dashboardv2beta1\TimeSettingsBuilder;
use Grafana\Foundation\Dashboardv2beta1\TabsLayoutBuilder;
use Grafana\Foundation\Dashboardv2beta1\TabsLayoutTabBuilder;
use Grafana\Foundation\Dashboardv2beta1\VariableHide;
use Grafana\Foundation\Dashboardv2beta1\VariableOption;
use Grafana\Foundation\Dashboardv2beta1\VariableRefresh;
use Grafana\Foundation\Dashboardv2beta1\VariableSort;

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
            ->query(
                MonitoringCommon::basicPrometheusQuery('label_values(node_uname_info{job="integrations/raspberrypi-node", sysname!="Darwin"}, instance)', ''),
            )
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
    ->tabsLayout(
        (new TabsLayoutBuilder())
            ->tab(
                (new TabsLayoutTabBuilder("CPU"))
                    ->autoGridLayout(
                        (new AutoGridLayoutBuilder())
                            ->item(new AutoGridLayoutItemBuilder("cpu_usage"))
                            ->item(new AutoGridLayoutItemBuilder("cpu_temp"))
                            ->item(new AutoGridLayoutItemBuilder("load_avg"))
                    )
            )
            ->tab(
                (new TabsLayoutTabBuilder("Memory"))
                    ->autoGridLayout(
                        (new AutoGridLayoutBuilder())
                            ->item(new AutoGridLayoutItemBuilder("mem_usage"))
                            ->item(new AutoGridLayoutItemBuilder("mem_usage_current"))
                    )
            )
            ->tab(
                (new TabsLayoutTabBuilder("Disk"))
                    ->autoGridLayout(
                        (new AutoGridLayoutBuilder())
                            ->item(new AutoGridLayoutItemBuilder("disk_io"))
                            ->item(new AutoGridLayoutItemBuilder("disk_usage"))
                    )
            )
            ->tab(
                (new TabsLayoutTabBuilder("Network"))
                    ->autoGridLayout(
                        (new AutoGridLayoutBuilder())
                            ->item(new AutoGridLayoutItemBuilder("network_in"))
                            ->item(new AutoGridLayoutItemBuilder("network_out"))
                    )
            )
            ->tab(
                (new TabsLayoutTabBuilder("Logs"))
                    ->autoGridLayout(
                        (new AutoGridLayoutBuilder())
                            ->item(new AutoGridLayoutItemBuilder("sys_error_logs"))
                            ->item(new AutoGridLayoutItemBuilder("auth_logs"))
                            ->item(new AutoGridLayoutItemBuilder("kernel_logs"))
                            ->item(new AutoGridLayoutItemBuilder("all_sys_logs"))
                    )
            )
    );

$jsonEncodedDashboard = json_encode($builder->build(), JSON_PRETTY_PRINT);

echo($jsonEncodedDashboard.PHP_EOL);
