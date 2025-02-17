<?php

use App\Monitoring\CPU;
use App\Monitoring\Disk;
use App\Monitoring\Logs;
use App\Monitoring\Memory;
use App\Monitoring\Network;
use Grafana\Foundation\Common;
use Grafana\Foundation\Dashboard\DashboardBuilder;
use Grafana\Foundation\Dashboard\DashboardCursorSync;
use Grafana\Foundation\Dashboard\DataSourceRef;
use Grafana\Foundation\Dashboard\DatasourceVariableBuilder;
use Grafana\Foundation\Dashboard\QueryVariableBuilder;
use Grafana\Foundation\Dashboard\RowBuilder;
use Grafana\Foundation\Dashboard\TimePickerBuilder;
use Grafana\Foundation\Dashboard\VariableHide;
use Grafana\Foundation\Dashboard\VariableOption;
use Grafana\Foundation\Dashboard\VariableRefresh;
use Grafana\Foundation\Dashboard\VariableSort;

require_once __DIR__.'/vendor/autoload.php';

$builder = (new DashboardBuilder(title: '[TEST] Node Exporter / Raspberry'))
    ->uid('test-dashboard-raspberry')
    ->tags(['generated', 'raspberrypi-node-integration'])
    ->refresh('30s')
    ->time('now-30m', 'now')
    ->timezone(Common\Constants::TIME_ZONE_BROWSER)
    ->timepicker(
        (new TimePickerBuilder())->refreshIntervals(['5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h', '2h', '1d'])
    )
    ->tooltip(DashboardCursorSync::crosshair())
    // 'Data Source' variable
    ->withVariable(
        (new DatasourceVariableBuilder('datasource'))
            ->label('Data Source')
            ->hide(VariableHide::dontHide())
            ->type('prometheus')
            ->current(new VariableOption(
                selected: true,
                text: 'grafanacloud-potatopi-prom',
                value: 'grafanacloud-prom',
            ))
    )
    // "Instance" variable
    ->withVariable(
        (new QueryVariableBuilder('instance'))
            ->label('Instance')
            ->hide(VariableHide::dontHide())
            ->refresh(VariableRefresh::onTimeRangeChanged())
            ->query('label_values(node_uname_info{job="integrations/raspberrypi-node", sysname!="Darwin"}, instance)')
            ->datasource(new DataSourceRef(uid: '$datasource', type: 'prometheus'))
            ->current(new VariableOption(
                selected: true,
                text: 'potato',
                value: 'potato',
            ))
            ->sort(VariableSort::disabled())
    )
    // CPU
    ->withRow(new RowBuilder('CPU'))
    ->withPanel(CPU::usageTimeseries())
    ->withPanel(CPU::temperatureGauge())
    ->withPanel(CPU::loadAverageTimeseries())
    // Memory
    ->withRow(new RowBuilder('Memory'))
    ->withPanel(Memory::usageTimeseries())
    ->withPanel(Memory::usageGauge())
    // Disk
    ->withRow(new RowBuilder('Disk'))
    ->withPanel(Disk::ioTimeseries())
    ->withPanel(Disk::spaceUsageTable())
    // Network
    ->withRow(new RowBuilder('Network'))
    ->withPanel(Network::receivedTimeseries())
    ->withPanel(Network::transmittedTimeseries())
    // Logs
    ->withRow(new RowBuilder('Logs'))
    ->withPanel(Logs::errorsInSystemLogs())
    ->withPanel(Logs::authLogs())
    ->withPanel(Logs::kernelLogs())
    ->withPanel(Logs::allSystemLogs())
;

echo(json_encode($builder->build(), JSON_PRETTY_PRINT).PHP_EOL);
