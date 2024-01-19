package intersections

import (
	externalPkg "github.com/grafana/cog/generated/externalPkg"
)

type Intersections struct {
	SomeStruct
	externalPkg.AnotherStruct

	FieldString string `json:"fieldString"`
	FieldInteger int32 `json:"fieldInteger"`
}

type SomeStruct struct {
	FieldBool bool `json:"fieldBool"`
}

