package main

import (
	logs "github.com/grafana/cog/generated/logs/panel"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func errorsInSystemLogs() *logs.Builder {
	return defaultLogs().
		Title("Errors in system logs").
		Targets([]types.Target{
			basicLokiQuery(`{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"`),
		})
}

func authLogs() *logs.Builder {
	return defaultLogs().
		Title("Auth logs").
		Targets([]types.Target{
			basicLokiQuery(`{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}`),
		})
}

func kernelLogs() *logs.Builder {
	return defaultLogs().
		Title("Kernel logs").
		Targets([]types.Target{
			basicLokiQuery(`{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}`),
		})
}
func allSystemLogs() *logs.Builder {
	return defaultLogs().
		Title("All system logs").
		Targets([]types.Target{
			basicLokiQuery(`{transport!="", job="integrations/raspberrypi-node", instance="$instance"}`),
			basicLokiQuery(`{filename!="", job="integrations/raspberrypi-node", instance="$instance"}`),
		})
}
