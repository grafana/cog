package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/dashboard/dashboard"
	"github.com/grafana/cog/generated/dashboard/rowpanel"
	"github.com/grafana/cog/generated/dashboard/timepicker"
	common "github.com/grafana/cog/generated/types/common"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func main() {
	builder := dashboard.New("[TEST] Node Exporter / Raspberry").
		Uid("test-dashboard-raspberry").
		Tags([]string{"generated", "raspberrypi-node-integration"}).
		Refresh("30s").
		Time("now-30m", "now").
		Timezone(common.TimeZoneBrowser).
		Timepicker(
			timepicker.New().
				RefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
				TimeOptions([]string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"}),
		).
		Tooltip(types.DashboardCursorSyncCrosshair).
		Templating(struct { // TODO: uh.
			List []types.VariableModel `json:"list,omitempty"`
		}{
			List: []types.VariableModel{
				{
					Type:        "datasource",
					Name:        "datasource",
					Label:       toPtr("Data Source"),
					Hide:        toPtr(types.VariableHideDontHide),
					SkipUrlSync: toPtr(false),
					Query: &types.StringOrAny{
						String: toPtr("prometheus"),
					},
					Datasource: &types.DataSourceRef{
						Type: toPtr("prometheus"),
						Uid:  toPtr("$datasource"),
					},
					Current: &types.VariableOption{
						Selected: toPtr(true),
						Text:     types.StringOrArrayOfString{String: toPtr("grafanacloud-potatopi-prom")},
						Value:    types.StringOrArrayOfString{String: toPtr("grafanacloud-prom")},
					},
					Multi:   toPtr(false),
					Refresh: toPtr(types.VariableRefreshOnDashboardLoad),
					Sort:    toPtr(types.VariableSortDisabled),
				},
				{
					Type:        "query",
					Name:        "instance",
					Label:       toPtr("Instance"),
					Hide:        toPtr(types.VariableHideDontHide),
					SkipUrlSync: toPtr(false),
					Query: &types.StringOrAny{
						String: toPtr("label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)"),
					},
					Datasource: &types.DataSourceRef{
						Type: toPtr("prometheus"),
						Uid:  toPtr("$datasource"),
					},
					Current: &types.VariableOption{
						Selected: toPtr(false),
						Text:     types.StringOrArrayOfString{String: toPtr("potato")},
						Value:    types.StringOrArrayOfString{String: toPtr("potato")},
					},
					Multi:   toPtr(false),
					Refresh: toPtr(types.VariableRefreshOnTimeRangeChanged),
					Sort:    toPtr(types.VariableSortDisabled),
				},
				{
					Type:        "query",
					Name:        "instance",
					Label:       toPtr("Instance"),
					Hide:        toPtr(types.VariableHideDontHide),
					SkipUrlSync: toPtr(false),
					Query: &types.StringOrAny{
						String: toPtr("label_values(node_uname_info{job=\"integrations/raspberrypi-node\", sysname!=\"Darwin\"}, instance)"),
					},
					Datasource: &types.DataSourceRef{
						Type: toPtr("prometheus"),
						Uid:  toPtr("$datasource"),
					},
					Current: &types.VariableOption{
						Selected: toPtr(false),
						Text:     types.StringOrArrayOfString{String: toPtr("potato")},
						Value:    types.StringOrArrayOfString{String: toPtr("potato")},
					},
					Multi:   toPtr(false),
					Refresh: toPtr(types.VariableRefreshOnTimeRangeChanged),
					Sort:    toPtr(types.VariableSortDisabled),
				},
			},
		}).
		WithRow(rowpanel.New("CPU").GridPos(types.GridPos{H: 1, W: 24})).
		WithPanel(cpuUsageTimeseries().GridPos(types.GridPos{H: 7, W: 18})).    // TODO: painful, not intuitive
		WithPanel(cpuTemperatureGauge().GridPos(types.GridPos{H: 7, W: 6})).    // TODO: painful, not intuitive
		WithPanel(loadAverageTimeseries().GridPos(types.GridPos{H: 7, W: 18})). // TODO: painful, not intuitive
		WithRow(rowpanel.New("Memory").GridPos(types.GridPos{H: 1, W: 24})).    // TODO: painful, not intuitive
		WithPanel(memoryUsageTimeseries().GridPos(types.GridPos{H: 7, W: 18})).
		WithPanel(memoryUsageGauge().GridPos(types.GridPos{H: 7, W: 6})).
		WithRow(rowpanel.New("Disk")).
		WithPanel(diskIOTimeseries().GridPos(types.GridPos{H: 7, W: 12})).
		WithPanel(diskSpaceUsageTable().GridPos(types.GridPos{H: 7, W: 12})).
		WithRow(rowpanel.New("Network")).
		WithPanel(networkReceivedTimeseries().GridPos(types.GridPos{H: 7, W: 12})).
		WithPanel(networkTransmittedTimeseries().GridPos(types.GridPos{H: 7, W: 12})).
		WithRow(rowpanel.New("Logs")).
		WithPanel(errorsInSystemLogs().GridPos(types.GridPos{H: 7, W: 24})).
		WithPanel(authLogs().GridPos(types.GridPos{H: 7, W: 24})).
		WithPanel(kernelLogs().GridPos(types.GridPos{H: 7, W: 24})).
		WithPanel(allSystemLogs().GridPos(types.GridPos{H: 7, W: 24}))

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
