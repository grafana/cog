package main

import (
	dashboard "github.com/grafana/cog/generated/go/dashboardv2beta1"
	"github.com/grafana/cog/generated/go/units"
)

func networkReceivedTimeseries() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Network Received").
		Description("Network received (bits/s)").
		Visualization(
			defaultTimeseries().
				Min(0).
				Unit(units.BitsPerSecondSI).
				FillOpacity(0),
		).
		Data(
			dashboard.NewQueryGroupBuilder().Target(
				dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}")).RefId("A"),
			),
		)
}

func networkTransmittedTimeseries() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Network Transmitted").
		Description("Network transmitted (bits/s)").
		Visualization(
			defaultTimeseries().
				Min(0).
				Unit(units.BitsPerSecondSI).
				FillOpacity(0),
		).
		Data(
			dashboard.NewQueryGroupBuilder().Target(
				dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}")).RefId("A"),
			),
		)
}
