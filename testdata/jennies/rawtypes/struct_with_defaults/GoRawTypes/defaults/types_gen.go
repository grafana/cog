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

func (resource *SomeStruct) StrictUnmarshalJSON(raw []byte) error {
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
	// Field "fieldString"
	if fields["fieldString"] != nil {
		if string(fields["fieldString"]) != "null" {
			if err := json.Unmarshal(fields["fieldString"], &resource.FieldString); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fieldString", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("fieldString", errors.New("required field is null"))...)
		
		}
		delete(fields, "fieldString")
	} else {errs = append(errs, cog.MakeBuildErrors("fieldString", errors.New("required field is missing from input"))...)
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
	} else {errs = append(errs, cog.MakeBuildErrors("FieldFloat32", errors.New("required field is missing from input"))...)
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
	} else {errs = append(errs, cog.MakeBuildErrors("FieldInt32", errors.New("required field is missing from input"))...)
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


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeStruct) Validate() error {
	return nil
}


