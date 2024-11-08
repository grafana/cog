package defaults

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type SomeStruct struct {
    FieldBool bool `json:"fieldBool"`
    FieldString string `json:"fieldString"`
    FieldStringWithConstantValue string `json:"FieldStringWithConstantValue"`
    FieldFloat32 float32 `json:"FieldFloat32"`
    FieldInt32 int32 `json:"FieldInt32"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		FieldBool: true,
		FieldString: "foo",
		FieldStringWithConstantValue: "auto",
		FieldFloat32: 42.42,
		FieldInt32: 42,
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
	// Field "fieldBool"
	if fields["fieldBool"] != nil {
		if string(fields["fieldBool"]) != "null" {
			if err := json.Unmarshal(fields["fieldBool"], &resource.FieldBool); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fieldBool", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("fieldBool", errors.New("required field is null"))...)
		
		}
		delete(fields, "fieldBool")
	
	}
	// Field "fieldString"
	if fields["fieldString"] != nil {
		if string(fields["fieldString"]) != "null" {
			if err := json.Unmarshal(fields["fieldString"], &resource.FieldString); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fieldString", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("fieldString", errors.New("required field is null"))...)
		
		}
		delete(fields, "fieldString")
	
	}
	// Field "FieldStringWithConstantValue"
	if fields["FieldStringWithConstantValue"] != nil {
		if string(fields["FieldStringWithConstantValue"]) != "null" {
			if err := json.Unmarshal(fields["FieldStringWithConstantValue"], &resource.FieldStringWithConstantValue); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldStringWithConstantValue", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldStringWithConstantValue", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldStringWithConstantValue")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldStringWithConstantValue", errors.New("required field is missing from input"))...)
	}
	// Field "FieldFloat32"
	if fields["FieldFloat32"] != nil {
		if string(fields["FieldFloat32"]) != "null" {
			if err := json.Unmarshal(fields["FieldFloat32"], &resource.FieldFloat32); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldFloat32", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldFloat32", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldFloat32")
	
	}
	// Field "FieldInt32"
	if fields["FieldInt32"] != nil {
		if string(fields["FieldInt32"]) != "null" {
			if err := json.Unmarshal(fields["FieldInt32"], &resource.FieldInt32); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldInt32", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldInt32", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldInt32")
	
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
		if resource.FieldBool != other.FieldBool {
			return false
		}
		if resource.FieldString != other.FieldString {
			return false
		}
		if resource.FieldStringWithConstantValue != other.FieldStringWithConstantValue {
			return false
		}
		if resource.FieldFloat32 != other.FieldFloat32 {
			return false
		}
		if resource.FieldInt32 != other.FieldInt32 {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


