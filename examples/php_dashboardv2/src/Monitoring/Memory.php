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

class Memory
{
    public static function usageTimeseries(): PanelBuilder
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

        return (new PanelBuilder())
            ->title('Memory Usage')
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
                    ->decimals(2)
                    ->unit('bytes')
            )
            ->data((new QueryGroupBuilder())
                ->targets([
                    (new TargetBuilder())->query(Common::basicPrometheusQuery($query, 'Used'))->refId("A"),
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('node_memory_Buffers_bytes{job="integrations/raspberrypi-node", instance="$instance"}', 'Buffers'))->refId("B"),
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('node_memory_Cached_bytes{job="integrations/raspberrypi-node", instance="$instance"}', 'Cached'))->refId("C"),
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('node_memory_MemFree_bytes{job="integrations/raspberrypi-node", instance="$instance"}', 'Free'))->refId("D"),
                ])
            );
    }

    public static function usageGauge(): PanelBuilder
    {
        $query = <<<'PROMQL'
100 - (
    avg(node_memory_MemAvailable_bytes{job="integrations/raspberrypi-node", instance="$instance"}) /
    avg(node_memory_MemTotal_bytes{job="integrations/raspberrypi-node", instance="$instance"})
    * 100
)
PROMQL;


        return (new PanelBuilder())
            ->title('Memory Usage')
            ->visualization(
                Common::defaultGauge()
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
            )
            ->data((new QueryGroupBuilder())
                ->target(
                    (new TargetBuilder())->query(Common::basicPrometheusQuery($query, ''))->refId("A"),
                )
            );
    }
}
