<?php

namespace App\Monitoring;

use Grafana\Foundation\Common\StackingConfigBuilder;
use Grafana\Foundation\Common\StackingMode;
use Grafana\Foundation\Dashboard\Threshold;
use Grafana\Foundation\Dashboard\ThresholdsConfigBuilder;
use Grafana\Foundation\Dashboard\ThresholdsMode;
use Grafana\Foundation\Gauge;
use Grafana\Foundation\Timeseries;

class CPU
{
    public static function usageTimeseries(): Timeseries\PanelBuilder
    {
        $query = <<<'PROMQL'
(
 (1 - sum without (mode) (rate(node_cpu_seconds_total{job="integrations/raspberrypi-node", mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
 / ignoring(cpu) group_left
 count without (cpu, mode) (node_cpu_seconds_total{job="integrations/raspberrypi-node", mode="idle", instance="$instance"})
)
PROMQL;
        
        return Common::defaultTimeseries()
            ->title('CPU Usage')
            ->span(18)
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
            ->withTarget(
                Common::basicPrometheusQuery($query, '{{ cpu }}'),
            );
    }

    public static function temperatureGauge(): Gauge\PanelBuilder
    {
        return Common::defaultGauge()
            ->title('CPU Temperature')
            ->span(6)
            ->min(30)
            ->max(100)
            ->unit('celsius')
            ->thresholds(
                (new ThresholdsConfigBuilder())
                    ->mode(ThresholdsMode::absolute())
                    ->steps([
                        new Threshold(color: 'rgba(50, 172, 45, 0.97)'),
                        new Threshold(color: 'rgba(237, 129, 40, 0.89)', value: 65.0),
                        new Threshold(color: 'rgba(245, 54, 54, 0.9)', value: 85.0),
                    ])
            )
            ->withTarget(
                Common::basicPrometheusQuery('avg(node_hwmon_temp_celsius{job="integrations/raspberrypi-node", instance="$instance"})', ''),
            );
    }

    public static function loadAverageTimeseries(): Timeseries\PanelBuilder
    {
        return Common::defaultTimeseries()
            ->title('Load Average')
            ->span(18)
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
            ->withTarget(
                Common::basicPrometheusQuery('node_load1{job="integrations/raspberrypi-node", instance="$instance"}', '1m load average'),
            )
            ->withTarget(
                Common::basicPrometheusQuery('node_load5{job="integrations/raspberrypi-node", instance="$instance"}', '5m load average'),
            )
            ->withTarget(
                Common::basicPrometheusQuery('node_load15{job="integrations/raspberrypi-node", instance="$instance"}', '15m load average'),
            )
            ->withTarget(
                Common::basicPrometheusQuery('count(node_cpu_seconds_total{job="integrations/raspberrypi-node", instance="$instance", mode="idle"})', 'logical cores'),
            );
    }
}
