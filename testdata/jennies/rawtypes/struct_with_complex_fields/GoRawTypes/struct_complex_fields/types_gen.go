package struct_complex_fields

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
	FieldRefToConstant string `json:"fieldRefToConstant"`
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

func (resource StringOrBool) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}

	return nil, fmt.Errorf("no value for disjunction of scalars")
}


func (resource *StringOrBool) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// String
	var String string
	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
		resource.String = nil
	} else {
		resource.String = &String
		return nil
	}

	// Bool
	var Bool bool
	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
		resource.Bool = nil
	} else {
		resource.Bool = &Bool
		return nil
	}

	return errors.Join(errList...)
}


type StringOrSomeOtherStruct struct {
	String *string `json:"String,omitempty"`
	SomeOtherStruct *SomeOtherStruct `json:"SomeOtherStruct,omitempty"`
}

