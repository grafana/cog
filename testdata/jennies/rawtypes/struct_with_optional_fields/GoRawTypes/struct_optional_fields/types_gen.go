package struct_optional_fields

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"fmt"
	"errors"
	"reflect"
)

type SomeStruct struct {
	FieldRef *SomeOtherStruct `json:"FieldRef,omitempty"`
	FieldString *string `json:"FieldString,omitempty"`
	Operator *SomeStructOperator `json:"Operator,omitempty"`
	FieldArrayOfStrings []string `json:"FieldArrayOfStrings,omitempty"`
	FieldAnonymousStruct *StructOptionalFieldsSomeStructFieldAnonymousStruct `json:"FieldAnonymousStruct,omitempty"`
}

func (resource *SomeStruct) StrictUnmarshalJSON(raw []byte) error {
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
			
			resource.FieldRef = &SomeOtherStruct{}
			if err := resource.FieldRef.StrictUnmarshalJSON(fields["FieldRef"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldRef", err)...)
			}
		
		}
		delete(fields, "FieldRef")
	
	}
	// Field "FieldString"
	if fields["FieldString"] != nil {
		if string(fields["FieldString"]) != "null" {
			if err := json.Unmarshal(fields["FieldString"], &resource.FieldString); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldString", err)...)
			}
		
		}
		delete(fields, "FieldString")
	
	}
	// Field "Operator"
	if fields["Operator"] != nil {
		if string(fields["Operator"]) != "null" {
			if err := json.Unmarshal(fields["Operator"], &resource.Operator); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Operator", err)...)
			}
		
		}
		delete(fields, "Operator")
	
	}
	// Field "FieldArrayOfStrings"
	if fields["FieldArrayOfStrings"] != nil {
		if string(fields["FieldArrayOfStrings"]) != "null" {
			
			if err := json.Unmarshal(fields["FieldArrayOfStrings"], &resource.FieldArrayOfStrings); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldArrayOfStrings", err)...)
			}
		
		}
		delete(fields, "FieldArrayOfStrings")
	
	}
	// Field "FieldAnonymousStruct"
	if fields["FieldAnonymousStruct"] != nil {
		if string(fields["FieldAnonymousStruct"]) != "null" {
			
			resource.FieldAnonymousStruct = &StructOptionalFieldsSomeStructFieldAnonymousStruct{}
			if err := resource.FieldAnonymousStruct.StrictUnmarshalJSON(fields["FieldAnonymousStruct"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("FieldAnonymousStruct", err)...)
			}
		
		}
		delete(fields, "FieldAnonymousStruct")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource SomeStruct) Equals(other SomeStruct) bool {
		if resource.FieldRef == nil && other.FieldRef != nil || resource.FieldRef != nil && other.FieldRef == nil {
			return false
		}

		if resource.FieldRef != nil {
		if !resource.FieldRef.Equals(*other.FieldRef) {
			return false
		}
		}
		if resource.FieldString == nil && other.FieldString != nil || resource.FieldString != nil && other.FieldString == nil {
			return false
		}

		if resource.FieldString != nil {
		if *resource.FieldString != *other.FieldString {
			return false
		}
		}
		if resource.Operator == nil && other.Operator != nil || resource.Operator != nil && other.Operator == nil {
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
		if resource.FieldAnonymousStruct == nil && other.FieldAnonymousStruct != nil || resource.FieldAnonymousStruct != nil && other.FieldAnonymousStruct == nil {
			return false
		}

		if resource.FieldAnonymousStruct != nil {
		if !resource.FieldAnonymousStruct.Equals(*other.FieldAnonymousStruct) {
			return false
		}
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeStruct) Validate() error {
	var errs cog.BuildErrors
		if resource.FieldRef != nil {
		if err := resource.FieldRef.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("FieldRef", err)...)
		}
		}
		if resource.FieldAnonymousStruct != nil {
		if err := resource.FieldAnonymousStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("FieldAnonymousStruct", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type SomeOtherStruct struct {
	FieldAny any `json:"FieldAny"`
}

func (resource *SomeOtherStruct) StrictUnmarshalJSON(raw []byte) error {
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


func (resource SomeOtherStruct) Equals(other SomeOtherStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeOtherStruct) Validate() error {
	return nil
}


type SomeStructOperator string
const (
	SomeStructOperatorGreaterThan SomeStructOperator = ">"
	SomeStructOperatorLessThan SomeStructOperator = "<"
)


type StructOptionalFieldsSomeStructFieldAnonymousStruct struct {
	FieldAny any `json:"FieldAny"`
}

func (resource *StructOptionalFieldsSomeStructFieldAnonymousStruct) StrictUnmarshalJSON(raw []byte) error {
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
		errs = append(errs, cog.MakeBuildErrors("StructOptionalFieldsSomeStructFieldAnonymousStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource StructOptionalFieldsSomeStructFieldAnonymousStruct) Equals(other StructOptionalFieldsSomeStructFieldAnonymousStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource StructOptionalFieldsSomeStructFieldAnonymousStruct) Validate() error {
	return nil
}


