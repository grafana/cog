<?php

namespace App\Monitoring;

use Grafana\Foundation\Dashboardv2beta1\PanelBuilder;
use Grafana\Foundation\Dashboardv2beta1\QueryGroupBuilder;
use Grafana\Foundation\Dashboardv2beta1\TargetBuilder;

class Network
{
    public static function receivedTimeseries(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Network Received')
            ->description('Network received (bits/s)')
            ->visualization(
                Common::defaultTimeseries()
                    ->min(0)
                    ->unit('bps')
                    ->fillOpacity(0)
            )
            ->data((new QueryGroupBuilder())
                ->target(
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8', '{{ device }}'))->refId("A"),
                )
            );
    }

    public static function transmittedTimeseries(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Network Transmitted')
            ->description('Network transmitted (bits/s)')
            ->visualization(
                Common::defaultTimeseries()
                    ->min(0)
                    ->unit('bps')
                    ->fillOpacity(0)
            )
            ->data((new QueryGroupBuilder())
                ->target(
                    (new TargetBuilder())->query(Common::basicPrometheusQuery('rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8', '{{ device }}'))->refId("A"),
                )
            );
    }
}
