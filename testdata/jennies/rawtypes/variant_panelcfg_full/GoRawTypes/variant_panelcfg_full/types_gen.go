package variant_panelcfg_full

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	variants "github.com/grafana/cog/generated/cog/variants"
)

type Options struct {
	TimeseriesOption string `json:"timeseries_option"`
}

func (resource *Options) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "timeseries_option"
	if fields["timeseries_option"] != nil {
		if string(fields["timeseries_option"]) != "null" {
			if err := json.Unmarshal(fields["timeseries_option"], &resource.TimeseriesOption); err != nil {
				errs = append(errs, cog.MakeBuildErrors("timeseries_option", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("timeseries_option", errors.New("required field is null"))...)
		
		}
		delete(fields, "timeseries_option")
	} else {errs = append(errs, cog.MakeBuildErrors("timeseries_option", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Options", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
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

func (resource *FieldConfig) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "timeseries_field_config_option"
	if fields["timeseries_field_config_option"] != nil {
		if string(fields["timeseries_field_config_option"]) != "null" {
			if err := json.Unmarshal(fields["timeseries_field_config_option"], &resource.TimeseriesFieldConfigOption); err != nil {
				errs = append(errs, cog.MakeBuildErrors("timeseries_field_config_option", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("timeseries_field_config_option", errors.New("required field is null"))...)
		
		}
		delete(fields, "timeseries_field_config_option")
	} else {errs = append(errs, cog.MakeBuildErrors("timeseries_field_config_option", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("FieldConfig", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
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
		StrictOptionsUnmarshaler: func (raw []byte) (any, error) {
			options := &Options{}

			if err := options.UnmarshalJSONStrict(raw); err != nil {
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
		StrictFieldConfigUnmarshaler: func (raw []byte) (any, error) {
			fieldConfig := &FieldConfig{}

			if err := fieldConfig.UnmarshalJSONStrict(raw); err != nil {
                return nil, err
            }

			return fieldConfig, nil
		},
	}
}

