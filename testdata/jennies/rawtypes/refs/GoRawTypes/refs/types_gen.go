package refs

import (
	"reflect"
	cog "github.com/grafana/cog/generated/cog"
	otherpkg "github.com/grafana/cog/generated/otherpkg"
)

type SomeStruct struct {
	FieldAny any `json:"FieldAny"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
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


type RefToSomeStruct = SomeStruct

type RefToSomeStructFromOtherPackage = otherpkg.SomeDistantStruct

