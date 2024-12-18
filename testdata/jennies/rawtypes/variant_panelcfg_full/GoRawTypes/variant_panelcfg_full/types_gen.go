package variant_panelcfg_full

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type Options struct {
    TimeseriesOption string `json:"timeseries_option"`
}

// NewOptions creates a new Options object.
func NewOptions() *Options {
	return &Options{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Options` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
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


// Equals tests the equality of two `Options` objects.
func (resource Options) Equals(other Options) bool {
		if resource.TimeseriesOption != other.TimeseriesOption {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Options` fields for violations and returns them.
func (resource Options) Validate() error {
	return nil
}


type FieldConfig struct {
    TimeseriesFieldConfigOption string `json:"timeseries_field_config_option"`
}

// NewFieldConfig creates a new FieldConfig object.
func NewFieldConfig() *FieldConfig {
	return &FieldConfig{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `FieldConfig` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
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


// Equals tests the equality of two `FieldConfig` objects.
func (resource FieldConfig) Equals(other FieldConfig) bool {
		if resource.TimeseriesFieldConfigOption != other.TimeseriesFieldConfigOption {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `FieldConfig` fields for violations and returns them.
func (resource FieldConfig) Validate() error {
	return nil
}


