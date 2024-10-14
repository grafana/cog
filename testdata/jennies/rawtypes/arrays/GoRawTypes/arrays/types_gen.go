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


func (resource SomeStruct) Validate() error {
	return nil
}


type ArrayOfRefs []SomeStruct

type ArrayOfArrayOfNumbers [][]int64

