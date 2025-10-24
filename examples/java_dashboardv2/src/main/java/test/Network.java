package test;

import com.grafana.foundation.dashboardv2beta1.PanelBuilder;
import com.grafana.foundation.dashboardv2beta1.QueryGroupBuilder;
import com.grafana.foundation.dashboardv2beta1.TargetBuilder;
import com.grafana.foundation.units.Constants;

public class Network {

        public static PanelBuilder networkReceivedTimeseries() {
                return new PanelBuilder()
                                .title("Network Received")
                                .description("Network received (bits/s)")
                                .visualization(Common.defaultTimeSeries()
                                                .min(0.0)
                                                .unit(Constants.BitsPerSecondSI)
                                                .fillOpacity(0.0))
                                .data(new QueryGroupBuilder()
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "rate(node_network_receive_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"lo\"}[$__rate_interval]) * 8",
                                                                                "{{ device }}"))
                                                                .refId("A")));
        }

        public static PanelBuilder networkTransmittedTimeseries() {
                return new PanelBuilder()
                                .title("Network Transmitted")
                                .description("Network transmitted (bits/s)")
                                .visualization(Common.defaultTimeSeries()
                                                .min(0.0)
                                                .unit(Constants.BitsPerSecondSI)
                                                .fillOpacity(0.0))
                                .data(new QueryGroupBuilder()
                                                .target(new TargetBuilder()
                                                                .query(Common.basicPrometheusQuery(
                                                                                "rate(node_network_transmit_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"lo\"}[$__rate_interval]) * 8",
                                                                                "{{ device }}"))
                                                                .refId("A")));
        }
}
