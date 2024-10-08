package variant_panelcfg_full

import (
	variants "github.com/grafana/cog/generated/cog/variants"
)

type Options struct {
	TimeseriesOption string `json:"timeseries_option" yaml:"timeseries_option"`
}

func (resource Options) Equals(other Options) bool {
		if resource.TimeseriesOption != other.TimeseriesOption {
			return false
		}

	return true
}


type FieldConfig struct {
	TimeseriesFieldConfigOption string `json:"timeseries_field_config_option" yaml:"timeseries_field_config_option"`
}

func (resource FieldConfig) Equals(other FieldConfig) bool {
		if resource.TimeseriesFieldConfigOption != other.TimeseriesFieldConfigOption {
			return false
		}

	return true
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

