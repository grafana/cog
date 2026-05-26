package enums_as_map_index

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type StringEnum string
const (
	StringEnumA StringEnum = "a"
	StringEnumB StringEnum = "b"
	StringEnumC StringEnum = "c"
)


type StringEnumWithDefault string
const (
	StringEnumWithDefaultA StringEnumWithDefault = "a"
	StringEnumWithDefaultB StringEnumWithDefault = "b"
	StringEnumWithDefaultC StringEnumWithDefault = "c"
)


type SomeStruct struct {
    Data map[StringEnum]string `json:"data"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		Data: map[StringEnum]string{},
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
	// Field "data"
	if fields["data"] != nil {
		if string(fields["data"]) != "null" {
			
			if err := json.Unmarshal(fields["data"], &resource.Data); err != nil {
				errs = append(errs, cog.MakeBuildErrors("data", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("data", errors.New("required field is null"))...)
		
		}
		delete(fields, "data")
	} else {errs = append(errs, cog.MakeBuildErrors("data", errors.New("required field is missing from input"))...)
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

		if len(resource.Data) != len(other.Data) {
			return false
		}

		for key1 := range resource.Data {
		if resource.Data[key1] != other.Data[key1] {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


type SomeStructWithDefaultEnum struct {
    Data map[StringEnumWithDefault]string `json:"data"`
}

// NewSomeStructWithDefaultEnum creates a new SomeStructWithDefaultEnum object.
func NewSomeStructWithDefaultEnum() *SomeStructWithDefaultEnum {
	return &SomeStructWithDefaultEnum{
		Data: map[StringEnumWithDefault]string{},
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `SomeStructWithDefaultEnum` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *SomeStructWithDefaultEnum) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "data"
	if fields["data"] != nil {
		if string(fields["data"]) != "null" {
			
			if err := json.Unmarshal(fields["data"], &resource.Data); err != nil {
				errs = append(errs, cog.MakeBuildErrors("data", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("data", errors.New("required field is null"))...)
		
		}
		delete(fields, "data")
	} else {errs = append(errs, cog.MakeBuildErrors("data", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeStructWithDefaultEnum", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `SomeStructWithDefaultEnum` objects.
func (resource SomeStructWithDefaultEnum) Equals(other SomeStructWithDefaultEnum) bool {

		if len(resource.Data) != len(other.Data) {
			return false
		}

		for key1 := range resource.Data {
		if resource.Data[key1] != other.Data[key1] {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStructWithDefaultEnum` fields for violations and returns them.
func (resource SomeStructWithDefaultEnum) Validate() error {
	return nil
}


