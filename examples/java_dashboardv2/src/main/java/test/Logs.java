package test;

import com.grafana.foundation.dashboardv2beta1.PanelBuilder;
import com.grafana.foundation.dashboardv2beta1.QueryGroupBuilder;
import com.grafana.foundation.dashboardv2beta1.TargetBuilder;

public class Logs {
    public static PanelBuilder errorsInSystemLogs() {
        return new PanelBuilder()
                .title("Errors in system logs")
                .visualization(Common.defaultLogs())
                .data(new QueryGroupBuilder()
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{level=~\"err|crit|alert|emerg\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"))
                                .refId("A"))
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{filename=~\"/var/log/syslog*|/var/log/messages*\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"} |~\".+(?i)error(?-i).+\""))
                                .refId("B")));
    }

    public static PanelBuilder authLogs() {
        return new PanelBuilder()
                .title("Auth logs")
                .visualization(Common.defaultLogs())
                .data(new QueryGroupBuilder()
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{unit=\"ssh.service\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"))
                                .refId("A"))
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{filename=~\"/var/log/auth.log|/var/log/secure\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"))
                                .refId("B")));
    }

    public static PanelBuilder kernelLogs() {
        return new PanelBuilder()
                .title("Kernel logs")
                .visualization(Common.defaultLogs())
                .data(new QueryGroupBuilder()
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{transport=\"kernel\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"))
                                .refId("A"))
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{filename=\"/var/log/kern.log\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"))
                                .refId("B")));
    }

    public static PanelBuilder allSystemLogs() {
        return new PanelBuilder()
                .title("All system logs")
                .visualization(Common.defaultLogs())
                .data(new QueryGroupBuilder()
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{transport!=\"\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"))
                                .refId("A"))
                        .target(new TargetBuilder()
                                .query(Common.basicLokiQuery(
                                        "{filename!=\"\", job=\"integrations/raspberrypi-node\", instance=\"$instance\"}"))
                                .refId("B")));
    }
}
