package arrays

import (
	"reflect"
)

// List of tags, maybe?
type ArrayOfStrings []string

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


type ArrayOfRefs []SomeStruct

type ArrayOfArrayOfNumbers [][]int64

