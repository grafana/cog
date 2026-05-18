package defaults

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"reflect"
)

type SomeStruct struct {
    Options map[string]any `json:"options,omitempty"`
    Items []string `json:"items,omitempty"`
    Extra any `json:"extra"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		Options: map[string]any{},
		Items: []string{},
		Extra: map[string]interface{}{},
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `SomeStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *SomeStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "options"
	if fields["options"] != nil {
		if string(fields["options"]) != "null" {
			
			if err := json.Unmarshal(fields["options"], &resource.Options); err != nil {
				errs = append(errs, cog.MakeBuildErrors("options", err)...)
			}
		
		}
		delete(fields, "options")
	
	}
	// Field "items"
	if fields["items"] != nil {
		if string(fields["items"]) != "null" {
			
			if err := json.Unmarshal(fields["items"], &resource.Items); err != nil {
				errs = append(errs, cog.MakeBuildErrors("items", err)...)
			}
		
		}
		delete(fields, "items")
	
	}
	// Field "extra"
	if fields["extra"] != nil {
		if string(fields["extra"]) != "null" {
			if err := json.Unmarshal(fields["extra"], &resource.Extra); err != nil {
				errs = append(errs, cog.MakeBuildErrors("extra", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("extra", errors.New("required field is null"))...)
		
		}
		delete(fields, "extra")
	} else {errs = append(errs, cog.MakeBuildErrors("extra", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	return errs
}


// Equals tests the equality of two `SomeStruct` objects.
func (resource SomeStruct) Equals(other SomeStruct) bool {

		if len(resource.Options) != len(other.Options) {
			return false
		}

		for key1 := range resource.Options {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.Options[key1], other.Options[key1]) {
			return false
		}
		}

		if len(resource.Items) != len(other.Items) {
			return false
		}

		for i1 := range resource.Items {
		if resource.Items[i1] != other.Items[i1] {
			return false
		}
		}
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.Extra, other.Extra) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


