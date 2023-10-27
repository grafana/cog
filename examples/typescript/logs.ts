import {PanelBuilder as LogsPanelBuilder} from "../../generated/logs";
import {basicLokiQuery, defaultLogs} from "./common";

export const errorsInSystemLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("Errors in system logs")
        .targets([
            basicLokiQuery(`{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}`),
            basicLokiQuery(`{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"`),
        ]);
};

export const authLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("Auth logs")
        .targets([
            basicLokiQuery(`{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}`),
            basicLokiQuery(`{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}`),
        ]);
};

export const kernelLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("Kernel logs")
        .targets([
            basicLokiQuery(`{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}`),
            basicLokiQuery(`{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}`),
        ]);
};

export const allSystemLogs = (): LogsPanelBuilder => {
    return defaultLogs()
        .title("All system logs")
        .targets([
            basicLokiQuery(`{transport!="", job="integrations/raspberrypi-node", instance="$instance"}`),
            basicLokiQuery(`{filename!="", job="integrations/raspberrypi-node", instance="$instance"}`),
        ]);
};
