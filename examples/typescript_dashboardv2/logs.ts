import {basicLokiQuery, defaultLogs} from "./common";
import {PanelBuilder, QueryGroupBuilder, TargetBuilder} from "../../generated/typescript/src/dashboardv2beta1";

export const errorsInSystemLogs = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Errors in system logs")
        .visualization(defaultLogs())
        .data(new QueryGroupBuilder().targets([
            new TargetBuilder().query(basicLokiQuery(`{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}`).refId("A")),
            new TargetBuilder().query(basicLokiQuery(`{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"`).refId("B")),
        ]));
};

export const authLogs = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Auth logs")
        .visualization(defaultLogs())
        .data(new QueryGroupBuilder().targets([
            new TargetBuilder().query(basicLokiQuery(`{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}`).refId("A")),
            new TargetBuilder().query(basicLokiQuery(`{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}`).refId("B")),
        ]));
};

export const kernelLogs = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Kernel logs")
        .visualization(defaultLogs())
        .data(new QueryGroupBuilder().targets([
            new TargetBuilder().query(basicLokiQuery(`{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}`).refId("A")),
            new TargetBuilder().query(basicLokiQuery(`{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}`).refId("B")),
        ]));
};

export const allSystemLogs = (): PanelBuilder => {
    return new PanelBuilder()
        .title("All system logs")
        .visualization(defaultLogs())
        .data(new QueryGroupBuilder().targets([
            new TargetBuilder().query(basicLokiQuery(`{transport!="", job="integrations/raspberrypi-node", instance="$instance"}`).refId("A")),
            new TargetBuilder().query(basicLokiQuery(`{filename!="", job="integrations/raspberrypi-node", instance="$instance"}`).refId("B")),
        ]));
};
