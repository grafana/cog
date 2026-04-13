package disjunction_of_numbers

import (
	"encoding/json"
	"errors"
	cog "github.com/grafana/cog/generated/cog"
)

type Numbers = Int64OrFloat64OrFloat32

// NewNumbers creates a new Numbers object.
func NewNumbers() *Numbers {
	return NewInt64OrFloat64OrFloat32()
}
type Int64OrFloat64OrFloat32 struct {
    Int64 *int64 `json:"Int64,omitempty"`
    Float64 *float64 `json:"Float64,omitempty"`
    Float32 *float32 `json:"Float32,omitempty"`
}

// NewInt64OrFloat64OrFloat32 creates a new Int64OrFloat64OrFloat32 object.
func NewInt64OrFloat64OrFloat32() *Int64OrFloat64OrFloat32 {
	return &Int64OrFloat64OrFloat32{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `Int64OrFloat64OrFloat32` as JSON.
func (resource Int64OrFloat64OrFloat32) MarshalJSON() ([]byte, error) {
	if resource.Int64 != nil {
		return json.Marshal(resource.Int64)
	}

	if resource.Float64 != nil {
		return json.Marshal(resource.Float64)
	}

	if resource.Float32 != nil {
		return json.Marshal(resource.Float32)
	}


	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `Int64OrFloat64OrFloat32` from JSON.
func (resource *Int64OrFloat64OrFloat32) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// Int64
	var Int64 int64
	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
		resource.Int64 = nil
	} else {
		resource.Int64 = &Int64
		return nil
	}

	// Float64
	var Float64 float64
	if err := json.Unmarshal(raw, &Float64); err != nil {
		errList = append(errList, err)
		resource.Float64 = nil
	} else {
		resource.Float64 = &Float64
		return nil
	}

	// Float32
	var Float32 float32
	if err := json.Unmarshal(raw, &Float32); err != nil {
		errList = append(errList, err)
		resource.Float32 = nil
	} else {
		resource.Float32 = &Float32
		return nil
	}

	return errors.Join(errList...)
}


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Int64OrFloat64OrFloat32` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Int64OrFloat64OrFloat32) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// Int64
	var Int64 int64

	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
	} else {
		resource.Int64 = &Int64
		return nil
	}

	// Float64
	var Float64 float64

	if err := json.Unmarshal(raw, &Float64); err != nil {
		errList = append(errList, err)
	} else {
		resource.Float64 = &Float64
		return nil
	}

	// Float32
	var Float32 float32

	if err := json.Unmarshal(raw, &Float32); err != nil {
		errList = append(errList, err)
	} else {
		resource.Float32 = &Float32
		return nil
	}


	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("Int64OrFloat64OrFloat32", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `Int64OrFloat64OrFloat32` objects.
func (resource Int64OrFloat64OrFloat32) Equals(other Int64OrFloat64OrFloat32) bool {
		if resource.Int64 == nil && other.Int64 != nil || resource.Int64 != nil && other.Int64 == nil {
			return false
		}

		if resource.Int64 != nil {
		if *resource.Int64 != *other.Int64 {
			return false
		}
		}
		if resource.Float64 == nil && other.Float64 != nil || resource.Float64 != nil && other.Float64 == nil {
			return false
		}

		if resource.Float64 != nil {
		if *resource.Float64 != *other.Float64 {
			return false
		}
		}
		if resource.Float32 == nil && other.Float32 != nil || resource.Float32 != nil && other.Float32 == nil {
			return false
		}

		if resource.Float32 != nil {
		if *resource.Float32 != *other.Float32 {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Int64OrFloat64OrFloat32` fields for violations and returns them.
func (resource Int64OrFloat64OrFloat32) Validate() error {
	return nil
}


