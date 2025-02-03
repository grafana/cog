package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/cog/plugins"
	"github.com/grafana/cog/generated/go/common"
	dashboard "github.com/grafana/cog/generated/go/dashboardv2alpha0"
	"github.com/grafana/cog/generated/go/prometheus"
)

func dashboardBuilder() []byte {
	builder := dashboard.NewDashboardBuilder("[TEST] Node Exporter / Raspberry").
		//Uid("test-dashboard-raspberry"). // no more dashboard UID? (is it because the schema is just for a dashboard "spec"?)
		Tags([]string{"generated", "raspberrypi-node-integration"}).
		CursorSync(dashboard.DashboardCursorSyncCrosshair).
		TimeSettings(dashboard.NewTimeSettingsBuilder().
			AutoRefresh("30s").
			AutoRefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
			From("now-30m").
			To("now").
			Timezone("browser"),
		).
		// "Data Source" variable
		DatasourceVariable(dashboard.NewDatasourceVariableBuilder("datasource").
			Label("Data Source").
			Hide(dashboard.VariableHideDontHide).
			PluginId("prometheus").
			Current(dashboard.VariableOption{
				Selected: cog.ToPtr(true),
				Text:     dashboard.StringOrArrayOfString{String: cog.ToPtr("grafanacloud-potatopi-prom")},
				Value:    dashboard.StringOrArrayOfString{String: cog.ToPtr("grafanacloud-prom")},
			}),
		).
		// "Instance" variable
		QueryVariable(dashboard.NewQueryVariableBuilder("instance").
			Label("Instance").
			Hide(dashboard.VariableHideDontHide).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged).
			Query(prometheus.NewQueryBuilder().
				Expr("label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)"),
			).
			Datasource(common.DataSourceRef{
				Type: cog.ToPtr("prometheus"),
				Uid:  cog.ToPtr("$datasource"),
			}).
			Current(dashboard.VariableOption{
				Selected: cog.ToPtr(false),
				Text:     dashboard.StringOrArrayOfString{String: cog.ToPtr("potato")},
				Value:    dashboard.StringOrArrayOfString{String: cog.ToPtr("potato")},
			}).
			Sort(dashboard.VariableSortDisabled),
		).
		// CPU
		Panel("cpu_usage", cpuUsageTimeseries()).
		Panel("cpu_temp", cpuTemperatureGauge()).
		Panel("load_avg", loadAverageTimeseries()).
		// Memory
		Panel("mem_usage", memoryUsageTimeseries()).
		Panel("mem_usage_current", memoryUsageGauge()).
		// Disk
		Panel("disk_io", diskIOTimeseries()).
		Panel("disk_usage", diskSpaceUsageTable()).
		// Network
		Panel("network_in", networkReceivedTimeseries()).
		Panel("network_out", networkTransmittedTimeseries()).
		// Logs
		Panel("sys_error_logs", errorsInSystemLogs()).
		Panel("auth_logs", authLogs()).
		Panel("kernel_logs", kernelLogs()).
		Panel("all_sys_logs", allSystemLogs()).
		// Layout building
		Layout(dashboard.NewGridLayoutBuilder().
			Row(dashboard.NewGridLayoutRowBuilder("CPU")).
			Item(dashboard.NewGridLayoutItemBuilder("cpu_usage").
				X(0). // TODO: X/Y calculations based on height and width?
				Y(0).
				Height(200).
				Width(200),
			).
			Row(dashboard.NewGridLayoutRowBuilder("Memory")).
			Row(dashboard.NewGridLayoutRowBuilder("Disk")).
			Row(dashboard.NewGridLayoutRowBuilder("Network")).
			Row(dashboard.NewGridLayoutRowBuilder("Logs")),
		)

	sampleDashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}
	dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
	if err != nil {
		panic(err)
	}

	return dashboardJson
}

func main() {
	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	dashboardAsBytes := dashboardBuilder()

	fmt.Println(string(dashboardAsBytes))

	dash := dashboard.DashboardV2Spec{}
	if err := json.Unmarshal(dashboardAsBytes, &dash); err != nil {
		panic(err)
	} else {
		fmt.Printf("%#v\n", dash)
	}
}
