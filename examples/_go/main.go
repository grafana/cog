package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/cog/plugins"
	"github.com/grafana/cog/generated/go/common"
	"github.com/grafana/cog/generated/go/dashboard"
)

func dashboardBuilder() []byte {
	builder := dashboard.NewDashboardBuilder("[TEST] Node Exporter / Raspberry").
		Uid("test-dashboard-raspberry").
		Tags([]string{"generated", "raspberrypi-node-integration"}).
		Refresh("30s").
		Time("now-30m", "now").
		Timezone(common.TimeZoneBrowser).
		Timepicker(
			dashboard.NewTimePickerBuilder().
				RefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}),
		).
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		// "Data Source" variable
		WithVariable(dashboard.NewDatasourceVariableBuilder("datasource").
			Label("Data Source").
			Hide(dashboard.VariableHideDontHide).
			Type("prometheus").
			Current(dashboard.VariableOption{
				Selected: cog.ToPtr(true),
				Text:     dashboard.StringOrArrayOfString{String: cog.ToPtr("grafanacloud-potatopi-prom")},
				Value:    dashboard.StringOrArrayOfString{String: cog.ToPtr("grafanacloud-prom")},
			}),
		).
		// "Instance" variable
		WithVariable(dashboard.NewQueryVariableBuilder("instance").
			Label("Instance").
			Hide(dashboard.VariableHideDontHide).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged).
			Query(dashboard.StringOrMap{
				String: cog.ToPtr("label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)"),
			}).
			Datasource(dashboard.DataSourceRef{
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
		WithRow(dashboard.NewRowBuilder("CPU")).
		WithPanel(cpuUsageTimeseries()).
		WithPanel(cpuTemperatureGauge()).
		WithPanel(loadAverageTimeseries()).
		// Memory
		WithRow(dashboard.NewRowBuilder("Memory")).
		WithPanel(memoryUsageTimeseries()).
		WithPanel(memoryUsageGauge()).
		// Disk
		WithRow(dashboard.NewRowBuilder("Disk")).
		WithPanel(diskIOTimeseries()).
		WithPanel(diskSpaceUsageTable()).
		// Network
		WithRow(dashboard.NewRowBuilder("Network")).
		WithPanel(networkReceivedTimeseries()).
		WithPanel(networkTransmittedTimeseries()).
		// Logs
		WithRow(dashboard.NewRowBuilder("Logs")).
		WithPanel(errorsInSystemLogs()).
		WithPanel(authLogs()).
		WithPanel(kernelLogs()).
		WithPanel(allSystemLogs())

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
}
