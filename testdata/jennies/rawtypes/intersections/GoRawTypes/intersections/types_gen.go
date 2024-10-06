package intersections

import (
	externalpkg "github.com/grafana/cog/generated/externalpkg"
)

type Intersections struct {
	SomeStruct
	externalpkg.AnotherStruct

	FieldString string `json:"fieldString" yaml:"fieldString"`
	FieldInteger int32 `json:"fieldInteger" yaml:"fieldInteger"`
}

type SomeStruct struct {
	FieldBool bool `json:"fieldBool" yaml:"fieldBool"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		if resource.FieldBool != other.FieldBool {
			return false
		}

	return true
}


