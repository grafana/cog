package intersections

import (
	externalpkg "github.com/grafana/cog/generated/externalpkg"
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type Intersections struct {
	SomeStruct
	externalpkg.AnotherStruct

	FieldString string `json:"fieldString"`
	FieldInteger int32 `json:"fieldInteger"`
}

type SomeStruct struct {
    FieldBool bool `json:"fieldBool"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		FieldBool: true,
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

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


