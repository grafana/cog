<?php

namespace App\Monitoring;

use Grafana\Foundation\Dashboardv2beta1\PanelBuilder;
use Grafana\Foundation\Dashboardv2beta1\QueryGroupBuilder;
use Grafana\Foundation\Dashboardv2beta1\TargetBuilder;

class Logs
{
    public static function errorsInSystemLogs(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Errors in system logs')
            ->visualization(Common::defaultLogs())
            ->data((new QueryGroupBuilder())
                ->targets([
                    (new TargetBuilder())->query(Common::basicLokiQuery('{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}'))->refId("A"),
                    (new TargetBuilder())->query(Common::basicLokiQuery('{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"'))->refId("B")
                ])
            );
    }

    public static function authLogs(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Auth logs')
            ->visualization(Common::defaultLogs())
            ->data((new QueryGroupBuilder())
                ->targets([
                    (new TargetBuilder())->query(Common::basicLokiQuery('{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}'))->refId("A"),
                    (new TargetBuilder())->query(Common::basicLokiQuery('{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}'))->refId("B")
                ])
            );
    }

    public static function kernelLogs(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('Kernel logs')
            ->visualization(Common::defaultLogs())
            ->data((new QueryGroupBuilder())
                ->targets([
                    (new TargetBuilder())->query(Common::basicLokiQuery('{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}'))->refId("A"),
                    (new TargetBuilder())->query(Common::basicLokiQuery('{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}'))->refId("B")
                ])
            );
    }

    public static function allSystemLogs(): PanelBuilder
    {
        return (new PanelBuilder())
            ->title('All system logs')
            ->visualization(Common::defaultLogs())
            ->data((new QueryGroupBuilder())
                ->targets([
                    (new TargetBuilder())->query(Common::basicLokiQuery('{transport!="", job="integrations/raspberrypi-node", instance="$instance"}'))->refId("A"),
                    (new TargetBuilder())->query(Common::basicLokiQuery('{filename!="", job="integrations/raspberrypi-node", instance="$instance"}'))->refId("B")
                ])
            );
    }
}
