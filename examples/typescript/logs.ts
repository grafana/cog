import {PanelBuilder as LogsPanelBuilder} from "../../generated/typescript/src/logs";
import {basicLokiQuery, defaultLogs} from "./common";

export const errorsInSystemLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("Errors in system logs")
        .withTarget(
            basicLokiQuery(`{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}`),
        )
        .withTarget(
            basicLokiQuery(`{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"`),
        );
};

export const authLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("Auth logs")
        .withTarget(
            basicLokiQuery(`{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}`),
        )
        .withTarget(
            basicLokiQuery(`{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}`),
        );
};

export const kernelLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("Kernel logs")
        .withTarget(
            basicLokiQuery(`{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}`),
        )
        .withTarget(
            basicLokiQuery(`{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}`),
        );
};

export const allSystemLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("All system logs")
        .withTarget(
            basicLokiQuery(`{transport!="", job="integrations/raspberrypi-node", instance="$instance"}`),
        )
        .withTarget(
            basicLokiQuery(`{filename!="", job="integrations/raspberrypi-node", instance="$instance"}`),
        );
};
