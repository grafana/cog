package main

import (
	"github.com/grafana/cog/generated/common"
	"github.com/grafana/cog/generated/dashboard"
	"github.com/grafana/cog/generated/table"
	"github.com/grafana/cog/generated/timeseries"
)

func diskIOTimeseries() *timeseries.PanelBuilder {
	return defaultTimeseries().
		Title("Disk I/O").
		FillOpacity(0).
		Unit("Bps").
		Targets([]dashboard.Target{
			basicPrometheusQuery(`rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} read"),
			basicPrometheusQuery(`rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} written"),
			basicPrometheusQuery(`rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} IO time"),
		}).
		// Overrides configuration
		WithOverride(
			// TODO: not very intuitive
			// we could have "factory" functions:
			// - dashboard.OverrideByName("Mounted on")
			// - dashboard.OverrideByRegexp("/ regex /")
			// - ...
			// Also: knowing what to set in the Value field is far from obvious
			dashboard.MatcherConfig{
				Id:      "byRegexp", // TODO: we don't have constants for these?
				Options: "/ io time/",
			},
			[]dashboard.DynamicConfigValue{
				{Id: "unit", Value: "percentunit"},
			},
		)
}

func diskSpaceUsageTable() *table.PanelBuilder {
	return table.NewPanelBuilder().
		Title("Disk Space Usage").
		Align(common.FieldTextAlignmentAuto).
		CellOptions(common.TableCellOptions{
			TableAutoCellOptions: &common.TableAutoCellOptions{
				Type: "auto",
			},
		}).
		Unit("decbytes").
		CellHeight(common.TableCellHeightSm).
		Footer(common.NewTableFooterOptionsBuilder().CountRows(false).Reducer([]string{"sum"})).
		Targets([]dashboard.Target{
			tablePrometheusQuery(`max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "A"),
			tablePrometheusQuery(`max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "B"),
		}).
		Transformations([]dashboard.DataTransformerConfig{
			{
				Id: "groupBy",
				Options: map[string]any{
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
				},
			},
			{
				Id:      "merge",
				Options: map[string]any{},
			},
			{
				Id: "calculateField",
				Options: map[string]any{
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
				},
			},
			{
				Id: "calculateField",
				Options: map[string]any{
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
				},
			},
			{
				Id: "organize",
				Options: map[string]any{
					"excludeByName": map[string]any{},
					"indexByName":   map[string]any{},
					"renameByName": map[string]any{
						"Value #A (lastNotNull)": "Size",
						"Value #B (lastNotNull)": "Available",
						"mountpoint":             "Mounted on",
					},
				},
			}, {
				Id: "sortBy",
				Options: map[string]any{
					"fields": map[string]any{},
					"sort": []map[string]any{
						{"field": "Mounted on"},
					},
				},
			},
		}).
		// Overrides configuration
		WithOverride(
			dashboard.MatcherConfig{Id: "byName", Options: "Mounted on"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.width", Value: 260},
			},
		).
		WithOverride(
			dashboard.MatcherConfig{Id: "byName", Options: "Size"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.width", Value: 93},
			},
		).
		WithOverride(
			dashboard.MatcherConfig{Id: "byName", Options: "Used"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.width", Value: 72},
			},
		).
		WithOverride(
			dashboard.MatcherConfig{Id: "byName", Options: "Available"},
			[]dashboard.DynamicConfigValue{
				{Id: "custom.width", Value: 88},
			},
		).
		WithOverride(
			dashboard.MatcherConfig{Id: "byName", Options: "Used, %"},
			[]dashboard.DynamicConfigValue{
				{Id: "unit", Value: "percentunit"},
				{Id: "custom.cellOptions", Value: struct {
					Mode string `json:"mode"`
					Type string `json:"type"`
				}{
					Mode: "gradient",
					Type: "gauge",
				}},
				{Id: "min", Value: 0},
				{Id: "max", Value: 1},
			},
		)
}
