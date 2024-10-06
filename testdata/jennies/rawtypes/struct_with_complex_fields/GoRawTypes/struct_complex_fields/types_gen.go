package struct_complex_fields

// This struct does things.
type SomeStruct struct {
	FieldRef SomeOtherStruct `json:"FieldRef" yaml:"FieldRef"`
	FieldDisjunctionOfScalars StringOrBool `json:"FieldDisjunctionOfScalars" yaml:"FieldDisjunctionOfScalars"`
	FieldMixedDisjunction StringOrSomeOtherStruct `json:"FieldMixedDisjunction" yaml:"FieldMixedDisjunction"`
	FieldDisjunctionWithNull *string `json:"FieldDisjunctionWithNull" yaml:"FieldDisjunctionWithNull"`
	Operator SomeStructOperator `json:"Operator" yaml:"Operator"`
	FieldArrayOfStrings []string `json:"FieldArrayOfStrings" yaml:"FieldArrayOfStrings"`
	FieldMapOfStringToString map[string]string `json:"FieldMapOfStringToString" yaml:"FieldMapOfStringToString"`
	FieldAnonymousStruct StructComplexFieldsSomeStructFieldAnonymousStruct `json:"FieldAnonymousStruct" yaml:"FieldAnonymousStruct"`
	FieldRefToConstant string `json:"fieldRefToConstant" yaml:"fieldRefToConstant"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		if !resource.FieldRef.Equals(other.FieldRef) {
			return false
		}
		if !resource.FieldDisjunctionOfScalars.Equals(other.FieldDisjunctionOfScalars) {
			return false
		}
		if !resource.FieldMixedDisjunction.Equals(other.FieldMixedDisjunction) {
			return false
		}
		if resource.FieldDisjunctionWithNull == nil && other.FieldDisjunctionWithNull != nil || resource.FieldDisjunctionWithNull != nil && other.FieldDisjunctionWithNull == nil {
			return false
		}

		if resource.FieldDisjunctionWithNull != nil {
		if *resource.FieldDisjunctionWithNull != *other.FieldDisjunctionWithNull {
			return false
		}
		}
		if resource.Operator != other.Operator {
			return false
		}

		if len(resource.FieldArrayOfStrings) != len(other.FieldArrayOfStrings) {
			return false
		}

		for i1 := range resource.FieldArrayOfStrings {
		if resource.FieldArrayOfStrings[i1] != other.FieldArrayOfStrings[i1] {
			return false
		}
		}

		if len(resource.FieldMapOfStringToString) != len(other.FieldMapOfStringToString) {
			return false
		}

		for key1 := range resource.FieldMapOfStringToString {
		if resource.FieldMapOfStringToString[key1] != other.FieldMapOfStringToString[key1] {
			return false
		}
		}
		if !resource.FieldAnonymousStruct.Equals(other.FieldAnonymousStruct) {
			return false
		}
		if resource.FieldRefToConstant != other.FieldRefToConstant {
			return false
		}

	return true
}


const ConnectionPath = "straight"

type SomeOtherStruct struct {
	FieldAny any `json:"FieldAny" yaml:"FieldAny"`
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


type StructComplexFieldsSomeStructFieldAnonymousStruct struct {
	FieldAny any `json:"FieldAny" yaml:"FieldAny"`
}

func (resource StructComplexFieldsSomeStructFieldAnonymousStruct) Equals(other StructComplexFieldsSomeStructFieldAnonymousStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


type StringOrBool struct {
	String *string `json:"String,omitempty" yaml:"String,omitempty"`
	Bool *bool `json:"Bool,omitempty" yaml:"Bool,omitempty"`
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

func (resource StringOrBool) MarshalYAML() (any, error) {
	if resource.String != nil {
		return resource.String, nil
	}

	if resource.Bool != nil {
		return resource.Bool, nil
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


func (resource StringOrBool) Equals(other StringOrBool) bool {
		if resource.String == nil && other.String != nil || resource.String != nil && other.String == nil {
			return false
		}

		if resource.String != nil {
		if *resource.String != *other.String {
			return false
		}
		}
		if resource.Bool == nil && other.Bool != nil || resource.Bool != nil && other.Bool == nil {
			return false
		}

		if resource.Bool != nil {
		if *resource.Bool != *other.Bool {
			return false
		}
		}

	return true
}


type StringOrSomeOtherStruct struct {
	String *string `json:"String,omitempty" yaml:"String,omitempty"`
	SomeOtherStruct *SomeOtherStruct `json:"SomeOtherStruct,omitempty" yaml:"SomeOtherStruct,omitempty"`
}

func (resource StringOrSomeOtherStruct) Equals(other StringOrSomeOtherStruct) bool {
		if resource.String == nil && other.String != nil || resource.String != nil && other.String == nil {
			return false
		}

		if resource.String != nil {
		if *resource.String != *other.String {
			return false
		}
		}
		if resource.SomeOtherStruct == nil && other.SomeOtherStruct != nil || resource.SomeOtherStruct != nil && other.SomeOtherStruct == nil {
			return false
		}

		if resource.SomeOtherStruct != nil {
		if !resource.SomeOtherStruct.Equals(*other.SomeOtherStruct) {
			return false
		}
		}

	return true
}


