package constant_reference_as_default

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

const ConstantRefString = "AString"

type MyStruct struct {
    AString string `json:"aString"`
    OptString *string `json:"optString,omitempty"`
}

// NewMyStruct creates a new MyStruct object.
func NewMyStruct() *MyStruct {
	return &MyStruct{
		AString: ConstantRefString,
		OptString: (func (input string) *string { return &input })(ConstantRefString),
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `MyStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *MyStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "aString"
	if fields["aString"] != nil {
		if string(fields["aString"]) != "null" {
			if err := json.Unmarshal(fields["aString"], &resource.AString); err != nil {
				errs = append(errs, cog.MakeBuildErrors("aString", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("aString", errors.New("required field is null"))...)
		
		}
		delete(fields, "aString")
	} else {errs = append(errs, cog.MakeBuildErrors("aString", errors.New("required field is missing from input"))...)
	}
	// Field "optString"
	if fields["optString"] != nil {
		if string(fields["optString"]) != "null" {
			if err := json.Unmarshal(fields["optString"], &resource.OptString); err != nil {
				errs = append(errs, cog.MakeBuildErrors("optString", err)...)
			}
		
		}
		delete(fields, "optString")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("MyStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `MyStruct` objects.
func (resource MyStruct) Equals(other MyStruct) bool {
        if resource.AString != other.AString {
            return false
        }
		if resource.OptString == nil && other.OptString != nil || resource.OptString != nil && other.OptString == nil {
			return false
		}

		if resource.OptString != nil {
        if *resource.OptString != *other.OptString {
            return false
        }
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `MyStruct` fields for violations and returns them.
func (resource MyStruct) Validate() error {
	return nil
}
