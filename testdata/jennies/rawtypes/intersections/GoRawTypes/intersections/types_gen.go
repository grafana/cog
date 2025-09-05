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


// Base properties for all metrics
type Common struct {
    // The metric name
    Name string `json:"name"`
    // The metric type
    Type CommonType `json:"type"`
    // The type of data the metric contains
    Contains CommonContains `json:"contains"`
}

// NewCommon creates a new Common object.
func NewCommon() *Common {
	return &Common{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Common` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Common) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "name"
	if fields["name"] != nil {
		if string(fields["name"]) != "null" {
			if err := json.Unmarshal(fields["name"], &resource.Name); err != nil {
				errs = append(errs, cog.MakeBuildErrors("name", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("name", errors.New("required field is null"))...)
		
		}
		delete(fields, "name")
	} else {errs = append(errs, cog.MakeBuildErrors("name", errors.New("required field is missing from input"))...)
	}
	// Field "type"
	if fields["type"] != nil {
		if string(fields["type"]) != "null" {
			if err := json.Unmarshal(fields["type"], &resource.Type); err != nil {
				errs = append(errs, cog.MakeBuildErrors("type", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("type", errors.New("required field is null"))...)
		
		}
		delete(fields, "type")
	} else {errs = append(errs, cog.MakeBuildErrors("type", errors.New("required field is missing from input"))...)
	}
	// Field "contains"
	if fields["contains"] != nil {
		if string(fields["contains"]) != "null" {
			if err := json.Unmarshal(fields["contains"], &resource.Contains); err != nil {
				errs = append(errs, cog.MakeBuildErrors("contains", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("contains", errors.New("required field is null"))...)
		
		}
		delete(fields, "contains")
	} else {errs = append(errs, cog.MakeBuildErrors("contains", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Common", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `Common` objects.
func (resource Common) Equals(other Common) bool {
		if resource.Name != other.Name {
			return false
		}
		if resource.Type != other.Type {
			return false
		}
		if resource.Contains != other.Contains {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Common` fields for violations and returns them.
func (resource Common) Validate() error {
	return nil
}


// Counter metric combining common properties with specific values
type Counter struct {
	Common

	// Counter metric values
Values struct {
    // Total count of events
    Count float64 `json:"count"`
} `json:"values"`
}

type CommonType string
const (
	CommonTypeCounter CommonType = "counter"
	CommonTypeGauge CommonType = "gauge"
)


type CommonContains string
const (
	CommonContainsDefault CommonContains = "default"
	CommonContainsTime CommonContains = "time"
)


