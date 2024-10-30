package main

import (
	"github.com/grafana/cog/generated/go/timeseries"
)

func networkReceivedTimeseries() *timeseries.PanelBuilder {
	return defaultTimeseries().
		Title("Network Received").
		Description("Network received (bits/s)").
		Min(0).
		Unit("bps").
		FillOpacity(0).
		WithTarget(
			basicPrometheusQuery(`rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}"),
		)
}

func networkTransmittedTimeseries() *timeseries.PanelBuilder {
	return defaultTimeseries().
		Title("Network Transmitted").
		Description("Network transmitted (bits/s)").
		Min(0).
		Unit("bps").
		FillOpacity(0).
		WithTarget(
			basicPrometheusQuery(`rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}"),
		)
}
