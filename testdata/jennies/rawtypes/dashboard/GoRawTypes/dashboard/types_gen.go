package dashboard

import (
	cogvariants "github.com/grafana/cog/generated/cog/variants"
)

type Dashboard struct {
	Title string `json:"title"`
	Panels []Panel `json:"panels,omitempty"`
}

type DataSourceRef struct {
	Type *string `json:"type,omitempty"`
	Uid *string `json:"uid,omitempty"`
}

type FieldConfigSource struct {
	Defaults *FieldConfig `json:"defaults,omitempty"`
}

type FieldConfig struct {
	Unit *string `json:"unit,omitempty"`
	Custom any `json:"custom,omitempty"`
}

type Panel struct {
	Title string `json:"title"`
	Type string `json:"type"`
	Datasource *DataSourceRef `json:"datasource,omitempty"`
	Options any `json:"options,omitempty"`
	Targets []cogvariants.Dataquery `json:"targets,omitempty"`
	FieldConfig *FieldConfigSource `json:"fieldConfig,omitempty"`
}

