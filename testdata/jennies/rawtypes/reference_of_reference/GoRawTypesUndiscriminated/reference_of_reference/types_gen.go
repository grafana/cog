package reference_of_reference

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"fmt"
	"errors"
)

type MyStruct struct {
    Field *OtherStruct `json:"field,omitempty"`
}

// NewMyStruct creates a new MyStruct object.
func NewMyStruct() *MyStruct {
	return &MyStruct{
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
	// Field "field"
	if fields["field"] != nil {
		if string(fields["field"]) != "null" {
			
			resource.Field = &OtherStruct{}
			if err := resource.Field.UnmarshalJSONStrict(fields["field"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("field", err)...)
			}
		
		}
		delete(fields, "field")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("MyStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Validate checks all the validation constraints that may be defined on `MyStruct` fields for violations and returns them.
func (resource MyStruct) Validate() error {
	var errs cog.BuildErrors
		if resource.Field != nil {
		if err := resource.Field.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("field", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type OtherStruct = AnotherStruct

// NewOtherStruct creates a new OtherStruct object.
func NewOtherStruct() *OtherStruct {
	return NewAnotherStruct()
}
type AnotherStruct struct {
    A string `json:"a"`
}

// NewAnotherStruct creates a new AnotherStruct object.
func NewAnotherStruct() *AnotherStruct {
	return &AnotherStruct{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `AnotherStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *AnotherStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "a"
	if fields["a"] != nil {
		if string(fields["a"]) != "null" {
			if err := json.Unmarshal(fields["a"], &resource.A); err != nil {
				errs = append(errs, cog.MakeBuildErrors("a", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("a", errors.New("required field is null"))...)
		
		}
		delete(fields, "a")
	} else {errs = append(errs, cog.MakeBuildErrors("a", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("AnotherStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Validate checks all the validation constraints that may be defined on `AnotherStruct` fields for violations and returns them.
func (resource AnotherStruct) Validate() error {
	return nil
}


