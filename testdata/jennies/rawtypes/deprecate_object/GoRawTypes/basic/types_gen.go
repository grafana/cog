package basic

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

// Deprecated: This object is deprecated, use NewStruct instead.
type SomeStruct struct {
    FieldString string `json:"FieldString"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
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
	// Field "FieldString"
	if fields["FieldString"] != nil {
		if string(fields["FieldString"]) != "null" {
			if err := json.Unmarshal(fields["FieldString"], &resource.FieldString); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldString", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldString", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldString")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldString", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `SomeStruct` objects.
func (resource SomeStruct) Equals(other SomeStruct) bool {
		if resource.FieldString != other.FieldString {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


