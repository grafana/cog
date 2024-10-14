package intersections

import (
	externalpkg "github.com/grafana/cog/generated/externalpkg"
	cog "github.com/grafana/cog/generated/cog"
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


func (resource SomeStruct) Validate() error {
	var errs cog.BuildErrors

	if len(errs) == 0 {
		return nil
	}

	return errs
}


