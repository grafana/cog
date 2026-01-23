package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/cog/plugins"
	dashboard "github.com/grafana/cog/generated/go/dashboardv2beta1"
	"github.com/grafana/cog/generated/go/prometheus"
	"github.com/grafana/cog/generated/go/resource"
)

func dashboardBuilder() []byte {
	builder := dashboard.NewDashboardBuilder("[TEST] Node Exporter / Raspberry").
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
		// Tabs layout
		TabsLayout(dashboard.NewTabsLayoutBuilder().
			Tab(dashboard.NewTabsLayoutTabBuilder("CPU").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("cpu_usage")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("cpu_temp")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("load_avg")),
				),
			).
			Tab(dashboard.NewTabsLayoutTabBuilder("Memory").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("mem_usage")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("mem_usage_current")),
				),
			).
			Tab(dashboard.NewTabsLayoutTabBuilder("Disk").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("disk_io")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("disk_usage")),
				),
			).
			Tab(dashboard.NewTabsLayoutTabBuilder("Network").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("network_in")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("network_out")),
				),
			).
			Tab(dashboard.NewTabsLayoutTabBuilder("Logs").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("sys_error_logs")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("auth_logs")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("kernel_logs")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("all_sys_logs")),
				),
			),
		)
	// Rows layout
	/*
		RowsLayout(dashboard.NewRowsLayoutBuilder().
			Row(dashboard.NewRowsLayoutRowBuilder("CPU").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("cpu_usage")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("cpu_temp")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("load_avg")),
				),
			).
			Row(dashboard.NewRowsLayoutRowBuilder("Memory").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("mem_usage")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("mem_usage_current")),
				),
			).
			Row(dashboard.NewRowsLayoutRowBuilder("Disk").
				Collapse(true).
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("disk_io")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("disk_usage")),
				),
			).
			Row(dashboard.NewRowsLayoutRowBuilder("Network").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("network_in")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("network_out")),
				),
			).
			Row(dashboard.NewRowsLayoutRowBuilder("Logs").
				AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
					Item(dashboard.NewAutoGridLayoutItemBuilder("sys_error_logs")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("auth_logs")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("kernel_logs")).
					Item(dashboard.NewAutoGridLayoutItemBuilder("all_sys_logs")),
				),
			),
		)
	*/
	// Auto grid layout
	/*
		AutoGridLayout(dashboard.NewAutoGridLayoutBuilder().
			Item(dashboard.NewAutoGridLayoutItemBuilder("cpu_usage")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("cpu_temp")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("load_avg")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("mem_usage")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("mem_usage_current")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("disk_io")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("disk_usage")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("network_in")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("network_out")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("sys_error_logs")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("auth_logs")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("kernel_logs")).
			Item(dashboard.NewAutoGridLayoutItemBuilder("all_sys_logs")),
		)
	*/
	// "Manual" grid layout
	/*
		GridLayout(dashboard.NewGridLayoutBuilder().
			Row(dashboard.NewGridLayoutRowBuilder("CPU")).
			Item(dashboard.NewGridLayoutItemBuilder("cpu_usage").
				X(0). // TODO: X/Y calculations based on height and width? Or we leave it to the user to handle since we have the "AutoGrid" layout?
				Y(0).
				Height(200).
				Width(200),
			).
			Item(dashboard.NewGridLayoutItemBuilder("cpu_temp")).
			Item(dashboard.NewGridLayoutItemBuilder("load_avg")).
			Row(dashboard.NewGridLayoutRowBuilder("Memory")).
			Row(dashboard.NewGridLayoutRowBuilder("Disk")).
			Row(dashboard.NewGridLayoutRowBuilder("Network")).
			Row(dashboard.NewGridLayoutRowBuilder("Logs")),
		)
	*/

	sampleDashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}

	manifest := resource.Manifest{
		ApiVersion: resource.DashboardV2Beta1,
		Kind:       resource.DashboardKind,
		Metadata: resource.Metadata{
			Name: "test-v2-go",
		},
		Spec: sampleDashboard,
	}

	manifestJson, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		panic(err)
	}

	return manifestJson
}

func main() {
	// Required to correctly unmarshal panels and dataqueries
	plugins.RegisterDefaultPlugins()

	dashboardAsBytes := dashboardBuilder()

	fmt.Println(string(dashboardAsBytes))
}
