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
		Overrides([]struct { // TODO: paaaaaain
			Matcher    dashboard.MatcherConfig        `json:"matcher"`
			Properties []dashboard.DynamicConfigValue `json:"properties"`
		}{
			{
				Matcher: dashboard.MatcherConfig{ // TODO: not intuitive
					Id:      "byRegexp",
					Options: "/ io time/",
				},
				Properties: []dashboard.DynamicConfigValue{
					{Id: "unit", Value: "percentunit"},
				},
			},
		}).
		Unit("Bps").
		Targets([]dashboard.Target{
			basicPrometheusQuery(`rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} read"),
			basicPrometheusQuery(`rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} written"),
			basicPrometheusQuery(`rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} IO time"),
		})
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
		Overrides([]struct { // TODO: paaaaaain
			Matcher    dashboard.MatcherConfig        `json:"matcher"`
			Properties []dashboard.DynamicConfigValue `json:"properties"`
		}{
			{
				Matcher: dashboard.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Mounted on",
				},
				Properties: []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 260},
				},
			},
			{
				Matcher: dashboard.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Size",
				},
				Properties: []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 93},
				},
			},
			{
				Matcher: dashboard.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Used",
				},
				Properties: []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 72},
				},
			},
			{
				Matcher: dashboard.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Available",
				},
				Properties: []dashboard.DynamicConfigValue{
					{Id: "custom.width", Value: 88},
				},
			},
			{
				Matcher: dashboard.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Used, %",
				},
				Properties: []dashboard.DynamicConfigValue{
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
			},
		})
}
