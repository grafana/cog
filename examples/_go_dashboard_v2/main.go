package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/cog/plugins"
	"github.com/grafana/cog/generated/go/dashboard"
	"github.com/grafana/cog/generated/go/timeseries"
)

func dashboardBuilder() []byte {
	builder := dashboard.NewDashboardBuilder("[TEST] Node Exporter / Raspberry").
		//Uid("test-dashboard-raspberry"). // no more dashboard UID? (is it because the schema is just for a dashboard "spec"?)
		Tags([]string{"generated", "raspberrypi-node-integration"}).
		TimeSettings(dashboard.NewTimeSettingsBuilder().
			AutoRefresh("30s").
			AutoRefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
			From("now-30m").
			To("now").
			Timezone("browser"),
		).
		CursorSync(dashboard.DashboardCursorSyncCrosshair).
		Layout(dashboard.NewGridLayoutBuilder().
			WithItem(dashboard.NewGridLayoutItemBuilder().
				X(0).
				Y(0).
				Height(200).
				Width(200).
				Element(dashboard.NewElementReferenceBuilder().Name("somePanel")),
			),
		).
		WithElement("somePanel", dashboard.NewPanelKindBuilder().Spec(
			dashboard.NewPanelBuilder().
				Uid("somePanel").
				Title("Some panel").
				Description("veeery descriptive").
				VizConfig(
					timeseries.NewVizConfigBuilder().
						Unit("s").
						PointSize(5).
						WithOverride(
							dashboard.MatcherConfig{Id: "byName", Options: "Available"},
							[]dashboard.DynamicConfigValue{
								{Id: "custom.width", Value: 88},
							},
						),
				).
				Data(dashboard.NewQueryGroupKindBuilder().Spec(
					dashboard.NewQueryGroupSpecBuilder().
						Queries([]cog.Builder[dashboard.PanelQueryKind]{
							dashboard.NewPanelQueryKindBuilder().Spec(
								dashboard.NewPanelQuerySpecBuilder().
									Query(dashboard.NewDataQueryKindBuilder().
										Kind("prometheus").
										Spec(map[string]string{
											"expr": "query here",
										}),
									).
									Datasource(dashboard.DataSourceRef{
										Type: cog.ToPtr("prometheus"),
										Uid:  cog.ToPtr("some-prometheus-datasource"),
									}).
									RefId("A"), // this field is also present in dataquery schemas
							),
						}).
						WithTransformation(dashboard.NewTransformationKindBuilder().
							Kind("transformation ID. eg: `sortBy`").
							Spec(
								dashboard.NewDataTransformerConfigBuilder().
									Id("what's this ID?").
									Topic(dashboard.DataTransformerConfigTopicSeries).
									Options(map[string]any{
										"fields": map[string]any{},
										"sort": []map[string]any{
											{"field": "Mounted on"},
										},
									}),
							)),
				)),
		))

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
