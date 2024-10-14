package variant_panelcfg_full

import (
	cog "github.com/grafana/cog/generated/cog"
	variants "github.com/grafana/cog/generated/cog/variants"
	"encoding/json"
)

type Options struct {
	TimeseriesOption string `json:"timeseries_option"`
}

func (resource Options) Equals(other Options) bool {
		if resource.TimeseriesOption != other.TimeseriesOption {
			return false
		}

	return true
}


func (resource Options) Validate() error {
	var errs cog.BuildErrors

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type FieldConfig struct {
	TimeseriesFieldConfigOption string `json:"timeseries_field_config_option"`
}

func (resource FieldConfig) Equals(other FieldConfig) bool {
		if resource.TimeseriesFieldConfigOption != other.TimeseriesFieldConfigOption {
			return false
		}

	return true
}


func (resource FieldConfig) Validate() error {
	var errs cog.BuildErrors

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func VariantConfig() variants.PanelcfgConfig {
	return variants.PanelcfgConfig{
		Identifier: "timeseries",
		OptionsUnmarshaler: func (raw []byte) (any, error) {
			options := &Options{}

			if err := json.Unmarshal(raw, options); err != nil {
				return nil, err
			}

			return options, nil
		},
		FieldConfigUnmarshaler: func (raw []byte) (any, error) {
			fieldConfig := &FieldConfig{}

			if err := json.Unmarshal(raw, fieldConfig); err != nil {
				return nil, err
			}

			return fieldConfig, nil
		},
	}
}

