package withdashes

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"reflect"
)

type SomeStruct struct {
    FieldAny any `json:"FieldAny"`
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
	// Field "FieldAny"
	if fields["FieldAny"] != nil {
		if string(fields["FieldAny"]) != "null" {
			if err := json.Unmarshal(fields["FieldAny"], &resource.FieldAny); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldAny", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldAny", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldAny")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldAny", errors.New("required field is missing from input"))...)
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
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


// Refresh rate or disabled.
type RefreshRate = StringOrBool

// NewRefreshRate creates a new RefreshRate object.
func NewRefreshRate() *RefreshRate {
	return NewStringOrBool()
}
type StringOrBool struct {
    String *string `json:"String,omitempty"`
    Bool *bool `json:"Bool,omitempty"`
}

// NewStringOrBool creates a new StringOrBool object.
func NewStringOrBool() *StringOrBool {
	return &StringOrBool{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `StringOrBool` as JSON.
func (resource StringOrBool) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}


	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StringOrBool` from JSON.
func (resource *StringOrBool) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// String
	var String string
	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
		resource.String = nil
	} else {
		resource.String = &String
		return nil
	}

	// Bool
	var Bool bool
	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
		resource.Bool = nil
	} else {
		resource.Bool = &Bool
		return nil
	}

	return errors.Join(errList...)
}


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StringOrBool` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StringOrBool) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// String
	var String string

	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
	} else {
		resource.String = &String
		return nil
	}

	// Bool
	var Bool bool

	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
	} else {
		resource.Bool = &Bool
		return nil
	}


	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StringOrBool", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StringOrBool` objects.
func (resource StringOrBool) Equals(other StringOrBool) bool {
		if resource.String == nil && other.String != nil || resource.String != nil && other.String == nil {
			return false
		}

		if resource.String != nil {
		if *resource.String != *other.String {
			return false
		}
		}
		if resource.Bool == nil && other.Bool != nil || resource.Bool != nil && other.Bool == nil {
			return false
		}

		if resource.Bool != nil {
		if *resource.Bool != *other.Bool {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StringOrBool` fields for violations and returns them.
func (resource StringOrBool) Validate() error {
	return nil
}


