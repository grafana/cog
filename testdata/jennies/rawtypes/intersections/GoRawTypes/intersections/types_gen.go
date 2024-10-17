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
	} else {errs = append(errs, cog.MakeBuildErrors("fieldBool", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource SomeStruct) Equals(other SomeStruct) bool {
		if resource.FieldBool != other.FieldBool {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeStruct) Validate() error {
	return nil
}


