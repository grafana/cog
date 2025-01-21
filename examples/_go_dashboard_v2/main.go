package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/cog/plugins"
	"github.com/grafana/cog/generated/go/common"
	"github.com/grafana/cog/generated/go/dashboard"
	"github.com/grafana/cog/generated/go/loki"
	"github.com/grafana/cog/generated/go/timeseries"
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
		WithElement("somePanel", dashboard.NewPanelBuilder().
			Title("Some panel").
			Description("veeery descriptive").
			Visualization(timeseries.NewVisualizationBuilder().
				Unit("s").
				PointSize(5).
				/*
					// This would be nice
					overrideByName("Available", []dashboard.DynamicConfigValue{
							{Id: "custom.width", Value: 88},
						},
					))
				*/
				Override(
					dashboard.MatcherConfig{Id: "byName", Options: "Available"},
					[]dashboard.DynamicConfigValue{
						{Id: "custom.width", Value: 88},
					},
				),
			).
			Data(dashboard.NewQueryGroupBuilder().
				Target(dashboard.NewTargetBuilder().
					Query(loki.NewQueryBuilder().
						RefId("A").  // TODO: also present on the target builder
						Hide(false). // TODO: also present on the target builder
						Expr("loki expr"),
					).
					Datasource(dashboard.DataSourceRef{
						Type: cog.ToPtr("loki"),
						Uid:  cog.ToPtr("some-loki-datasource"),
					}).
					Hidden(false).
					RefId("A"), // this field is also present in dataquery schemas
				).
				// TODO: simplify this
				Transformation(dashboard.NewTransformationKindBuilder().
					Kind("transformation ID. eg: `sortBy`").
					Spec(dashboard.NewDataTransformerConfigBuilder().
						Id("what's this ID?").
						Topic(common.DataTopicSeries).
						Options(map[string]any{
							"fields": map[string]any{},
							"sort": []map[string]any{
								{"field": "Mounted on"},
							},
						}),
					),
				),
			),
		).
		// TODO: rows?
		Layout(dashboard.NewGridLayoutBuilder().
			Item(dashboard.NewGridLayoutItemBuilder().
				X(0). // TODO: X/Y calculations based on height and width?
				Y(0).
				Height(200).
				Width(200).
				// TODO: proper references
				Element(dashboard.NewElementReferenceBuilder().Name("somePanel")),
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
}
