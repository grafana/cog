<?php

namespace App\Monitoring;

use Grafana\Foundation\Timeseries;

class Network
{
    public static function receivedTimeseries(): Timeseries\PanelBuilder
    {
        return Common::defaultTimeseries()
            ->title('Network Received')
            ->description('Network received (bits/s)')
            ->min(0)
            ->unit('bps')
            ->fillOpacity(0)
            ->withTarget(
                Common::basicPrometheusQuery('rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8', '{{ device }}'),
            );
    }

    public static function transmittedTimeseries(): Timeseries\PanelBuilder
    {
        return Common::defaultTimeseries()
            ->title('Network Transmitted')
            ->description('Network transmitted (bits/s)')
            ->min(0)
            ->unit('bps')
            ->fillOpacity(0)
            ->withTarget(
                Common::basicPrometheusQuery('rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8', '{{ device }}'),
            );
    }
}