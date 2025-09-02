package struct_complex_fields

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"reflect"
	"bytes"
)

// This struct does things.
type SomeStruct struct {
    FieldRef SomeOtherStruct `json:"FieldRef"`
    FieldDisjunctionOfScalars StringOrBool `json:"FieldDisjunctionOfScalars"`
    FieldMixedDisjunction StringOrSomeOtherStruct `json:"FieldMixedDisjunction"`
    FieldDisjunctionWithNull *string `json:"FieldDisjunctionWithNull"`
    Operator SomeStructOperator `json:"Operator"`
    FieldArrayOfStrings []string `json:"FieldArrayOfStrings"`
    FieldMapOfStringToString map[string]string `json:"FieldMapOfStringToString"`
    FieldAnonymousStruct StructComplexFieldsSomeStructFieldAnonymousStruct `json:"FieldAnonymousStruct"`
    FieldRefToConstant string `json:"fieldRefToConstant"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		FieldRef: *NewSomeOtherStruct(),
		FieldDisjunctionOfScalars: *NewStringOrBool(),
		FieldMixedDisjunction: *NewStringOrSomeOtherStruct(),
		FieldArrayOfStrings: []string{},
		FieldMapOfStringToString: map[string]string{},
		FieldAnonymousStruct: *NewStructComplexFieldsSomeStructFieldAnonymousStruct(),
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `SomeStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *SomeStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "FieldRef"
	if fields["FieldRef"] != nil {
		if string(fields["FieldRef"]) != "null" {
			
			resource.FieldRef = SomeOtherStruct{}
			if err := resource.FieldRef.UnmarshalJSONStrict(fields["FieldRef"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldRef", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldRef", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldRef")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldRef", errors.New("required field is missing from input"))...)
	}
	// Field "FieldDisjunctionOfScalars"
	if fields["FieldDisjunctionOfScalars"] != nil {
		if string(fields["FieldDisjunctionOfScalars"]) != "null" {
			
			resource.FieldDisjunctionOfScalars = StringOrBool{}
			if err := resource.FieldDisjunctionOfScalars.UnmarshalJSONStrict(fields["FieldDisjunctionOfScalars"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldDisjunctionOfScalars", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldDisjunctionOfScalars", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldDisjunctionOfScalars")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldDisjunctionOfScalars", errors.New("required field is missing from input"))...)
	}
	// Field "FieldMixedDisjunction"
	if fields["FieldMixedDisjunction"] != nil {
		if string(fields["FieldMixedDisjunction"]) != "null" {
			
			resource.FieldMixedDisjunction = StringOrSomeOtherStruct{}
			if err := resource.FieldMixedDisjunction.UnmarshalJSONStrict(fields["FieldMixedDisjunction"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldMixedDisjunction", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldMixedDisjunction", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldMixedDisjunction")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldMixedDisjunction", errors.New("required field is missing from input"))...)
	}
	// Field "FieldDisjunctionWithNull"
	if fields["FieldDisjunctionWithNull"] != nil {
		if string(fields["FieldDisjunctionWithNull"]) != "null" {
			if err := json.Unmarshal(fields["FieldDisjunctionWithNull"], &resource.FieldDisjunctionWithNull); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldDisjunctionWithNull", err)...)
			}
		
		}
		delete(fields, "FieldDisjunctionWithNull")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldDisjunctionWithNull", errors.New("required field is missing from input"))...)
	}
	// Field "Operator"
	if fields["Operator"] != nil {
		if string(fields["Operator"]) != "null" {
			if err := json.Unmarshal(fields["Operator"], &resource.Operator); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Operator", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("Operator", errors.New("required field is null"))...)
		
		}
		delete(fields, "Operator")
	} else {errs = append(errs, cog.MakeBuildErrors("Operator", errors.New("required field is missing from input"))...)
	}
	// Field "FieldArrayOfStrings"
	if fields["FieldArrayOfStrings"] != nil {
		if string(fields["FieldArrayOfStrings"]) != "null" {
			
			if err := json.Unmarshal(fields["FieldArrayOfStrings"], &resource.FieldArrayOfStrings); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldArrayOfStrings", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldArrayOfStrings", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldArrayOfStrings")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldArrayOfStrings", errors.New("required field is missing from input"))...)
	}
	// Field "FieldMapOfStringToString"
	if fields["FieldMapOfStringToString"] != nil {
		if string(fields["FieldMapOfStringToString"]) != "null" {
			
			if err := json.Unmarshal(fields["FieldMapOfStringToString"], &resource.FieldMapOfStringToString); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldMapOfStringToString", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldMapOfStringToString", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldMapOfStringToString")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldMapOfStringToString", errors.New("required field is missing from input"))...)
	}
	// Field "FieldAnonymousStruct"
	if fields["FieldAnonymousStruct"] != nil {
		if string(fields["FieldAnonymousStruct"]) != "null" {
			
			resource.FieldAnonymousStruct = StructComplexFieldsSomeStructFieldAnonymousStruct{}
			if err := resource.FieldAnonymousStruct.UnmarshalJSONStrict(fields["FieldAnonymousStruct"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldAnonymousStruct", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldAnonymousStruct", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldAnonymousStruct")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldAnonymousStruct", errors.New("required field is missing from input"))...)
	}
	// Field "fieldRefToConstant"
	if fields["fieldRefToConstant"] != nil {
		if string(fields["fieldRefToConstant"]) != "null" {
			if err := json.Unmarshal(fields["fieldRefToConstant"], &resource.FieldRefToConstant); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fieldRefToConstant", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("fieldRefToConstant", errors.New("required field is null"))...)
		
		}
		delete(fields, "fieldRefToConstant")
	} else {errs = append(errs, cog.MakeBuildErrors("fieldRefToConstant", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `SomeStruct` objects.
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


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	var errs cog.BuildErrors
		if err := resource.FieldRef.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("FieldRef", err)...)
		}
		if err := resource.FieldDisjunctionOfScalars.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("FieldDisjunctionOfScalars", err)...)
		}
		if err := resource.FieldMixedDisjunction.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("FieldMixedDisjunction", err)...)
		}
		if err := resource.FieldAnonymousStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("FieldAnonymousStruct", err)...)
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


