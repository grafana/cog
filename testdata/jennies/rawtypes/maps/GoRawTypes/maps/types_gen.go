package maps

// String to... something.
type MapOfStringToAny map[string]any

type MapOfStringToString map[string]string

type SomeStruct struct {
	FieldAny any `json:"FieldAny"`
}

type MapOfStringToRef map[string]SomeStruct

type MapOfStringToMapOfStringToBool map[string]map[string]bool

