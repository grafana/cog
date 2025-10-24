package main

import (
	dashboard "github.com/grafana/cog/generated/go/dashboardv2beta1"
)

func errorsInSystemLogs() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Errors in system logs").
		Visualization(defaultLogs()).
		Data(
			dashboard.NewQueryGroupBuilder().
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{level=~"err|crit|alert|emerg", job="integrations/raspberrypi-node", instance="$instance"}`)).RefId("A"),
				).
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{filename=~"/var/log/syslog*|/var/log/messages*", job="integrations/raspberrypi-node", instance="$instance"} |~".+(?i)error(?-i).+"`)).RefId("B"),
				),
		)
}

func authLogs() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Auth logs").
		Visualization(defaultLogs()).
		Data(
			dashboard.NewQueryGroupBuilder().
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{unit="ssh.service", job="integrations/raspberrypi-node", instance="$instance"}`)).RefId("A"),
				).
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{filename=~"/var/log/auth.log|/var/log/secure", job="integrations/raspberrypi-node", instance="$instance"}`)).RefId("B"),
				),
		)
}

func kernelLogs() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Kernel logs").
		Visualization(defaultLogs()).
		Data(
			dashboard.NewQueryGroupBuilder().
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{transport="kernel", job="integrations/raspberrypi-node", instance="$instance"}`)).RefId("A"),
				).
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{filename="/var/log/kern.log", job="integrations/raspberrypi-node", instance="$instance"}`)).RefId("B"),
				),
		)
}

func allSystemLogs() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("All system logs").
		Visualization(defaultLogs()).
		Data(
			dashboard.NewQueryGroupBuilder().
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{transport!="", job="integrations/raspberrypi-node", instance="$instance"}`)).RefId("A"),
				).
				Target(
					dashboard.NewTargetBuilder().Query(basicLokiQuery(`{filename!="", job="integrations/raspberrypi-node", instance="$instance"}`)).RefId("B"),
				),
		)
}
