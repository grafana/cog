package variant_panelcfg_full

import (
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


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource Options) Validate() error {
	return nil
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


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource FieldConfig) Validate() error {
	return nil
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

