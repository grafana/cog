package variant_panelcfg_only_options

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type Options struct {
    Content string `json:"content"`
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


// Equals tests the equality of two `Options` objects.
func (resource Options) Equals(other Options) bool {
		if resource.Content != other.Content {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Options` fields for violations and returns them.
func (resource Options) Validate() error {
	return nil
}


