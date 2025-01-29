<?php

namespace App\Monitoring;

use Grafana\Foundation\Common\FieldTextAlignment;
use Grafana\Foundation\Common\TableCellHeight;
use Grafana\Foundation\Common\TableFooterOptionsBuilder;
use Grafana\Foundation\Dashboard\DataTransformerConfig;
use Grafana\Foundation\Dashboard\DynamicConfigValue;
use Grafana\Foundation\Dashboard\MatcherConfig;
use Grafana\Foundation\Table;
use Grafana\Foundation\Timeseries;

class Disk
{
    public static function ioTimeseries(): Timeseries\PanelBuilder
    {
        return Common::defaultTimeseries()
            ->title('Disk I/O')
            ->unit('Bps')
            ->targets([
                Common::basicPrometheusQuery('rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', '{{ device }} read'),
                Common::basicPrometheusQuery('rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', '{{ device }} written'),
                Common::basicPrometheusQuery('rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])', '{{ device }} IO time'),
            ])
            ->overrideByRegexp('/ io time/', [
                new DynamicConfigValue(id: 'unit', value: 'percentunit'),
            ]);
    }

    public static function spaceUsageTable(): Table\PanelBuilder
    {
        return (new Table\PanelBuilder())
            ->title('Disk Space Usage')
            ->align(FieldTextAlignment::auto())
            ->unit('decbytes')
            ->cellHeight(TableCellHeight::sm())
            ->footer(
                (new TableFooterOptionsBuilder())
                    ->countRows(false)
                    ->reducer(['sum'])
            )
            ->targets([
                Common::tablePrometheusQuery('max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})', 'A'),
                Common::tablePrometheusQuery('max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})', 'B'),
            ])
            // Transformations
            ->withTransformation(new DataTransformerConfig(
                id: 'groupBy',
                options: [
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
                ],
            ))
            ->withTransformation(new DataTransformerConfig(id: 'merge',options: []))
            ->withTransformation(new DataTransformerConfig(
                id: 'calculateField',
                options: [
                    "alias" => "Used",
                    "binary" => [
                        "left" => "Value #A (lastNotNull)",
                        "operator" => "-",
                        "reducer" => "sum",
                        "right" => "Value #B (lastNotNull)",
                    ],
                    "mode" => "binary",
                    "reduce" => ["reducer" => "sum"],
                ],
            ))
            ->withTransformation(new DataTransformerConfig(
                id: 'calculateField',
                options: [
                    "alias" => "Used, %",
                    "binary" => [
                        "left" => "Used",
                        "operator" => "/",
                        "reducer" => "sum",
                        "right" => "Value #A (lastNotNull)",
                    ],
                    "mode" => "binary",
                    "reduce" => ["reducer" => "sum"],
                ],
            ))
            ->withTransformation(new DataTransformerConfig(
                id: 'organize',
                options: [
                    "excludeByName" => [],
                    "indexByName" => [],
                    "renameByName" => [
                        "Value #A (lastNotNull)" => "Size",
                        "Value #B (lastNotNull)" => "Available",
                        "mountpoint" => "Mounted on",
                    ],
                ],
            ))
            ->withTransformation(new DataTransformerConfig(
                id: 'sortBy',
                options: [
                    "fields" => [],
                    "sort" => [["field" => "Mounted on"]],
                ],
            ))
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
            ]);
    }
}
