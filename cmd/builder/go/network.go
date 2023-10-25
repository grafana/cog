package main

import (
	timeseries "github.com/grafana/cog/generated/timeseries/panel"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func networkReceivedTimeseries() *timeseries.Builder {
	return defaultTimeseries().
		Title("Network Received").
		Description("Network received (bits/s)").
		Min(0).
		Unit("bps").
		FillOpacity(0).
		Targets([]types.Target{
			basicPrometheusQuery(`rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}"),
		})
}

func networkTransmittedTimeseries() *timeseries.Builder {
	return defaultTimeseries().
		Title("Network Transmitted").
		Description("Network transmitted (bits/s)").
		Min(0).
		Unit("bps").
		FillOpacity(0).
		Targets([]types.Target{
			basicPrometheusQuery(`rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}"),
		})
}
