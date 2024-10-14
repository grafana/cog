package maps

import (
	"reflect"
)

// String to... something.
type MapOfStringToAny map[string]any

type MapOfStringToString map[string]string

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


type MapOfStringToRef map[string]SomeStruct

type MapOfStringToMapOfStringToBool map[string]map[string]bool

