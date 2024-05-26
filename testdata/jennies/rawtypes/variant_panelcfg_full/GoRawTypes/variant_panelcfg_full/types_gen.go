package variant_panelcfg_full

import (
	cogvariants "github.com/grafana/cog/generated/cog/variants"
)

type Options struct {
	TimeseriesOption string `json:"timeseries_option"`
}

type FieldConfig struct {
	TimeseriesFieldConfigOption string `json:"timeseries_field_config_option"`
}

func VariantConfig() cogvariants.PanelcfgConfig {
	return cogvariants.PanelcfgConfig{
		Identifier: "timeseries",
		OptionsUnmarshaler: func (raw []byte) (any, error) {
			options := Options{}

			if err := json.Unmarshal(raw, &options); err != nil {
				return nil, err
			}

			return options, nil
		},
		FieldConfigUnmarshaler: func (raw []byte) (any, error) {
			fieldConfig := FieldConfig{}

			if err := json.Unmarshal(raw, &fieldConfig); err != nil {
				return nil, err
			}

			return fieldConfig, nil
		},
	}
}

