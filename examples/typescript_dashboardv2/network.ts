import {basicPrometheusQuery, defaultTimeseries} from "./common";
import {PanelBuilder, QueryGroupBuilder, TargetBuilder} from "../../generated/typescript/src/dashboardv2beta1";

export const networkReceivedTimeseries = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Network Received")
        .description("Network received (bits/s)")
        .visualization(defaultTimeseries()
            .min(0)
            .unit("bps")
            .fillOpacity(0)
        )
        .data(new QueryGroupBuilder().target(
            new TargetBuilder().query(basicPrometheusQuery(`rate(node_network_receive_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}").refId("A")),
        ));
};

export const networkTransmittedTimeseries = (): PanelBuilder => {
    return new PanelBuilder()
        .title("Network Transmitted")
        .description("Network transmitted (bits/s)")
        .visualization(defaultTimeseries()
            .min(0)
            .unit("bps")
            .fillOpacity(0)
        )
        .data(new QueryGroupBuilder().target(
            new TargetBuilder().query(basicPrometheusQuery(`rate(node_network_transmit_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!="lo"}[$__rate_interval]) * 8`, "{{ device }}").refId("B")),
        ));
};
