package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/cog/plugins"
	"github.com/grafana/cog/generated/go/dashboard"
)

func dashboardBuilder() []byte {
	builder := dashboard.NewDashboardBuilder("[TEST] Node Exporter / Raspberry").
		//Uid("test-dashboard-raspberry").
		Tags([]string{"generated", "raspberrypi-node-integration"}).
		TimeSettings(dashboard.NewTimeSettingsBuilder().
			AutoRefresh("30s").
			AutoRefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
			From("now-30m").
			To("now").
			Timezone("browser"),
		).
		CursorSync(dashboard.DashboardCursorSyncCrosshair).
		Layout(dashboard.NewGridLayoutKindBuilder().
			Spec(
				dashboard.NewGridLayoutSpecBuilder().
					Items([]cog.Builder[dashboard.GridLayoutItemKind]{
						dashboard.NewGridLayoutItemKindBuilder().Spec(
							dashboard.NewGridLayoutItemSpecBuilder().
								X(0).
								Y(0).
								Height(200).
								Width(200).
								Element(dashboard.NewElementReferenceBuilder().Name("somePanel")),
						),
					}),
			),
		).
		Elements(map[string]dashboard.PanelKind{
			"somePanel": {
				Kind: "Panel",
				Spec: dashboard.Panel{
					Uid:         "somePanel",
					Title:       "Some panel",
					Description: "veeery descriptive",
					//Data:        dashboard.QueryGroupKind{},
					VizConfig: dashboard.VizConfigKind{
						Kind: "", // plugin ID
						Spec: dashboard.VizConfigSpec{
							// PluginVersion: "", // TODO?
							Options: nil, // plugin dependent
							FieldConfig: dashboard.FieldConfigSource{
								Defaults: dashboard.FieldConfig{
									DisplayName:       nil,
									DisplayNameFromDS: nil,
									Description:       nil,
									Path:              nil,
									Writeable:         nil,
									Filterable:        nil,
									Unit:              nil,
									Decimals:          nil,
									Min:               nil,
									Max:               nil,
									Mappings:          nil,
									Thresholds:        nil,
									Color:             nil,
									Links:             nil,
									NoValue:           nil,
									Custom:            nil, // plugin dependent
								},
								Overrides: nil,
							},
						},
					},
				},
			},
		})
	// "Data Source" variable
	/*
			WithVariable(dashboard.NewDatasourceVariableBuilder("datasource").
				Label("Data Source").
				Hide(dashboard.VariableHideDontHide).
				Type("prometheus").
				Current(dashboard.VariableOption{
					Selected: toPtr(true),
					Text:     dashboard.StringOrArrayOfString{String: toPtr("grafanacloud-potatopi-prom")},
					Value:    dashboard.StringOrArrayOfString{String: toPtr("grafanacloud-prom")},
				}),
			).
			// "Instance" variable
			WithVariable(dashboard.NewQueryVariableBuilder("instance").
				Label("Instance").
				Hide(dashboard.VariableHideDontHide).
				Refresh(dashboard.VariableRefreshOnTimeRangeChanged).
				Query(dashboard.StringOrMap{
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
	*/

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
