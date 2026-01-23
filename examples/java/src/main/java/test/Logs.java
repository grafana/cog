package test;

import com.grafana.foundation.logs.LogsPanelBuilder;

public class Logs {
    public static LogsPanelBuilder errorsInSystemLogs() {
        return Common.defaultLogs().
                title("Errors in system logs").
                withTarget(Common.basicLokiQuery("{level=~\"err|crit|alert|emerg\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                withTarget(Common.basicLokiQuery("{filename=~\"/var/log/syslog*|/var/log/messages*\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"} |~\".+(?i)error(?-i).+\""));
    }
    public static LogsPanelBuilder authLogs() {
        return Common.defaultLogs().
                title("Auth logs").
                withTarget(Common.basicLokiQuery("{unit=\"ssh.service\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                withTarget(Common.basicLokiQuery("{filename=~\"/var/log/auth.log|/var/log/secure\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"));
    }
    public static LogsPanelBuilder kernelLogs() {
        return Common.defaultLogs().
                title("Kernel logs").
                withTarget(Common.basicLokiQuery("{transport=\"kernel\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                withTarget(Common.basicLokiQuery("{filename=\"/var/log/kern.log\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"));
    }
    public static LogsPanelBuilder allSystemLogs() {
        return Common.defaultLogs().
                title("All system logs").
                withTarget(Common.basicLokiQuery("{transport!=\"\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}")).
                withTarget(Common.basicLokiQuery("{filename!=\"\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"));
    }
}
