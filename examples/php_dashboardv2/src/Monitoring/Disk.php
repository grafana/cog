<?php

namespace App\Monitoring;

use Grafana\Foundation\Common\FieldTextAlignment;
use Grafana\Foundation\Common\TableCellHeight;
use Grafana\Foundation\Common\TableFooterOptionsBuilder;
use Grafana\Foundation\Dashboardv2beta1\DataTransformerConfig;
use Grafana\Foundation\Dashboardv2beta1\DynamicConfigValue;
use Grafana\Foundation\Dashboardv2beta1\PanelBuilder;
use Grafana\Foundation\Dashboardv2beta1\QueryGroupBuilder;
use Grafana\Foundation\Dashboardv2beta1\TargetBuilder;
use Grafana\Foundation\Dashboardv2beta1\TransformationBuilder;
use Grafana\Foundation\Table;

class Disk
{
    public static function ioTimeseries(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Disk I/O')
            ->visualization(
                Common::defaultTimeseries()
                    ->unit('Bps')
                    ->overrideByRegexp('/ io time/', [
                        new DynamicConfigValue(id: 'unit', value: 'percentunit'),
                    ])
            )
            ->data((new QueryGroupBuilder())
                ->targets([
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', '{{ device }} read'))->refId("A"),
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', '{{ device }} written'))->refId("B"),
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', '{{ device }} IO time'))->refId("C"),
                ])
            );
    }

    public static function spaceUsageTable(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Disk Space Usage')
            ->visualization(
                (new Table\VisualizationBuilder())
                    ->align(FieldTextAlignment::auto())
                    ->unit('decbytes')
                    ->cellHeight(TableCellHeight::sm())
                    ->footer(
                        (new TableFooterOptionsBuilder())
                            ->countRows(false)
                            ->reducer(['sum'])
                    )
                    // Overrides
                    ->overrideByName('Mounted on', [
                        new DynamicConfigValue(id: 'custom.width', value: 260),
                    ])
                    ->overrideByName('Size', [
                        new DynamicConfigValue(id: 'custom.width', value: 93),
                    ])
                    ->overrideByName('Used', [
                        new DynamicConfigValue(id: 'custom.width', value: 72),
                    ])
                    ->overrideByName('Available', [
                        new DynamicConfigValue(id: 'custom.width', value: 88),
                    ])
                    ->overrideByName('Used, %', [
                        new DynamicConfigValue(id: 'unit', value: 'percentunit'),
                        new DynamicConfigValue(id: 'custom.cellOptions', value: [
                            'mode' => 'gradient',
                            'type' => 'gauge',
                        ]),
                        new DynamicConfigValue(id: 'min', value: 0),
                        new DynamicConfigValue(id: 'max', value: 1),
                    ])
            )
            ->data((new QueryGroupBuilder())
                ->targets([
                    (new TargetBuilder())->query(Common::tablePrometheusQuery('max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})'))->refId("A"),
                    (new TargetBuilder())->query(Common::tablePrometheusQuery('max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})'))->refId("B"),
                ])
                // Transformations
                ->transformation((new TransformationBuilder())
                    ->id('groupBy')
                    ->kind('groupBy')
                    ->options([
                        "fields" => [
                            "Value #A"=> [
                                "aggregations"=> ["lastNotNull"],
                                "operation"=> "aggregate",
                            ],
                            "Value #B"=> [
                                "aggregations"=> ["lastNotNull"],
                                "operation"=> "aggregate",
                            ],
                            "mountpoint"=> [
                                "aggregations"=> [],
                                "operation"=> "groupby",
                            ],
                        ]
                    ])
                )
                ->transformation((new TransformationBuilder())
                    ->id('merge')
                    ->kind('merge')
                    ->options([])
                )
                ->transformation((new TransformationBuilder())
                    ->id('calculateField')
                    ->kind('calculateField')
                    ->options([
                        "alias" => "Used",
                        "binary" => [
                            "left" => "Value #A (lastNotNull)",
                            "operator" => "-",
                            "reducer" => "sum",
                            "right" => "Value #B (lastNotNull)",
                        ],
                        "mode" => "binary",
                        "reduce" => ["reducer" => "sum"],
                    ])
                )
                ->transformation((new TransformationBuilder())
                    ->id('calculateField')
                    ->kind('calculateField')
                    ->options([
                        "alias" => "Used, %",
                        "binary" => [
                            "left" => "Used",
                            "operator" => "/",
                            "reducer" => "sum",
                            "right" => "Value #A (lastNotNull)",
                        ],
                        "mode" => "binary",
                        "reduce" => ["reducer" => "sum"],
                    ])
                )
                ->transformation((new TransformationBuilder())
                    ->id('organize')
                    ->kind('organize')
                    ->options([
                        "excludeByName" => [],
                        "indexByName" => [],
                        "renameByName" => [
                            "Value #A (lastNotNull)" => "Size",
                            "Value #B (lastNotNull)" => "Available",
                            "mountpoint" => "Mounted on",
                        ],
                    ])
                )
                ->transformation((new TransformationBuilder())
                    ->id('sortBy')
                    ->kind('sortBy')
                    ->options([
                        "fields" => [],
                        "sort" => [["field" => "Mounted on"]],
                    ])
                )
            );
    }
}
