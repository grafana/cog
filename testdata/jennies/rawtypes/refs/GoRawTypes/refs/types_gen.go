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


type RefToSomeStruct = SomeStruct

type RefToSomeStructFromOtherPackage = otherpkg.SomeDistantStruct

