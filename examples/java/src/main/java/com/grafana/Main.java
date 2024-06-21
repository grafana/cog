package com.grafana;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.dashboard.*;

import java.util.List;

public class Main {
    public static void main(String[] args) {
        Dashboard dashboard = new Dashboard.Builder("[TEST] Node Exporter / Raspberry").
                Uid("test-dashboard-raspberry").
                Tags(List.of("generated", "raspberrypi-node-integration")).
                Refresh("30s").
                Time(new DashboardDashboardTime.Builder().From("now-30m").To("now")).
                Timezone("browser").
                Timepicker(new TimePickerConfig.Builder().
                        RefreshIntervals(List.of("5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d")).
                        TimeOptions(List.of("5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"))
                ).
                Tooltip(DashboardCursorSync.CROSSHAIR).
                WithVariable(datasourceVariable()).
                WithVariable(queryVariable()).
                // CPU
                WithRow(new RowPanel.Builder("CPU")).
                WithPanel(CPU.cpuUsageTimeseries()).
                WithPanel(CPU.cpuTemperatureGauge()).
                WithPanel(CPU.loadAverageTimeseries()).
                // Memory
                WithRow(new RowPanel.Builder("Memory")).
                WithPanel(Memory.memoryUsageTimeseries()).
                WithPanel(Memory.memoryUsageGauge()).
                // Disk
                WithRow(new RowPanel.Builder("Disk")).
                WithPanel(Disk.diskIOTimeseries()).
                WithPanel(Disk.diskSpaceUsageTable()).
                // Network
                WithRow(new RowPanel.Builder("Network")).
                WithPanel(Network.networkReceivedTimeseries()).
                WithPanel(Network.networkTransmittedTimeseries()).
                // Logs
                WithRow(new RowPanel.Builder("Logs")).
                WithPanel(Logs.errorsInSystemLogs()).
                WithPanel(Logs.authLogs()).
                WithPanel(Logs.kernelLogs()).
                WithPanel(Logs.allSystemLogs()).
                build();
        try {
            System.out.println(dashboard.toJSON());
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    private static Builder<VariableModel> datasourceVariable() {
        VariableOption current = new VariableOption();
        current.selected = true;
        StringOrArrayOfString text = new StringOrArrayOfString();
        text.string = "grafanacloud-potatopi-prom";
        current.text = text;
        StringOrArrayOfString value = new StringOrArrayOfString();
        value.string = "grafanacloud-prom";
        current.value = value;

        return new VariableModel.DatasourceVariableBuilder("datasource").
                Label("Data Source").
                Hide(VariableHide.DONT_HIDE).
                Type("prometheus").
                Current(current);
    }

    private static Builder<VariableModel> queryVariable() {
        VariableOption current = new VariableOption();
        current.selected = false;
        StringOrArrayOfString text = new StringOrArrayOfString();
        text.string = "potato";
        current.text = text;
        StringOrArrayOfString value = new StringOrArrayOfString();
        value.string = "potato";
        current.value = value;

        StringOrMap query = new StringOrMap();
        query.string = "label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)";

        DataSourceRef datasource = new DataSourceRef();
        datasource.uid = "$datasource";
        datasource.type = "prometheus";

        return new VariableModel.QueryVariableBuilder("instance").
                Label("Instance").
                Hide(VariableHide.DONT_HIDE).
                Refresh(VariableRefresh.ON_TIME_RANGE_CHANGED).
                Query(query).
                Datasource(datasource).
                Current(current).
                Sort(VariableSort.DISABLED);
    }

}