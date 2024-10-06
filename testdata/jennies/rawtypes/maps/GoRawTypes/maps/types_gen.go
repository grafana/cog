package maps

// String to... something.
type MapOfStringToAny map[string]any

type MapOfStringToString map[string]string

type SomeStruct struct {
	FieldAny any `json:"FieldAny" yaml:"FieldAny"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


type MapOfStringToRef map[string]SomeStruct

type MapOfStringToMapOfStringToBool map[string]map[string]bool

