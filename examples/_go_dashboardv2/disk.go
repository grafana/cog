package main

import (
	"github.com/grafana/cog/generated/go/cog"
	"github.com/grafana/cog/generated/go/common"
	dashboard "github.com/grafana/cog/generated/go/dashboardv2beta1"
	"github.com/grafana/cog/generated/go/table"
	"github.com/grafana/cog/generated/go/units"
)

func diskIOTimeseries() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Disk I/O").
		Visualization(
			defaultTimeseries().
				FillOpacity(0).
				Unit(units.BytesPerSecondSI).
				OverrideByRegexp("/ io time/", []dashboard.DynamicConfigValue{
					{Id: "unit", Value: units.PercentUnit},
				}),
		).
		Data(
			dashboard.NewQueryGroupBuilder().
				Targets([]cog.Builder[dashboard.PanelQueryKind]{
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} read")).RefId("A"),
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} written")).RefId("B"),
					dashboard.NewTargetBuilder().Query(basicPrometheusQuery(`rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} IO time")).RefId("C"),
				}),
		)
}

func diskSpaceUsageTable() *dashboard.PanelBuilder {
	return dashboard.NewPanelBuilder().
		Title("Disk Space Usage").
		Visualization(
			table.NewVisualizationBuilder().
				Align(common.FieldTextAlignmentAuto).
				Unit(units.BytesSI).
				CellHeight(common.TableCellHeightSm).
				Footer(common.NewTableFooterOptionsBuilder().CountRows(false).Reducer([]string{"sum"})).
				// Overrides configuration
				OverrideByName("Mounted on", []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 260},
				}).
				OverrideByName("Size", []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 93},
				}).
				OverrideByName("Used", []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 72},
				}).
				OverrideByName("Available", []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 88},
				}).
				OverrideByName("Used, %", []dashboard.DynamicConfigValue{
					{Id: "unit", Value: units.PercentUnit},
					{Id: "custom.cellOptions", Value: struct {
						Mode string `json:"mode"`
						Type string `json:"type"`
					}{
						Mode: "gradient",
						Type: "gauge",
					}},
					{Id: "min", Value: 0},
					{Id: "max", Value: 1},
				}),
		).
		Data(
			dashboard.NewQueryGroupBuilder().
				Targets([]cog.Builder[dashboard.PanelQueryKind]{
					dashboard.NewTargetBuilder().
						RefId("A"). // TODO: unsure which refId should be set (this one, or on the query itself)
						Query(tablePrometheusQuery(`max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`)),
					dashboard.NewTargetBuilder().
						RefId("B").
						Query(tablePrometheusQuery(`max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`)),
				}).
				// TODO: transformations are clunky
				Transformation(dashboard.NewTransformationBuilder().
					Kind("groupBy").
					Id("groupBy"). // TODO: what's this ID?
					Options(map[string]any{
						"fields": map[string]any{
							"Value #A": map[string]any{
								"aggregations": []string{"lastNotNull"},
								"operation":    "aggregate",
							},
							"Value #B": map[string]any{
								"aggregations": []string{"lastNotNull"},
								"operation":    "aggregate",
							},
							"mountpoint": map[string]any{
								"aggregations": []string{},
								"operation":    "groupby",
							},
						},
					}),
				).
				Transformation(dashboard.NewTransformationBuilder().
					Kind("merge").
					Id("merge").
					Options(map[string]any{}),
				).
				Transformation(dashboard.NewTransformationBuilder().
					Kind("calculateField").
					Id("calculateField").
					Options(map[string]any{
						"alias": "Used",
						"binary": map[string]any{
							"left":     "Value #A (lastNotNull)",
							"operator": "-",
							"reducer":  "sum",
							"right":    "Value #B (lastNotNull)",
						},
						"mode": "binary",
						"reduce": map[string]any{
							"reducer": "sum",
						},
					}),
				).
				Transformation(dashboard.NewTransformationBuilder().
					Kind("calculateField").
					Id("calculateField").
					Options(map[string]any{
						"alias": "Used, %",
						"binary": map[string]any{
							"left":     "Used",
							"operator": "/",
							"reducer":  "sum",
							"right":    "Value #A (lastNotNull)",
						},
						"mode": "binary",
						"reduce": map[string]any{
							"reducer": "sum",
						},
					}),
				).
				Transformation(dashboard.NewTransformationBuilder().
					Kind("organize").
					Id("organize").
					Options(map[string]any{
						"excludeByName": map[string]any{},
						"indexByName":   map[string]any{},
						"renameByName": map[string]any{
							"Value #A (lastNotNull)": "Size",
							"Value #B (lastNotNull)": "Available",
							"mountpoint":             "Mounted on",
						},
					}),
				).
				Transformation(dashboard.NewTransformationBuilder().
					Kind("sortBy").
					Id("sortBy").
					Options(map[string]any{
						"fields": map[string]any{},
						"sort": []map[string]any{
							{"field": "Mounted on"},
						},
					}),
				),
		)
}
