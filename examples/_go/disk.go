package main

import (
	"github.com/grafana/cog/generated/common/tablefooteroptions"
	table "github.com/grafana/cog/generated/table/panel"
	timeseries "github.com/grafana/cog/generated/timeseries/panel"
	common "github.com/grafana/cog/generated/types/common"
	types "github.com/grafana/cog/generated/types/dashboard"
)

func diskIOTimeseries() *timeseries.Builder {
	return defaultTimeseries().
		Title("Disk I/O").
		FillOpacity(0).
		Overrides([]struct { // TODO: paaaaaain
			Matcher    types.MatcherConfig        `json:"matcher"`
			Properties []types.DynamicConfigValue `json:"properties"`
		}{
			{
				Matcher: types.MatcherConfig{ // TODO: not intuitive
					Id:      "byRegexp",
					Options: "/ io time/",
				},
				Properties: []types.DynamicConfigValue{
					{Id: "unit", Value: "percentunit"},
				},
			},
		}).
		Unit("Bps").
		Targets([]types.Target{
			basicPrometheusQuery(`rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} read"),
			basicPrometheusQuery(`rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} written"),
			basicPrometheusQuery(`rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} IO time"),
		})
}

func diskSpaceUsageTable() *table.Builder {
	return table.New().
		Title("Disk Space Usage").
		Align(common.FieldTextAlignmentAuto).
		CellOptions(common.TableCellOptions{
			TableAutoCellOptions: &common.TableAutoCellOptions{
				Type: "auto",
			},
		}).
		Unit("decbytes").
		CellHeight(common.TableCellHeightSm).
		Footer(tablefooteroptions.New().CountRows(false).Reducer([]string{"sum"})).
		Targets([]types.Target{
			tablePrometheusQuery(`max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "A"),
			tablePrometheusQuery(`max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "B"),
		}).
		Transformations([]types.DataTransformerConfig{
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
			Matcher    types.MatcherConfig        `json:"matcher"`
			Properties []types.DynamicConfigValue `json:"properties"`
		}{
			{
				Matcher: types.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Mounted on",
				},
				Properties: []types.DynamicConfigValue{
					{Id: "custom.width", Value: 260},
				},
			},
			{
				Matcher: types.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Size",
				},
				Properties: []types.DynamicConfigValue{
					{Id: "custom.width", Value: 93},
				},
			},
			{
				Matcher: types.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Used",
				},
				Properties: []types.DynamicConfigValue{
					{Id: "custom.width", Value: 72},
				},
			},
			{
				Matcher: types.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Available",
				},
				Properties: []types.DynamicConfigValue{
					{Id: "custom.width", Value: 88},
				},
			},
			{
				Matcher: types.MatcherConfig{ // TODO: not intuitive
					Id:      "byName",
					Options: "Used, %",
				},
				Properties: []types.DynamicConfigValue{
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
