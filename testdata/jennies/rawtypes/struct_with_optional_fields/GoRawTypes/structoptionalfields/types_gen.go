package structoptionalfields

type SomeStruct struct {
	FieldRef *SomeOtherStruct `json:"FieldRef,omitempty"`
	FieldString *string `json:"FieldString,omitempty"`
	Operator *SomeStructOperator `json:"Operator,omitempty"`
	FieldArrayOfStrings []string `json:"FieldArrayOfStrings,omitempty"`
	FieldAnonymousStruct *struct {
	FieldAny any `json:"FieldAny"`
} `json:"FieldAnonymousStruct,omitempty"`
}

type SomeOtherStruct struct {
	FieldAny any `json:"FieldAny"`
}

type SomeStructOperator string
const (
	SomeStructOperatorGreaterThan SomeStructOperator = ">"
	SomeStructOperatorLessThan SomeStructOperator = "<"
)


