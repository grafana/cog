package main

import (
	"github.com/grafana/cog/generated/go/logs"
)

func errorsInSystemLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("Errors in system logs").
		WithTarget(
			basicLokiQuery(`{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}`),
		).
		WithTarget(
			basicLokiQuery(`{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"`),
		)
}

func authLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("Auth logs").
		WithTarget(
			basicLokiQuery(`{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}`),
		).
		WithTarget(
			basicLokiQuery(`{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}`),
		)
}

func kernelLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("Kernel logs").
		WithTarget(
			basicLokiQuery(`{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}`),
		).
		WithTarget(
			basicLokiQuery(`{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}`),
		)
}

func allSystemLogs() *logs.PanelBuilder {
	return defaultLogs().
		Title("All system logs").
		WithTarget(
			basicLokiQuery(`{transport!="", job="integrations/raspberrypi-node", instance="$instance"}`),
		).
		WithTarget(
			basicLokiQuery(`{filename!="", job="integrations/raspberrypi-node", instance="$instance"}`),
		)
}
