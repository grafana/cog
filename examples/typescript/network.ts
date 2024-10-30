import {PanelBuilder as TimeseriesPanelBuilder} from "../../generated/typescript/src/timeseries";
import {basicPrometheusQuery, defaultTimeseries} from "./common";

export const networkReceivedTimeseries = (): TimeseriesPanelBuilder => {
    return defaultTimeseries()
        .title("Network Received")
        .description("Network received (bits/s)")
        .min(0)
        .unit("bps")
        .fillOpacity(0)
        .withTarget(
            basicPrometheusQuery(`rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}"),
        );
};

export const networkTransmittedTimeseries = (): TimeseriesPanelBuilder => {
    return defaultTimeseries()
        .title("Network Transmitted")
        .description("Network transmitted (bits/s)")
        .min(0)
        .unit("bps")
        .fillOpacity(0)
        .withTarget(
            basicPrometheusQuery(`rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}"),
        );
};
