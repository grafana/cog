package test;

import com.grafana.foundation.logs.PanelBuilder;

public class Logs {
    public static PanelBuilder errorsInSystemLogs() {
        return Common.defaultLogs().
                Title("Errors in system logs").
                WithTarget(Common.basicLokiQuery("{level=~\"err|crit|alert|emerg\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                WithTarget(Common.basicLokiQuery("{filename=~\"/var/log/syslog*|/var/log/messages*\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"} |~\".+(?i)error(?-i).+\""));
    }
    public static PanelBuilder authLogs() {
        return Common.defaultLogs().
                Title("Auth logs").
                WithTarget(Common.basicLokiQuery("{unit=\"ssh.service\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                WithTarget(Common.basicLokiQuery("{filename=~\"/var/log/auth.log|/var/log/secure\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"));
    }
    public static PanelBuilder kernelLogs() {
        return Common.defaultLogs().
                Title("Kernel logs").
                WithTarget(Common.basicLokiQuery("{transport=\"kernel\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                WithTarget(Common.basicLokiQuery("{filename=\"/var/log/kern.log\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"));
    }
    public static PanelBuilder allSystemLogs() {
        return Common.defaultLogs().
                Title("All system logs").
                WithTarget(Common.basicLokiQuery("{transport!=\"\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                WithTarget(Common.basicLokiQuery("{filename!=\"\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"));
    }
}
