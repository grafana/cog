package intersections

import (
	externalpkg "github.com/grafana/cog/generated/externalpkg"
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


