package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog/plugins"
	dashboard "github.com/grafana/cog/generated/go/dashboardv2"
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
		// TODO: variables
		// TODO: Element() and Elements() should take builders
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
		// TODO: rows?
		Layout(dashboard.NewGridLayoutBuilder().
			Item(dashboard.NewGridLayoutItemBuilder().
				X(0). // TODO: X/Y calculations based on height and width?
				Y(0).
				Height(200).
				Width(200).
				// TODO: proper references
				Element(dashboard.NewElementReferenceBuilder().Name("cpu_usage")),
			),
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
		fmt.Printf("%+v\n", dash)
	}
}
