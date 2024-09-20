package struct_optional_fields

type SomeStruct struct {
	FieldRef *SomeOtherStruct `json:"FieldRef,omitempty"`
	FieldString *string `json:"FieldString,omitempty"`
	Operator *SomeStructOperator `json:"Operator,omitempty"`
	FieldArrayOfStrings []string `json:"FieldArrayOfStrings,omitempty"`
	FieldAnonymousStruct *struct {
	FieldAny any `json:"FieldAny"`
} `json:"FieldAnonymousStruct,omitempty"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		if !((resource.FieldRef == nil && other.FieldRef == nil) || (resource.FieldRef != nil && other.FieldRef != nil)) {
			return false
		}

		if resource.FieldRef != nil {
		if !resource.FieldRef.Equals(*other.FieldRef) {
			return false
		}
		}
		if !((resource.FieldString == nil && other.FieldString == nil) || (resource.FieldString != nil && other.FieldString != nil)) {
			return false
		}

		if resource.FieldString != nil {
		if *resource.FieldString != *other.FieldString {
			return false
		}
		}
		if !((resource.Operator == nil && other.Operator == nil) || (resource.Operator != nil && other.Operator != nil)) {
			return false
		}

		if resource.Operator != nil {
		if *resource.Operator != *other.Operator {
			return false
		}
		}

		if len(resource.FieldArrayOfStrings) != len(other.FieldArrayOfStrings) {
			return false
		}

		for i1 := range resource.FieldArrayOfStrings {
		if resource.FieldArrayOfStrings[i1] != other.FieldArrayOfStrings[i1] {
			return false
		}
		}
		if !((resource.FieldAnonymousStruct == nil && other.FieldAnonymousStruct == nil) || (resource.FieldAnonymousStruct != nil && other.FieldAnonymousStruct != nil)) {
			return false
		}

		if resource.FieldAnonymousStruct != nil {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAnonymousStruct.FieldAny, other.FieldAnonymousStruct.FieldAny) {
			return false
		}
		}

	return true
}


type SomeOtherStruct struct {
	FieldAny any `json:"FieldAny"`
}

func (resource SomeOtherStruct) Equals(other SomeOtherStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


type SomeStructOperator string
const (
	SomeStructOperatorGreaterThan SomeStructOperator = ">"
	SomeStructOperatorLessThan SomeStructOperator = "<"
)


