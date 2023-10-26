package main

import (
	"github.com/grafana/cog/generated/dashboard"
	"github.com/grafana/cog/generated/logs"
)

func errorsInSystemLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("Errors in system logs").
		Targets([]dashboard.Target{
			basicLokiQuery(`{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"`),
		})
}

func authLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("Auth logs").
		Targets([]dashboard.Target{
			basicLokiQuery(`{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}`),
		})
}

func kernelLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("Kernel logs").
		Targets([]dashboard.Target{
			basicLokiQuery(`{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}`),
		})
}
func allSystemLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("All system logs").
		Targets([]dashboard.Target{
			basicLokiQuery(`{transport!="", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename!="", job="integrations/raspberrypi-node", instance="$instance"}`),
		})
}
