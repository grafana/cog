<?php

namespace App\Monitoring;

use Grafana\Foundation\Common\StackingConfigBuilder;
use Grafana\Foundation\Common\StackingMode;
use Grafana\Foundation\Dashboardv2beta1\Threshold;
use Grafana\Foundation\Dashboardv2beta1\ThresholdsConfigBuilder;
use Grafana\Foundation\Dashboardv2beta1\ThresholdsMode;
use Grafana\Foundation\Dashboardv2beta1\PanelBuilder;
use Grafana\Foundation\Dashboardv2beta1\QueryGroupBuilder;
use Grafana\Foundation\Dashboardv2beta1\TargetBuilder;

class CPU
{
    public static function usageTimeseries(): PanelBuilder
    {
        $query = <<<'PROMQL'
(
 (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
 / ignoring(cpu) group_left
 count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
)
PROMQL;

        return (new PanelBuilder())
            ->title('CPU Usage')
            ->visualization(
                Common::defaultTimeseries()
                    ->stacking(
                        (new StackingConfigBuilder())
                            ->mode(StackingMode::normal())
                    )
                    ->thresholds(
                        (new ThresholdsConfigBuilder())
                            ->mode(ThresholdsMode::absolute())
                            ->steps([
                                new Threshold(color: 'green'),
                                new Threshold(color: 'red', value: 80.0),
                            ])
                    )
                    ->min(0)
                    ->max(1)
                    ->unit('percentunit')
            )
            ->data(
                (new QueryGroupBuilder())
                        ->target(
                            (new TargetBuilder())->query(Common::basicPrometheusQuery($query, '{{ cpu }}')->refId('A'))
                        )
            );
    }

    public static function temperatureGauge(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('CPU Temperature')
            ->visualization(
            Common::defaultGauge()
                ->thresholds(
                    (new ThresholdsConfigBuilder())
                        ->mode(ThresholdsMode::absolute())
                        ->steps([
                            new Threshold(color: 'rgba(50, 172, 45, 0.97)'),
                            new Threshold(color: 'rgba(237, 129, 40, 0.89)', value: 65.0),
                            new Threshold(color: 'rgba(245, 54, 54, 0.9)', value: 85.0),
                        ])
                )
                ->min(30)
                ->max(100)
                ->unit('celsius')
        )
        ->data(
            (new QueryGroupBuilder())
                    ->target(
                        Common::basicPrometheusQuery('avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})', '')
                    )
        );
    }

    public static function loadAverageTimeseries(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Load Average')
            ->visualization(
            Common::defaultTimeseries()
                ->thresholds(
                    (new ThresholdsConfigBuilder())
                        ->mode(ThresholdsMode::absolute())
                        ->steps([
                            new Threshold(color: 'green'),
                            new Threshold(color: 'red', value: 80.0),
                        ])
                )
                ->min(0)
                ->unit('short')
        )
        ->data(
            (new QueryGroupBuilder())
                    ->targets([
                        (new TargetBuilder())->query(Common::basicPrometheusQuery('node_load1{job="integrations/raspberrypi-node", instance="$instance"}', '1m load average'))->refId("A"),
                        (new TargetBuilder())->query(Common::basicPrometheusQuery('node_load5{job="integrations/raspberrypi-node", instance="$instance"}', '5m load average'))->refId("B"),
                        (new TargetBuilder())->query(Common::basicPrometheusQuery('node_load15{job="integrations/raspberrypi-node", instance="$instance"}', '15m load average'))->refId("C"),
                        (new TargetBuilder())->query(Common::basicPrometheusQuery('count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})', 'logical cores'))->refId("D"),
                    ])
        );
    }
}
