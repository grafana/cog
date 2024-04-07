package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/cog/generated/cog"
	"github.com/grafana/cog/generated/cog/plugins"
	"github.com/grafana/cog/generated/common"
	"github.com/grafana/cog/generated/dashboard"
	"github.com/grafana/cog/generated/heatmap"
	"github.com/grafana/cog/generated/piechart"
	"github.com/grafana/cog/generated/stat"
	"github.com/grafana/cog/generated/timeseries"
)

func main() {
	plugins.RegisterDefaultPlugins()

	dashboardJSON, err := os.ReadFile("/home/kevin/sandbox/work/cog/examples/converter/dashboard.json")
	if err != nil {
		panic(err)
	}

	dash := &dashboard.Dashboard{}
	if err := json.Unmarshal(dashboardJSON, dash); err != nil {
		panic(err)
	}

	convertedDash := dashboard.DashboardConverter(dash)
	fmt.Println(convertedDash)
}

func foo() {
	dashboard.NewDashboardBuilder("Blocky").
		Id(167).
		Uid("JvOqE4gRk").
		Tags([]string{"blocky"}).
		Editable().
		Tooltip(0).
		Time("now-30d", "now").
		Timepicker(dashboard.NewTimePickerBuilder().
			RefreshIntervals([]string{"5s",
				"10s",
				"30s",
				"1m",
				"5m",
				"15m",
				"30m",
				"1h",
				"2h",
				"1d"})).
		FiscalYearStartMonth(0x0).
		LiveNow(false).
		Refresh("1m").
		Version(0x12).
		WithPanel(stat.NewPanelBuilder().
			Id(0x37).
			Title("Version").
			Description("Blocky [version](https://github.com/0xERR0R/blocky) number").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			RepeatDirection("v").
			WithTransformation(dashboard.DataTransformerConfig{Id: "labelsToFields", Options: map[string]interface{}{}}).
			WithTransformation(dashboard.DataTransformerConfig{Id: "merge", Options: map[string]interface{}{}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("center").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"}).
				Fields("/^version$/")).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x1a).
			Title("State").
			Description("current service state").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{ValueMap: cog.ToPtr[dashboard.ValueMap](dashboard.ValueMap{Type: "value", Options: map[string]dashboard.ValueMappingResult{"0": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("down")}, "1": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("up")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "red"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "orange"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "green"}})).
			GraphMode("none").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x2b).
			Title("Blocking").
			Description("Is blocking enabled?").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{ValueMap: cog.ToPtr[dashboard.ValueMap](dashboard.ValueMap{Type: "value", Options: map[string]dashboard.ValueMappingResult{"0": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("off")}, "1": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("on")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "#d44a3a"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "rgba(237, 129, 40, 0.89)"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "green"}})).
			GraphMode("none").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("value").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x39).
			Title("Last list refresh").
			Description("Time since last list refresh").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			Unit("s").
			Decimals(0).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x4).
			Title("Query Count Total").
			Description("Number of all queries. Shows the last value").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			HideTimeOverride(true).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x18).
			Title("Avg response time").
			Description("Average query response time for all query types").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("ms").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x1e).
			Title("Blacklist entries total").
			Description("Number of blacklist entries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x2f).
			Title("Cache Hit/Miss ratio").
			Description("Cache Hit/Miss ratio. 100 % means, all queries could be answered from the cache, 0% - all queries must be resolved via external DNS").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			Unit("percentunit").
			Min(0).
			Max(1).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"mean"})).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x24).
			Title("Error count").
			Description("Number of occured errors").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("short").
			Decimals(0).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "orange"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x22).
			Title("Queries blocked").
			Description("Percentage of blocked queries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("percentunit").
			Decimals(2).
			Min(0).
			Max(1).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x2d).
			Title("Cache entries count").
			Description("Number of entries in the cache. Shows the last value").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"last"})).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x1c).
			Title("Memory allocated").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("bytes").
			Decimals(2).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("horizontal")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x35).
			Title("Prefetch count").
			Description("Amount of performed DNS queries to prefetch cached queries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x31).
			Title("Prefetch domain count").
			Description("Amount of unique domains in the prefetched cache").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x33).
			Title("Prefetch rate per min").
			Description("Amount of prefetch queries per minute").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(stat.NewPanelBuilder().
			Id(0x3a).
			Title("Prefetch Hit ratio").
			Description("How many of cached entries were prefetched automatically").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x3).
			Span(0x6).
			Unit("percentunit").
			Min(0).
			Max(1).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}})).
			GraphMode("area").
			ColorMode("value").
			JustifyMode("auto").
			TextMode("auto").
			WideLayout(true).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"mean"})).
			ShowPercentChange(false).
			Orientation("auto")).
		WithPanel(timeseries.NewPanelBuilder().
			Id(0x34).
			Title("Request rate per client").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x7).
			Span(0x18).
			Unit("reqpm").
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}})).
			Legend(common.NewVizLegendOptionsBuilder().
				DisplayMode("list").
				Placement("bottom").
				ShowLegend(true)).
			Tooltip(common.NewVizTooltipOptionsBuilder().
				Mode("single").
				Sort("none").
				MaxHeight(600)).
			DrawStyle("bars").
			GradientMode("opacity").
			ThresholdsStyle(common.NewGraphThresholdsStyleConfigBuilder().
				Mode("off")).
			LineWidth(1).
			LineInterpolation("linear").
			FillOpacity(100).
			ShowPoints("never").
			PointSize(5).
			AxisPlacement("auto").
			AxisColorMode("text").
			ScaleDistribution(common.NewScaleDistributionConfigBuilder().
				Type("linear")).
			AxisCenteredZero(false).
			BarAlignment(0).
			Stacking(common.NewStackingConfigBuilder().
				Mode("none").
				Group("A")).
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false)).
			SpanNulls(common.BoolOrFloat64{Bool: cog.ToPtr[bool](true)}).
			AxisBorderShow(false)).
		WithPanel(piechart.NewPanelBuilder().
			Id(0x2).
			Title("Query by type").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0).
			PieType("donut").
			Tooltip(common.NewVizTooltipOptionsBuilder().
				Mode("single").
				Sort("none").
				MaxHeight(600)).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"sum"})).
			Legend(piechart.PieChartLegendOptions{Values: []piechart.PieChartLegendValues{"value", "percent"}, DisplayMode: "table", Placement: "right", ShowLegend: true, Calcs: []string{}}).
			Orientation("").
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false))).
		WithPanel(piechart.NewPanelBuilder().
			Id(0x8).
			Title("Query per Client").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0).
			PieType("donut").
			Tooltip(common.NewVizTooltipOptionsBuilder().
				Mode("single").
				Sort("none").
				MaxHeight(600)).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			Legend(piechart.PieChartLegendOptions{Values: []piechart.PieChartLegendValues{"value", "percent"}, DisplayMode: "table", Placement: "right", ShowLegend: true, Calcs: []string{}}).
			Orientation("").
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false))).
		WithPanel(piechart.NewPanelBuilder().
			Id(0x20).
			Title("Blacklist by group").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0).
			PieType("donut").
			Tooltip(common.NewVizTooltipOptionsBuilder().
				Mode("single").
				Sort("none").
				MaxHeight(600)).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			Legend(piechart.PieChartLegendOptions{Values: []piechart.PieChartLegendValues{"value"}, DisplayMode: "table", Placement: "right", ShowLegend: true, Calcs: []string{}}).
			Orientation("").
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false))).
		WithPanel(piechart.NewPanelBuilder().
			Id(0x26).
			Title("Response Type").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0).
			PieType("donut").
			Tooltip(common.NewVizTooltipOptionsBuilder().
				Mode("single").
				Sort("none").
				MaxHeight(600)).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"sum"})).
			Legend(piechart.PieChartLegendOptions{Values: []piechart.PieChartLegendValues{"value", "percent"}, DisplayMode: "table", Placement: "right", ShowLegend: true, Calcs: []string{}}).
			Orientation("").
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false))).
		WithPanel(piechart.NewPanelBuilder().
			Id(0xe).
			Title("Response Reasons").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0).
			PieType("donut").
			Tooltip(common.NewVizTooltipOptionsBuilder().
				Mode("single").
				Sort("none").
				MaxHeight(600)).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			Legend(piechart.PieChartLegendOptions{Values: []piechart.PieChartLegendValues{"value", "percent"}, DisplayMode: "table", Placement: "right", ShowLegend: true, Calcs: []string{}}).
			Orientation("").
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false))).
		WithPanel(piechart.NewPanelBuilder().
			Id(0xc).
			Title("Response status").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0).
			PieType("donut").
			Tooltip(common.NewVizTooltipOptionsBuilder().
				Mode("single").
				Sort("none").
				MaxHeight(600)).
			ReduceOptions(common.NewReduceDataOptionsBuilder().
				Values(false).
				Calcs([]string{"lastNotNull"})).
			Legend(piechart.PieChartLegendOptions{Values: []piechart.PieChartLegendValues{"value", "percent"}, DisplayMode: "table", Placement: "right", ShowLegend: true, Calcs: []string{}}).
			Orientation("").
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false))).
		WithPanel(heatmap.NewPanelBuilder().
			Id(0x16).
			Title("request duration (upstream)").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(0x9).
			Span(0x18).
			Calculate(false).
			Calculation(common.NewHeatmapCalculationOptionsBuilder().
				YBuckets(common.NewHeatmapCalculationBucketConfigBuilder().
					Scale(common.NewScaleDistributionConfigBuilder().
						Type("linear")))).
			Color(heatmap.HeatmapColorOptions{Mode: cog.ToPtr[heatmap.HeatmapColorMode]("scheme"), Scheme: "RdBu", Fill: "#FADE2A", Scale: cog.ToPtr[heatmap.HeatmapColorScale]("exponential"), Exponent: 0.4000000059604645, Steps: 128, Reverse: false}).
			FilterValues(heatmap.FilterValueRange{Le: cog.ToPtr[float32](9.999999717180685e-10)}).
			RowsFrame(heatmap.RowsHeatmapOptions{Layout: cog.ToPtr[common.HeatmapCellLayout]("auto")}).
			ShowValue("never").
			CellGap(0x2).
			CellValues(heatmap.CellValues{}).
			YAxis(heatmap.YAxisConfig{Unit: cog.ToPtr[string]("ms"), Reverse: cog.ToPtr[bool](false), AxisPlacement: cog.ToPtr[common.AxisPlacement]("left")}).
			Legend(heatmap.HeatmapLegend{Show: true}).
			Tooltip(heatmap.HeatmapTooltip{Mode: "single", MaxHeight: cog.ToPtr[float64](600), YHistogram: cog.ToPtr[bool](false), ShowColorScale: cog.ToPtr[bool](false)}).
			ExemplarsColor(heatmap.ExemplarConfig{Color: "rgba(255,0,255,0.7)"}).
			ScaleDistribution(common.NewScaleDistributionConfigBuilder().
				Type("linear")).
			HideFrom(common.NewHideSeriesConfigBuilder().
				Tooltip(false).
				Legend(false).
				Viz(false))).
		WithRow(dashboard.NewRowBuilder("Overview").
			Collapsed(false).
			Id(0x3d)).
		WithRow(dashboard.NewRowBuilder("Activity").
			Collapsed(false).
			Id(0x3c)).
		WithVariable(dashboard.NewQueryVariableBuilder("job").
			Label("job").
			Hide(0).
			Query(dashboard.StringOrAny{Any: cog.ToPtr[interface{}](map[string]interface{}{"qryType": 1, "query": "label_values(blocky_blocking_enabled,job)", "refId": "PrometheusVariableQueryEditor-VariableQuery"})}).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Current(dashboard.VariableOption{Selected: cog.ToPtr[bool](false), Text: dashboard.StringOrArrayOfString{String: cog.ToPtr[string]("All")}, Value: dashboard.StringOrArrayOfString{String: cog.ToPtr[string]("$__all")}}).
			Multi(false).
			Refresh(1).
			Sort(1).
			IncludeAll(true)).
		Annotations(dashboard.NewAnnotationContainerBuilder().
			List([]cog.Builder[dashboard.AnnotationQuery]{dashboard.NewAnnotationQueryBuilder().
				Name("Annotations & Alerts").
				Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("datasource"), Uid: cog.ToPtr[string]("grafana")}).
				Enable(true).
				Hide(true).
				IconColor("rgba(0, 211, 255, 1)").
				Target(dashboard.NewAnnotationTargetBuilder().
					Limit(100).
					MatchAny(false).
					Type("dashboard")).
				Type("dashboard").
				BuiltIn(1),
				dashboard.NewAnnotationQueryBuilder().
					Name("Reboot").
					Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
					Enable(true).
					IconColor("yellow")})).
		Link(dashboard.NewDashboardLinkBuilder("blocky @ GitHub").
			Type("link").
			Icon("external link").
			Tooltip("open GitHub repo").
			Url("https://github.com/0xERR0R/blocky").
			AsDropdown(false).
			TargetBlank(false).
			IncludeVars(false).
			KeepTime(false))

}
