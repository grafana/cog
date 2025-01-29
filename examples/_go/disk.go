package main

import (
	"github.com/grafana/cog/generated/go/common"
	"github.com/grafana/cog/generated/go/dashboard"
	"github.com/grafana/cog/generated/go/table"
	"github.com/grafana/cog/generated/go/timeseries"
)

func diskIOTimeseries() *timeseries.PanelBuilder {
	return defaultTimeseries().
		Title("Disk I/O").
		FillOpacity(0).
		Unit("Bps").
		WithTarget(
			basicPrometheusQuery(`rate(node_disk_read_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} read"),
		).
		WithTarget(
			basicPrometheusQuery(`rate(node_disk_written_bytes_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} written"),
		).
		WithTarget(
			basicPrometheusQuery(`rate(node_disk_io_time_seconds_total{job="integrations/raspberrypi-node", instance="$instance", device!=""}[$__rate_interval])`, "{{device}} IO time"),
		).
		// Overrides configuration
		OverrideByRegexp("/ io time/", []dashboard.DynamicConfigValue{
			{Id: "unit", Value: "percentunit"},
		})
}

func diskSpaceUsageTable() *table.PanelBuilder {
	return table.NewPanelBuilder().
		Title("Disk Space Usage").
		Align(common.FieldTextAlignmentAuto).
		Unit("decbytes").
		CellHeight(common.TableCellHeightSm).
		Footer(common.NewTableFooterOptionsBuilder().CountRows(false).Reducer([]string{"sum"})).
		WithTarget(
			tablePrometheusQuery(`max by (mountpoint) (node_filesystem_size_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "A"),
		).
		WithTarget(
			tablePrometheusQuery(`max by (mountpoint) (node_filesystem_avail_bytes{job="integrations/raspberrypi-node", instance="$instance", fstype!=""})`, "B"),
		).
		// Transformations
		WithTransformation(dashboard.DataTransformerConfig{
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
		}).
		WithTransformation(dashboard.DataTransformerConfig{
			Id:      "merge",
			Options: map[string]any{},
		}).
		WithTransformation(dashboard.DataTransformerConfig{
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
		}).
		WithTransformation(dashboard.DataTransformerConfig{
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
		}).
		WithTransformation(dashboard.DataTransformerConfig{
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
		}).
		WithTransformation(dashboard.DataTransformerConfig{
			Id: "sortBy",
			Options: map[string]any{
				"fields": map[string]any{},
				"sort": []map[string]any{
					{"field": "Mounted on"},
				},
			},
		}).

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
		})
}
