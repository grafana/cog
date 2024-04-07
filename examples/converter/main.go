package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/cog/generated/cog"
	"github.com/grafana/cog/generated/cog/plugins"
	"github.com/grafana/cog/generated/dashboard"
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
		FiscalYearStartMonth(0).
		LiveNow(false).
		Refresh("1m").
		Version(18).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(55).
			Title("Version").
			Description("Blocky [version](https://github.com/0xERR0R/blocky) number").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			RepeatDirection("v").
			WithTransformation(dashboard.DataTransformerConfig{Id: "labelsToFields", Options: map[string]interface{}{}}).
			WithTransformation(dashboard.DataTransformerConfig{Id: "merge", Options: map[string]interface{}{}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(26).
			Title("State").
			Description("current service state").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{ValueMap: cog.ToPtr[dashboard.ValueMap](dashboard.ValueMap{Type: "value", Options: map[string]dashboard.ValueMappingResult{"0": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("down")}, "1": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("up")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "red"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "orange"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(43).
			Title("Blocking").
			Description("Is blocking enabled?").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{ValueMap: cog.ToPtr[dashboard.ValueMap](dashboard.ValueMap{Type: "value", Options: map[string]dashboard.ValueMappingResult{"0": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("off")}, "1": dashboard.ValueMappingResult{Text: cog.ToPtr[string]("on")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "#d44a3a"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "rgba(237, 129, 40, 0.89)"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(57).
			Title("Last list refresh").
			Description("Time since last list refresh").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			Unit("s").
			Decimals(0).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(4).
			Title("Query Count Total").
			Description("Number of all queries. Shows the last value").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			HideTimeOverride(true).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(24).
			Title("Avg response time").
			Description("Average query response time for all query types").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			Unit("ms").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(30).
			Title("Blacklist entries total").
			Description("Number of blacklist entries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(47).
			Title("Cache Hit/Miss ratio").
			Description("Cache Hit/Miss ratio. 100 % means, all queries could be answered from the cache, 0% - all queries must be resolved via external DNS").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			Unit("percentunit").
			Min(0).
			Max(1).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(36).
			Title("Error count").
			Description("Number of occured errors").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			Unit("short").
			Decimals(0).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](1), Color: "orange"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(34).
			Title("Queries blocked").
			Description("Percentage of blocked queries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			Unit("percentunit").
			Decimals(2).
			Min(0).
			Max(1).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(45).
			Title("Cache entries count").
			Description("Number of entries in the cache. Shows the last value").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(28).
			Title("Memory allocated").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			MaxDataPoints(100).
			Unit("bytes").
			Decimals(2).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: cog.ToPtr[string]("N/A")}}})}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(53).
			Title("Prefetch count").
			Description("Amount of performed DNS queries to prefetch cached queries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(49).
			Title("Prefetch domain count").
			Description("Amount of unique domains in the prefetched cache").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(51).
			Title("Prefetch rate per min").
			Description("Amount of prefetch queries per minute").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(58).
			Title("Prefetch Hit ratio").
			Description("How many of cached entries were prefetched automatically").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(3).
			Span(6).
			Unit("percentunit").
			Min(0).
			Max(1).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("timeseries").
			Id(52).
			Title("Request rate per client").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(7).
			Span(24).
			Unit("reqpm").
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode("absolute").
				Steps([]dashboard.Threshold{dashboard.Threshold{Color: "green"},
					dashboard.Threshold{Value: cog.ToPtr[float64](80), Color: "red"}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(2).
			Title("Query by type").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(8).
			Span(8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(8).
			Title("Query per Client").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(8).
			Span(8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(32).
			Title("Blacklist by group").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(8).
			Span(8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(38).
			Title("Response Type").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(8).
			Span(8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(14).
			Title("Response Reasons").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(8).
			Span(8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(12).
			Title("Response status").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(8).
			Span(8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("heatmap").
			Id(22).
			Title("request duration (upstream)").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{Type: cog.ToPtr[string]("prometheus"), Uid: cog.ToPtr[string]("grafanacloud-prom")}).
			Height(9).
			Span(24)).
		WithRow(dashboard.NewRowBuilder("Overview").
			Collapsed(false).
			Id(61)).
		WithRow(dashboard.NewRowBuilder("Activity").
			Collapsed(false).
			Id(60)).
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
