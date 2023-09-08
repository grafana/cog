package main

import (
	"fmt"

	"github.com/grafana/cog/generated/dashboard/dashboard"
	"github.com/grafana/cog/generated/dashboard/timepicker"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func main() {
	refresh := "1m"

	builder, err := dashboard.New(
		"Some title",
		dashboard.Uid("test-dashboard-codegen"),
		dashboard.Description("Some description"),
		dashboard.Time("now-3h", "now"),
		dashboard.Timepicker(
			timepicker.RefreshIntervals([]string{"30s", "1m", "5m"}),
		),
		dashboard.Style(types.StyleDark),
		dashboard.Timezone("utc"),
		dashboard.Tooltip(types.Crosshair),
		dashboard.Tags([]string{"generated", "from", "cue"}),
		dashboard.Links([]types.DashboardLink{
			{
				Title:       "Some link",
				Url:         "http://google.com",
				AsDropdown:  false,
				TargetBlank: true,
			},
		}),

		dashboard.Refresh(types.StringOrBool{ValString: &refresh}),
	)
	if err != nil {
		panic(err)
	}

	dashboardJson, err := builder.Internal().MarshalIndentJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dashboardJson))
}
