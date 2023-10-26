package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/common"
	"github.com/grafana/cog/generated/dashboard"
)

func main() {
	builder := dashboard.NewDashboardBuilder("[TEST] Node Exporter / Raspberry").
		Uid("test-dashboard-raspberry").
		Tags([]string{"generated", "raspberrypi-node-integration"}).
		Refresh("30s").
		Time("now-30m", "now").
		Timezone(common.TimeZoneBrowser).
		Timepicker(
			dashboard.NewTimePickerBuilder().
				RefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
				TimeOptions([]string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"}),
		).
		Tooltip(dashboard.DashboardCursorSyncCrosshair).
		// TODO: we should have specific builders for every possible variable type
		// "Data Source" variable
		WithVariable(dashboard.NewVariableModelBuilder().
			Type(dashboard.VariableTypeDatasource).
			Name("datasource").
			Label("Data Source").
			Hide(dashboard.VariableHideDontHide).
			Refresh(dashboard.VariableRefreshOnDashboardLoad).
			Query(dashboard.StringOrAny{
				String: toPtr("prometheus"),
			}).
			Datasource(dashboard.DataSourceRef{
				Type: toPtr("prometheus"),
				Uid:  toPtr("$datasource"),
			}).
			Current(dashboard.VariableOption{
				Selected: toPtr(true),
				Text:     dashboard.StringOrArrayOfString{String: toPtr("grafanacloud-potatopi-prom")},
				Value:    dashboard.StringOrArrayOfString{String: toPtr("grafanacloud-prom")},
			}).
			Sort(dashboard.VariableSortDisabled),
		).
		// "Instance" variable
		WithVariable(dashboard.NewVariableModelBuilder().
			Type(dashboard.VariableTypeQuery).
			Name("instance").
			Label("Instance").
			Hide(dashboard.VariableHideDontHide).
			Refresh(dashboard.VariableRefreshOnTimeRangeChanged).
			Query(dashboard.StringOrAny{
				String: toPtr("label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)"),
			}).
			Datasource(dashboard.DataSourceRef{
				Type: toPtr("prometheus"),
				Uid:  toPtr("$datasource"),
			}).
			Current(dashboard.VariableOption{
				Selected: toPtr(false),
				Text:     dashboard.StringOrArrayOfString{String: toPtr("potato")},
				Value:    dashboard.StringOrArrayOfString{String: toPtr("potato")},
			}).
			Sort(dashboard.VariableSortDisabled),
		).
		// CPU
		WithRow(dashboard.NewRowPanelBuilder("CPU").GridPos(dashboard.GridPos{H: 1, W: 24})).
		WithPanel(cpuUsageTimeseries().GridPos(dashboard.GridPos{H: 7, W: 18})).    // TODO: painful, not intuitive
		WithPanel(cpuTemperatureGauge().GridPos(dashboard.GridPos{H: 7, W: 6})).    // TODO: painful, not intuitive
		WithPanel(loadAverageTimeseries().GridPos(dashboard.GridPos{H: 7, W: 18})). // TODO: painful, not intuitive
		// Memory
		WithRow(dashboard.NewRowPanelBuilder("Memory").GridPos(dashboard.GridPos{H: 1, W: 24})). // TODO: painful, not intuitive
		WithPanel(memoryUsageTimeseries().GridPos(dashboard.GridPos{H: 7, W: 18})).
		WithPanel(memoryUsageGauge().GridPos(dashboard.GridPos{H: 7, W: 6})).
		// Disk
		WithRow(dashboard.NewRowPanelBuilder("Disk")).
		WithPanel(diskIOTimeseries().GridPos(dashboard.GridPos{H: 7, W: 12})).
		WithPanel(diskSpaceUsageTable().GridPos(dashboard.GridPos{H: 7, W: 12})).
		// Network
		WithRow(dashboard.NewRowPanelBuilder("Network")).
		WithPanel(networkReceivedTimeseries().GridPos(dashboard.GridPos{H: 7, W: 12})).
		WithPanel(networkTransmittedTimeseries().GridPos(dashboard.GridPos{H: 7, W: 12})).
		// Logs
		WithRow(dashboard.NewRowPanelBuilder("Logs")).
		WithPanel(errorsInSystemLogs().GridPos(dashboard.GridPos{H: 7, W: 24})).
		WithPanel(authLogs().GridPos(dashboard.GridPos{H: 7, W: 24})).
		WithPanel(kernelLogs().GridPos(dashboard.GridPos{H: 7, W: 24})).
		WithPanel(allSystemLogs().GridPos(dashboard.GridPos{H: 7, W: 24}))

	sampleDashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}
	dashboardJson, err := json.MarshalIndent(sampleDashboard, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dashboardJson))
}
