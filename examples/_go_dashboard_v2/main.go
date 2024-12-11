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
		TimeSettings(dashboard.NewTimeSettingsBuilder().
			AutoRefresh("30s").
			AutoRefreshIntervals([]string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}).
			From("now-30m").
			To("now").
			Timezone("browser"),
		).
		CursorSync(dashboard.DashboardCursorSyncCrosshair).
		WithElement("somePanel", dashboard.NewPanelBuilder().
			Uid("somePanelUid"). // TODO: should this be equal to the element ref, or unrelated?
			Title("Some panel").
			Description("veeery descriptive").
			// TODO: better method names. Maybe Vizualization(timeseries.NewVizualization())
			VizConfig(timeseries.NewVizConfigBuilder().
				Unit("s").
				PointSize(5).
				WithOverride(
					dashboard.MatcherConfig{Id: "byName", Options: "Available"},
					[]dashboard.DynamicConfigValue{
						{Id: "custom.width", Value: 88},
					},
				),
			).
			Data(dashboard.NewQueryGroupBuilder().
				// TODO: WithQuery() followed by Query() is a bit repetitive/confusing. Better names needed.
				// Maybe WithTarget(
				//   Query(),
				//   Datasource(),
				//   RefId(),
				// )
				WithQuery(dashboard.NewPanelQueryBuilder().
					Query(loki.NewQueryBuilder().
						Expr("loki expr"),
					).
					Datasource(dashboard.DataSourceRef{
						Type: cog.ToPtr("loki"),
						Uid:  cog.ToPtr("some-loki-datasource"),
					}).
					RefId("A"), // this field is also present in dataquery schemas
				).
				// TODO: simplify this
				WithTransformation(dashboard.NewTransformationKindBuilder().
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
		Layout(dashboard.NewGridLayoutBuilder().
			WithItem(dashboard.NewGridLayoutItemBuilder().
				X(0).
				Y(0).
				Height(200).
				Width(200).
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
