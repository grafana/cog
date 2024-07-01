import com.grafana.foundation.timeseries.PanelBuilder;

public class Network {

    public static PanelBuilder networkReceivedTimeseries() {
        return Common.defaultTimeSeries().
                Title("Network Received").
                Description("Network received (bits/s)").
                Min(0.0).
                Unit("bps").
                FillOpacity(0.0).
                WithTarget(
                        Common.basicPrometheusQuery("rate(node_network_receive_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"lo\"}[$__rate_interval]) * 8", "{{ device }}")
                );
    }
    public static PanelBuilder networkTransmittedTimeseries() {
        return Common.defaultTimeSeries().
                Title("Network Transmitted").
                Description("Network transmitted (bits/s)").
                Min(0.0).
                Unit("bps").
                FillOpacity(0.0).
                WithTarget(
                        Common.basicPrometheusQuery("rate(node_network_transmit_bytes_total{job=\"integrations/raspberrypi-node\", instance=\"$instance\", device!=\"lo\"}[$__rate_interval]) * 8", "{{ device }}")
                );
    }
}
