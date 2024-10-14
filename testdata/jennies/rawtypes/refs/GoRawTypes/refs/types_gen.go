package refs

import (
	"reflect"
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


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeStruct) Validate() error {
	return nil
}


type RefToSomeStruct = SomeStruct

type RefToSomeStructFromOtherPackage = otherpkg.SomeDistantStruct

