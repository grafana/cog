package basic

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"reflect"
	"bytes"
)

// This
// is
// a
// comment
type SomeStruct struct {
    // Anything can go in there.
    // Really, anything.
    FieldAny any `json:"FieldAny"`
    FieldBool bool `json:"FieldBool"`
    FieldBytes []byte `json:"FieldBytes"`
    FieldString string `json:"FieldString"`
    FieldStringWithConstantValue string `json:"FieldStringWithConstantValue"`
    FieldFloat32 float32 `json:"FieldFloat32"`
    FieldFloat64 float64 `json:"FieldFloat64"`
    FieldUint8 uint8 `json:"FieldUint8"`
    FieldUint16 uint16 `json:"FieldUint16"`
    FieldUint32 uint32 `json:"FieldUint32"`
    FieldUint64 uint64 `json:"FieldUint64"`
    FieldInt8 int8 `json:"FieldInt8"`
    FieldInt16 int16 `json:"FieldInt16"`
    FieldInt32 int32 `json:"FieldInt32"`
    FieldInt64 int64 `json:"FieldInt64"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		FieldStringWithConstantValue: "auto",
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
	// Field "FieldBool"
	if fields["FieldBool"] != nil {
		if string(fields["FieldBool"]) != "null" {
			if err := json.Unmarshal(fields["FieldBool"], &resource.FieldBool); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldBool", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldBool", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldBool")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldBool", errors.New("required field is missing from input"))...)
	}
	// Field "FieldBytes"
	if fields["FieldBytes"] != nil {
		if string(fields["FieldBytes"]) != "null" {
			if err := json.Unmarshal(fields["FieldBytes"], &resource.FieldBytes); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldBytes", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldBytes", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldBytes")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldBytes", errors.New("required field is missing from input"))...)
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
	// Field "FieldFloat64"
	if fields["FieldFloat64"] != nil {
		if string(fields["FieldFloat64"]) != "null" {
			if err := json.Unmarshal(fields["FieldFloat64"], &resource.FieldFloat64); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldFloat64", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldFloat64", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldFloat64")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldFloat64", errors.New("required field is missing from input"))...)
	}
	// Field "FieldUint8"
	if fields["FieldUint8"] != nil {
		if string(fields["FieldUint8"]) != "null" {
			if err := json.Unmarshal(fields["FieldUint8"], &resource.FieldUint8); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldUint8", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldUint8", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldUint8")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldUint8", errors.New("required field is missing from input"))...)
	}
	// Field "FieldUint16"
	if fields["FieldUint16"] != nil {
		if string(fields["FieldUint16"]) != "null" {
			if err := json.Unmarshal(fields["FieldUint16"], &resource.FieldUint16); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldUint16", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldUint16", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldUint16")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldUint16", errors.New("required field is missing from input"))...)
	}
	// Field "FieldUint32"
	if fields["FieldUint32"] != nil {
		if string(fields["FieldUint32"]) != "null" {
			if err := json.Unmarshal(fields["FieldUint32"], &resource.FieldUint32); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldUint32", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldUint32", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldUint32")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldUint32", errors.New("required field is missing from input"))...)
	}
	// Field "FieldUint64"
	if fields["FieldUint64"] != nil {
		if string(fields["FieldUint64"]) != "null" {
			if err := json.Unmarshal(fields["FieldUint64"], &resource.FieldUint64); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldUint64", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldUint64", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldUint64")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldUint64", errors.New("required field is missing from input"))...)
	}
	// Field "FieldInt8"
	if fields["FieldInt8"] != nil {
		if string(fields["FieldInt8"]) != "null" {
			if err := json.Unmarshal(fields["FieldInt8"], &resource.FieldInt8); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldInt8", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldInt8", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldInt8")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldInt8", errors.New("required field is missing from input"))...)
	}
	// Field "FieldInt16"
	if fields["FieldInt16"] != nil {
		if string(fields["FieldInt16"]) != "null" {
			if err := json.Unmarshal(fields["FieldInt16"], &resource.FieldInt16); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldInt16", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldInt16", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldInt16")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldInt16", errors.New("required field is missing from input"))...)
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
	// Field "FieldInt64"
	if fields["FieldInt64"] != nil {
		if string(fields["FieldInt64"]) != "null" {
			if err := json.Unmarshal(fields["FieldInt64"], &resource.FieldInt64); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldInt64", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldInt64", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldInt64")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldInt64", errors.New("required field is missing from input"))...)
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
		if resource.FieldBool != other.FieldBool {
			return false
		}
	    if !bytes.Equal(resource.FieldBytes, other.FieldBytes) {
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
		if resource.FieldFloat64 != other.FieldFloat64 {
			return false
		}
		if resource.FieldUint8 != other.FieldUint8 {
			return false
		}
		if resource.FieldUint16 != other.FieldUint16 {
			return false
		}
		if resource.FieldUint32 != other.FieldUint32 {
			return false
		}
		if resource.FieldUint64 != other.FieldUint64 {
			return false
		}
		if resource.FieldInt8 != other.FieldInt8 {
			return false
		}
		if resource.FieldInt16 != other.FieldInt16 {
			return false
		}
		if resource.FieldInt32 != other.FieldInt32 {
			return false
		}
		if resource.FieldInt64 != other.FieldInt64 {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


