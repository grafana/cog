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
		Tooltip(dashboard.DashboardCursorSyncOff).
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
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x37).
			Title("Version").
			Description("Blocky [version](https://github.com/0xERR0R/blocky) number").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			RepeatDirection(dashboard.PanelRepeatDirectionV).
			WithTransformation(dashboard.DataTransformerConfig{
				Id:      "labelsToFields",
				Filter:  nil,
				Topic:   nil,
				Options: map[string]interface{}{},
			}).
			WithTransformation(dashboard.DataTransformerConfig{
				Id:      "merge",
				Filter:  nil,
				Topic:   nil,
				Options: map[string]interface{}{},
			}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](80),
						Color: "red",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x1a).
			Title("State").
			Description("current service state").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				ValueMap: cog.ToPtr[dashboard.ValueMap](dashboard.ValueMap{Type: "value", Options: map[string]dashboard.ValueMappingResult{"0": dashboard.ValueMappingResult{Text: (*string)(0xc0001a80f0), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}, "1": dashboard.ValueMappingResult{Text: (*string)(0xc0001a8110), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}}}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "red",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](1),
						Color: "orange",
					},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](1),
						Color: "green",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x2b).
			Title("Blocking").
			Description("Is blocking enabled?").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				ValueMap: cog.ToPtr[dashboard.ValueMap](dashboard.ValueMap{Type: "value", Options: map[string]dashboard.ValueMappingResult{"0": dashboard.ValueMappingResult{Text: (*string)(0xc0001a8660), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}, "1": dashboard.ValueMappingResult{Text: (*string)(0xc0001a8680), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}}}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "#d44a3a",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](1),
						Color: "rgba(237, 129, 40, 0.89)",
					},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](1),
						Color: "green",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x39).
			Title("Last list refresh").
			Description("Time since last list refresh").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			Unit("s").
			Decimals(0).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x4).
			Title("Query Count Total").
			Description("Number of all queries. Shows the last value").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			HideTimeOverride(true).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: (*string)(0xc0001a9050), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}}}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x18).
			Title("Avg response time").
			Description("Average query response time for all query types").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("ms").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: (*string)(0xc0001a95c0), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}}}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](80),
						Color: "red",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x1e).
			Title("Blacklist entries total").
			Description("Number of blacklist entries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("none").
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: (*string)(0xc0001a9b20), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}}}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x2f).
			Title("Cache Hit/Miss ratio").
			Description("Cache Hit/Miss ratio. 100 % means, all queries could be answered from the cache, 0% - all queries must be resolved via external DNS").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			Unit("percentunit").
			Min(0).
			Max(1).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x24).
			Title("Error count").
			Description("Number of occured errors").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("short").
			Decimals(0).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: (*string)(0xc0001c85a0), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}}}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](1),
						Color: "orange",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x22).
			Title("Queries blocked").
			Description("Percentage of blocked queries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("percentunit").
			Decimals(2).
			Min(0).
			Max(1).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{Type: "special", Options: dashboard.DashboardSpecialValueMapOptions{Match: "null", Result: dashboard.ValueMappingResult{Text: (*string)(0xc0001c8b10), Color: (*string)(nil), Icon: (*string)(nil), Index: (*int32)(nil)}}}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](80),
						Color: "red",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x2d).
			Title("Cache entries count").
			Description("Number of entries in the cache. Shows the last value").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x1c).
			Title("Memory allocated").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			MaxDataPoints(100).
			Unit("bytes").
			Decimals(2).
			Mappings([]dashboard.ValueMapping{dashboard.ValueMapOrRangeMapOrRegexMapOrSpecialValueMap{
				SpecialValueMap: cog.ToPtr[dashboard.SpecialValueMap](dashboard.SpecialValueMap{
					Type: "special",
					Options: dashboard.DashboardSpecialValueMapOptions{
						Match: "null",
						Result: dashboard.ValueMappingResult{
							Text:  (*string)(0xc0001c9490),
							Color: (*string)(nil),
							Icon:  (*string)(nil),
							Index: (*int32)(nil),
						},
					},
				}),
			}}).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x35).
			Title("Prefetch count").
			Description("Amount of performed DNS queries to prefetch cached queries").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x31).
			Title("Prefetch domain count").
			Description("Amount of unique domains in the prefetched cache").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x33).
			Title("Prefetch rate per min").
			Description("Amount of prefetch queries per minute").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](80),
						Color: "red",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("stat").
			Id(0x3a).
			Title("Prefetch Hit ratio").
			Description("How many of cached entries were prefetched automatically").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Uid: cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x3).
			Span(0x6).
			Unit("percentunit").
			Min(0).
			Max(1).
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("timeseries").
			Id(0x34).
			Title("Request rate per client").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x7).
			Span(0x18).
			Unit("reqpm").
			Thresholds(dashboard.NewThresholdsConfigBuilder().
				Mode(dashboard.ThresholdsModeAbsolute).
				Steps([]dashboard.Threshold{dashboard.Threshold{
					Color: "green",
				},
					dashboard.Threshold{
						Value: cog.ToPtr[float64](80),
						Color: "red",
					}}))).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(0x2).
			Title("Query by type").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(0x8).
			Title("Query per Client").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(0x20).
			Title("Blacklist by group").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(0x26).
			Title("Response Type").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(0xe).
			Title("Response Reasons").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("piechart").
			Id(0xc).
			Title("Response status").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x8).
			Span(0x8).
			MaxDataPoints(3).
			Unit("short").
			Decimals(0)).
		WithPanel(dashboard.NewPanelBuilder().
			Type("heatmap").
			Id(0x16).
			Title("request duration (upstream)").
			Transparent(true).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Height(0x9).
			Span(0x18)).
		WithRow(dashboard.NewRowBuilder("Overview").
			Collapsed(false).
			Id(0x3d)).
		WithRow(dashboard.NewRowBuilder("Activity").
			Collapsed(false).
			Id(0x3c)).
		WithVariable(dashboard.NewQueryVariableBuilder("job").
			Label("job").
			Hide(dashboard.VariableHideDontHide).
			Query(dashboard.StringOrAny{
				Any: (*interface{})(0xc00020be40),
			}).
			Datasource(dashboard.DataSourceRef{
				Type: cog.ToPtr[string]("prometheus"),
				Uid:  cog.ToPtr[string]("grafanacloud-prom"),
			}).
			Current(dashboard.VariableOption{
				Selected: cog.ToPtr[bool](false),
				Text: dashboard.StringOrArrayOfString{
					String: cog.ToPtr[string]("All"),
				},
				Value: dashboard.StringOrArrayOfString{
					String: cog.ToPtr[string]("$__all"),
				},
			}).
			Multi(false).
			Refresh(dashboard.VariableRefreshOnDashboardLoad).
			Sort(dashboard.VariableSortAlphabeticalAsc).
			IncludeAll(true)).
		Annotations(dashboard.NewAnnotationContainerBuilder().
			List([]cog.Builder[dashboard.AnnotationQuery]{dashboard.NewAnnotationQueryBuilder().
				Name("Annotations & Alerts").
				Datasource(dashboard.DataSourceRef{
					Type: cog.ToPtr[string]("datasource"),
					Uid:  cog.ToPtr[string]("grafana"),
				}).
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
					Datasource(dashboard.DataSourceRef{
						Type: cog.ToPtr[string]("prometheus"),
						Uid:  cog.ToPtr[string]("grafanacloud-prom"),
					}).
					Enable(true).
					IconColor("yellow")})).
		Link(dashboard.NewDashboardLinkBuilder("blocky @ GitHub").
			Type(dashboard.DashboardLinkTypeLink).
			Icon("external link").
			Tooltip("open GitHub repo").
			Url("https://github.com/0xERR0R/blocky").
			AsDropdown(false).
			TargetBlank(false).
			IncludeVars(false).
			KeepTime(false))

}
