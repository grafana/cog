package test;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.grafana.foundation.cog.Builder;
import com.grafana.foundation.dashboard.*;

import java.util.List;

public class Main {
    public static void main(String[] args) {
        Dashboard dashboard = new Dashboard.Builder("[TEST] Node Exporter / Raspberry").
                uid("test-dashboard-raspberry").
                tags(List.of("generated", "raspberrypi-node-integration")).
                refresh("30s").
                time(new DashboardDashboardTime.Builder().from("now-30m").to("now")).
                timezone("browser").
                timepicker(new TimePickerConfig.Builder().
                        refreshIntervals(List.of("5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d")).
                        timeOptions(List.of("5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"))
                ).
                tooltip(DashboardCursorSync.CROSSHAIR).
                withVariable(datasourceVariable()).
                withVariable(queryVariable()).
                // CPU
                withRow(new RowPanel.Builder("CPU")).
                withPanel(CPU.cpuUsageTimeseries()).
                withPanel(CPU.cpuTemperatureGauge()).
                withPanel(CPU.loadAverageTimeseries()).
                // Memory
                withRow(new RowPanel.Builder("Memory")).
                withPanel(Memory.memoryUsageTimeseries()).
                withPanel(Memory.memoryUsageGauge()).
                // Disk
                withRow(new RowPanel.Builder("Disk")).
                withPanel(Disk.diskIOTimeseries()).
                withPanel(Disk.diskSpaceUsageTable()).
                // Network
                withRow(new RowPanel.Builder("Network")).
                withPanel(Network.networkReceivedTimeseries()).
                withPanel(Network.networkTransmittedTimeseries()).
                // Logs
                withRow(new RowPanel.Builder("Logs")).
                withPanel(Logs.errorsInSystemLogs()).
                withPanel(Logs.authLogs()).
                withPanel(Logs.kernelLogs()).
                withPanel(Logs.allSystemLogs()).
                build();
        try {
            System.out.println(dashboard.toJSON());
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    private static Builder<VariableModel> datasourceVariable() {
        return new VariableModel.DatasourceVariableBuilder("datasource").
                label("Data Source").
                hide(VariableHide.DONT_HIDE).
                type("prometheus").
                current(new VariableOption(true, StringOrArrayOfString.createString("grafanacloud-potatopi-prom"), StringOrArrayOfString.createString("grafanacloud-prom")));
    }

    private static Builder<VariableModel> queryVariable() {
        return new VariableModel.QueryVariableBuilder("instance").
                label("Instance").
                hide(VariableHide.DONT_HIDE).
                refresh(VariableRefresh.ON_TIME_RANGE_CHANGED).
                query(StringOrMap.createString("label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)")).
                datasource(new DataSourceRef("$datasource", "prometheus")).
                current(new VariableOption(false, StringOrArrayOfString.createString("potato"), StringOrArrayOfString.createString("potato"))).
                sort(VariableSort.DISABLED);
    }

}
