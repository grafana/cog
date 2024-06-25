<?php

namespace App\Monitoring;

use Grafana\Foundation\Logs as SDKLogs;

class Logs
{
    public static function errorsInSystemLogs(): SDKLogs\PanelBuilder
    {
        return Common::defaultLogs()
            ->title('Errors in system logs')
            ->targets([
                Common::basicLokiQuery('{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}'),
                Common::basicLokiQuery('{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"'),
            ]);
    }

    public static function authLogs(): SDKLogs\PanelBuilder
    {
        return Common::defaultLogs()
            ->title('Auth logs')
            ->targets([
                Common::basicLokiQuery('{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}'),
                Common::basicLokiQuery('{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}'),
            ]);
    }

    public static function kernelLogs(): SDKLogs\PanelBuilder
    {
        return Common::defaultLogs()
            ->title('Kernel logs')
            ->targets([
                Common::basicLokiQuery('{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}'),
                Common::basicLokiQuery('{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}'),
            ]);
    }

    public static function allSystemLogs(): SDKLogs\PanelBuilder
    {
        return Common::defaultLogs()
            ->title('All system logs')
            ->targets([
                Common::basicLokiQuery('{transport!="", job="integrations/raspberrypi-node", instance="$instance"}'),
                Common::basicLokiQuery('{filename!="", job="integrations/raspberrypi-node", instance="$instance"}'),
            ]);
    }
}