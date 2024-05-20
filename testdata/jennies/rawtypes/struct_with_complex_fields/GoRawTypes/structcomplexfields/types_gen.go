package structcomplexfields

// This struct does things.
type SomeStruct struct {
	FieldRef SomeOtherStruct `json:"FieldRef"`
	FieldDisjunctionOfScalars StringOrBool `json:"FieldDisjunctionOfScalars"`
	FieldMixedDisjunction StringOrSomeOtherStruct `json:"FieldMixedDisjunction"`
	FieldDisjunctionWithNull *string `json:"FieldDisjunctionWithNull"`
	Operator SomeStructOperator `json:"Operator"`
	FieldArrayOfStrings []string `json:"FieldArrayOfStrings"`
	FieldMapOfStringToString map[string]string `json:"FieldMapOfStringToString"`
	FieldAnonymousStruct struct {
	FieldAny any `json:"FieldAny"`
} `json:"FieldAnonymousStruct"`
	FieldRefToConstant ConnectionPath `json:"fieldRefToConstant"`
}

const ConnectionPath = "straight"

type SomeOtherStruct struct {
	FieldAny any `json:"FieldAny"`
}

type SomeStructOperator string
const (
	SomeStructOperatorGreaterThan SomeStructOperator = ">"
	SomeStructOperatorLessThan SomeStructOperator = "<"
)


type StringOrBool struct {
	String *string `json:"String,omitempty"`
	Bool *bool `json:"Bool,omitempty"`
}

type StringOrSomeOtherStruct struct {
	String *string `json:"String,omitempty"`
	SomeOtherStruct *SomeOtherStruct `json:"SomeOtherStruct,omitempty"`
}

