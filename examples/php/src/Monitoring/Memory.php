<?php

namespace App\Monitoring;

use Grafana\Foundation\Common\StackingConfigBuilder;
use Grafana\Foundation\Common\StackingMode;
use Grafana\Foundation\Dashboard\Threshold;
use Grafana\Foundation\Dashboard\ThresholdsConfigBuilder;
use Grafana\Foundation\Dashboard\ThresholdsMode;
use Grafana\Foundation\Gauge;
use Grafana\Foundation\Timeseries;

class Memory
{
    public static function usageTimeseries(): Timeseries\PanelBuilder
    {
        $query = <<<'PROMQL'
(
  node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"}
  -
  node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}
  -
  node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}
  -
  node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}
)
PROMQL;
        
        return Common::defaultTimeseries()
            ->title('Memory Usage')
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
            ->decimals(2)
            ->unit('bytes')
            ->withTarget(
                Common::basicPrometheusQuery($query, 'Used'),
            )
            ->withTarget(
                Common::basicPrometheusQuery('node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}', 'Buffers'),
            )
            ->withTarget(
                Common::basicPrometheusQuery('node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}', 'Cached'),
            )
            ->withTarget(
                Common::basicPrometheusQuery('node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}', 'Free'),
            );
    }

    public static function usageGauge(): Gauge\PanelBuilder
    {
        $query = <<<'PROMQL'
100 - (
    avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
    avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
    * 100
)
PROMQL;

        return Common::defaultGauge()
            ->title('Memory Usage')
            ->span(6)
            ->min(30)
            ->max(100)
            ->unit('percent')
            ->thresholds(
                (new ThresholdsConfigBuilder())
                    ->mode(ThresholdsMode::absolute())
                    ->steps([
                        new Threshold(color: 'rgba(50, 172, 45, 0.97)'),
                        new Threshold(color: 'rgba(237, 129, 40, 0.89)', value: 80.0),
                        new Threshold(color: 'rgba(245, 54, 54, 0.9)', value: 90.0),
                    ])
            )
            ->withTarget(
                Common::basicPrometheusQuery($query, ''),
            );
    }
}