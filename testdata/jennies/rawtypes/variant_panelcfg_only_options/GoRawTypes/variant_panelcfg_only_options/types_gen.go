package variant_panelcfg_only_options

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	variants "github.com/grafana/cog/generated/cog/variants"
)

type Options struct {
	Content string `json:"content"`
}

func (resource *Options) StrictUnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "content"
	if fields["content"] != nil {
		if string(fields["content"]) != "null" {
			if err := json.Unmarshal(fields["content"], &resource.Content); err != nil {
				errs = append(errs, cog.MakeBuildErrors("content", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("content", errors.New("required field is null"))...)
		
		}
		delete(fields, "content")
	} else {errs = append(errs, cog.MakeBuildErrors("content", errors.New("required field is missing from input"))...)
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
		if resource.Content != other.Content {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource Options) Validate() error {
	return nil
}


func VariantConfig() variants.PanelcfgConfig {
	return variants.PanelcfgConfig{
		Identifier: "text",
		OptionsUnmarshaler: func (raw []byte) (any, error) {
			options := &Options{}

			if err := json.Unmarshal(raw, options); err != nil {
				return nil, err
			}

			return options, nil
		},
		StrictOptionsUnmarshaler: func (raw []byte) (any, error) {
			options := &Options{}

			if err := options.StrictUnmarshalJSON(raw); err != nil {
                return nil, err
            }

			return options, nil
		},
	}
}