const ConnectionPath = "straight"

type SomeOtherStruct struct {
    FieldAny any `json:"FieldAny"`
}

// NewSomeOtherStruct creates a new SomeOtherStruct object.
func NewSomeOtherStruct() *SomeOtherStruct {
	return &SomeOtherStruct{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `SomeOtherStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *SomeOtherStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "FieldAny"
	if fields["FieldAny"] != nil {
		if string(fields["FieldAny"]) != "null" {
			if err := json.Unmarshal(fields["FieldAny"], &resource.FieldAny); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldAny", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldAny", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldAny")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldAny", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeOtherStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `SomeOtherStruct` objects.
func (resource SomeOtherStruct) Equals(other SomeOtherStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeOtherStruct` fields for violations and returns them.
func (resource SomeOtherStruct) Validate() error {
	return nil
}


type StructComplexFieldsSomeStructFieldAnonymousStruct struct {
    FieldAny any `json:"FieldAny"`
}

// NewStructComplexFieldsSomeStructFieldAnonymousStruct creates a new StructComplexFieldsSomeStructFieldAnonymousStruct object.
func NewStructComplexFieldsSomeStructFieldAnonymousStruct() *StructComplexFieldsSomeStructFieldAnonymousStruct {
	return &StructComplexFieldsSomeStructFieldAnonymousStruct{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StructComplexFieldsSomeStructFieldAnonymousStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StructComplexFieldsSomeStructFieldAnonymousStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "FieldAny"
	if fields["FieldAny"] != nil {
		if string(fields["FieldAny"]) != "null" {
			if err := json.Unmarshal(fields["FieldAny"], &resource.FieldAny); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldAny", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("FieldAny", errors.New("required field is null"))...)
		
		}
		delete(fields, "FieldAny")
	} else {errs = append(errs, cog.MakeBuildErrors("FieldAny", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("StructComplexFieldsSomeStructFieldAnonymousStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `StructComplexFieldsSomeStructFieldAnonymousStruct` objects.
func (resource StructComplexFieldsSomeStructFieldAnonymousStruct) Equals(other StructComplexFieldsSomeStructFieldAnonymousStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StructComplexFieldsSomeStructFieldAnonymousStruct` fields for violations and returns them.
func (resource StructComplexFieldsSomeStructFieldAnonymousStruct) Validate() error {
	return nil
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

// NewStringOrBool creates a new StringOrBool object.
func NewStringOrBool() *StringOrBool {
	return &StringOrBool{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `StringOrBool` as JSON.
func (resource StringOrBool) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}


	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StringOrBool` from JSON.
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


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StringOrBool` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StringOrBool) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// String
	var String string

	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
	} else {
		resource.String = &String
		return nil
	}

	// Bool
	var Bool bool

	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
	} else {
		resource.Bool = &Bool
		return nil
	}


	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StringOrBool", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StringOrBool` objects.
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


// Validate checks all the validation constraints that may be defined on `StringOrBool` fields for violations and returns them.
func (resource StringOrBool) Validate() error {
	return nil
}


type StringOrSomeOtherStruct struct {
    String *string `json:"String,omitempty"`
    SomeOtherStruct *SomeOtherStruct `json:"SomeOtherStruct,omitempty"`
}

// NewStringOrSomeOtherStruct creates a new StringOrSomeOtherStruct object.
func NewStringOrSomeOtherStruct() *StringOrSomeOtherStruct {
	return &StringOrSomeOtherStruct{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `StringOrSomeOtherStruct` as JSON.
func (resource StringOrSomeOtherStruct) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}
	if resource.SomeOtherStruct != nil {
		return json.Marshal(resource.SomeOtherStruct)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StringOrSomeOtherStruct` from JSON.
func (resource *StringOrSomeOtherStruct) UnmarshalJSON(raw []byte) error {
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

	// SomeOtherStruct
	var SomeOtherStruct SomeOtherStruct
    someOtherStructdec := json.NewDecoder(bytes.NewReader(raw))
    someOtherStructdec.DisallowUnknownFields()
    if err := someOtherStructdec.Decode(&SomeOtherStruct); err != nil {
        errList = append(errList, err)
        resource.SomeOtherStruct = nil
    } else {
        resource.SomeOtherStruct = &SomeOtherStruct
        return nil
    }

	return errors.Join(errList...)
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StringOrSomeOtherStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StringOrSomeOtherStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// String
	var String string
	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
	} else {
		resource.String = &String
		return nil
	}

	// SomeOtherStruct
	var SomeOtherStruct SomeOtherStruct
    someOtherStructdec := json.NewDecoder(bytes.NewReader(raw))
    someOtherStructdec.DisallowUnknownFields()
    if err := someOtherStructdec.Decode(&SomeOtherStruct); err != nil {
        errList = append(errList, err)
    } else {
        resource.SomeOtherStruct = &SomeOtherStruct
        return nil
    }

	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StringOrSomeOtherStruct", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StringOrSomeOtherStruct` objects.
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


// Validate checks all the validation constraints that may be defined on `StringOrSomeOtherStruct` fields for violations and returns them.
func (resource StringOrSomeOtherStruct) Validate() error {
	var errs cog.BuildErrors
		if resource.SomeOtherStruct != nil {
		if err := resource.SomeOtherStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("SomeOtherStruct", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


