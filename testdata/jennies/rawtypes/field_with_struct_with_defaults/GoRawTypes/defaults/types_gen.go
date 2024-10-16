package defaults

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type NestedStruct struct {
	StringVal string `json:"stringVal"`
	IntVal int64 `json:"intVal"`
}

func (resource *NestedStruct) StrictUnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "stringVal"
	if fields["stringVal"] != nil {
		if string(fields["stringVal"]) != "null" {
			if err := json.Unmarshal(fields["stringVal"], &resource.StringVal); err != nil {
				errs = append(errs, cog.MakeBuildErrors("stringVal", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("stringVal", errors.New("required field is null"))...)
		
		}
		delete(fields, "stringVal")
	} else {errs = append(errs, cog.MakeBuildErrors("stringVal", errors.New("required field is missing from input"))...)
	}
	// Field "intVal"
	if fields["intVal"] != nil {
		if string(fields["intVal"]) != "null" {
			if err := json.Unmarshal(fields["intVal"], &resource.IntVal); err != nil {
				errs = append(errs, cog.MakeBuildErrors("intVal", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("intVal", errors.New("required field is null"))...)
		
		}
		delete(fields, "intVal")
	} else {errs = append(errs, cog.MakeBuildErrors("intVal", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("NestedStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource NestedStruct) Equals(other NestedStruct) bool {
		if resource.StringVal != other.StringVal {
			return false
		}
		if resource.IntVal != other.IntVal {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource NestedStruct) Validate() error {
	return nil
}


type Struct struct {
	AllFields NestedStruct `json:"allFields"`
	PartialFields NestedStruct `json:"partialFields"`
	EmptyFields NestedStruct `json:"emptyFields"`
	ComplexField DefaultsStructComplexField `json:"complexField"`
	PartialComplexField DefaultsStructPartialComplexField `json:"partialComplexField"`
}

func (resource *Struct) StrictUnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "allFields"
	if fields["allFields"] != nil {
		if string(fields["allFields"]) != "null" {
			
			resource.AllFields = NestedStruct{}
			if err := resource.AllFields.StrictUnmarshalJSON(fields["allFields"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("allFields", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("allFields", errors.New("required field is null"))...)
		
		}
		delete(fields, "allFields")
	} else {errs = append(errs, cog.MakeBuildErrors("allFields", errors.New("required field is missing from input"))...)
	}
	// Field "partialFields"
	if fields["partialFields"] != nil {
		if string(fields["partialFields"]) != "null" {
			
			resource.PartialFields = NestedStruct{}
			if err := resource.PartialFields.StrictUnmarshalJSON(fields["partialFields"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("partialFields", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("partialFields", errors.New("required field is null"))...)
		
		}
		delete(fields, "partialFields")
	} else {errs = append(errs, cog.MakeBuildErrors("partialFields", errors.New("required field is missing from input"))...)
	}
	// Field "emptyFields"
	if fields["emptyFields"] != nil {
		if string(fields["emptyFields"]) != "null" {
			
			resource.EmptyFields = NestedStruct{}
			if err := resource.EmptyFields.StrictUnmarshalJSON(fields["emptyFields"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("emptyFields", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("emptyFields", errors.New("required field is null"))...)
		
		}
		delete(fields, "emptyFields")
	} else {errs = append(errs, cog.MakeBuildErrors("emptyFields", errors.New("required field is missing from input"))...)
	}
	// Field "complexField"
	if fields["complexField"] != nil {
		if string(fields["complexField"]) != "null" {
			
			resource.ComplexField = DefaultsStructComplexField{}
			if err := resource.ComplexField.StrictUnmarshalJSON(fields["complexField"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("complexField", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("complexField", errors.New("required field is null"))...)
		
		}
		delete(fields, "complexField")
	} else {errs = append(errs, cog.MakeBuildErrors("complexField", errors.New("required field is missing from input"))...)
	}
	// Field "partialComplexField"
	if fields["partialComplexField"] != nil {
		if string(fields["partialComplexField"]) != "null" {
			
			resource.PartialComplexField = DefaultsStructPartialComplexField{}
			if err := resource.PartialComplexField.StrictUnmarshalJSON(fields["partialComplexField"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("partialComplexField", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("partialComplexField", errors.New("required field is null"))...)
		
		}
		delete(fields, "partialComplexField")
	} else {errs = append(errs, cog.MakeBuildErrors("partialComplexField", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Struct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource Struct) Equals(other Struct) bool {
		if !resource.AllFields.Equals(other.AllFields) {
			return false
		}
		if !resource.PartialFields.Equals(other.PartialFields) {
			return false
		}
		if !resource.EmptyFields.Equals(other.EmptyFields) {
			return false
		}
		if !resource.ComplexField.Equals(other.ComplexField) {
			return false
		}
		if !resource.PartialComplexField.Equals(other.PartialComplexField) {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource Struct) Validate() error {
	var errs cog.BuildErrors
		if err := resource.AllFields.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("allFields", err)...)
		}
		if err := resource.PartialFields.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("partialFields", err)...)
		}
		if err := resource.EmptyFields.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("emptyFields", err)...)
		}
		if err := resource.ComplexField.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("complexField", err)...)
		}
		if err := resource.PartialComplexField.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("partialComplexField", err)...)
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type DefaultsStructComplexFieldNested struct {
	NestedVal string `json:"nestedVal"`
}

func (resource *DefaultsStructComplexFieldNested) StrictUnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "nestedVal"
	if fields["nestedVal"] != nil {
		if string(fields["nestedVal"]) != "null" {
			if err := json.Unmarshal(fields["nestedVal"], &resource.NestedVal); err != nil {
				errs = append(errs, cog.MakeBuildErrors("nestedVal", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("nestedVal", errors.New("required field is null"))...)
		
		}
		delete(fields, "nestedVal")
	} else {errs = append(errs, cog.MakeBuildErrors("nestedVal", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("DefaultsStructComplexFieldNested", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource DefaultsStructComplexFieldNested) Equals(other DefaultsStructComplexFieldNested) bool {
		if resource.NestedVal != other.NestedVal {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource DefaultsStructComplexFieldNested) Validate() error {
	return nil
}


type DefaultsStructComplexField struct {
	Uid string `json:"uid"`
	Nested DefaultsStructComplexFieldNested `json:"nested"`
	Array []string `json:"array"`
}

func (resource *DefaultsStructComplexField) StrictUnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "uid"
	if fields["uid"] != nil {
		if string(fields["uid"]) != "null" {
			if err := json.Unmarshal(fields["uid"], &resource.Uid); err != nil {
				errs = append(errs, cog.MakeBuildErrors("uid", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("uid", errors.New("required field is null"))...)
		
		}
		delete(fields, "uid")
	} else {errs = append(errs, cog.MakeBuildErrors("uid", errors.New("required field is missing from input"))...)
	}
	// Field "nested"
	if fields["nested"] != nil {
		if string(fields["nested"]) != "null" {
			
			resource.Nested = DefaultsStructComplexFieldNested{}
			if err := resource.Nested.StrictUnmarshalJSON(fields["nested"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("nested", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("nested", errors.New("required field is null"))...)
		
		}
		delete(fields, "nested")
	} else {errs = append(errs, cog.MakeBuildErrors("nested", errors.New("required field is missing from input"))...)
	}
	// Field "array"
	if fields["array"] != nil {
		if string(fields["array"]) != "null" {
			
			if err := json.Unmarshal(fields["array"], &resource.Array); err != nil {
				errs = append(errs, cog.MakeBuildErrors("array", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("array", errors.New("required field is null"))...)
		
		}
		delete(fields, "array")
	} else {errs = append(errs, cog.MakeBuildErrors("array", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("DefaultsStructComplexField", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource DefaultsStructComplexField) Equals(other DefaultsStructComplexField) bool {
		if resource.Uid != other.Uid {
			return false
		}
		if !resource.Nested.Equals(other.Nested) {
			return false
		}

		if len(resource.Array) != len(other.Array) {
			return false
		}

		for i1 := range resource.Array {
		if resource.Array[i1] != other.Array[i1] {
			return false
		}
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource DefaultsStructComplexField) Validate() error {
	var errs cog.BuildErrors
		if err := resource.Nested.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("nested", err)...)
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type DefaultsStructPartialComplexField struct {
	Uid string `json:"uid"`
	IntVal int64 `json:"intVal"`
}

func (resource *DefaultsStructPartialComplexField) StrictUnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "uid"
	if fields["uid"] != nil {
		if string(fields["uid"]) != "null" {
			if err := json.Unmarshal(fields["uid"], &resource.Uid); err != nil {
				errs = append(errs, cog.MakeBuildErrors("uid", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("uid", errors.New("required field is null"))...)
		
		}
		delete(fields, "uid")
	} else {errs = append(errs, cog.MakeBuildErrors("uid", errors.New("required field is missing from input"))...)
	}
	// Field "intVal"
	if fields["intVal"] != nil {
		if string(fields["intVal"]) != "null" {
			if err := json.Unmarshal(fields["intVal"], &resource.IntVal); err != nil {
				errs = append(errs, cog.MakeBuildErrors("intVal", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("intVal", errors.New("required field is null"))...)
		
		}
		delete(fields, "intVal")
	} else {errs = append(errs, cog.MakeBuildErrors("intVal", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("DefaultsStructPartialComplexField", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource DefaultsStructPartialComplexField) Equals(other DefaultsStructPartialComplexField) bool {
		if resource.Uid != other.Uid {
			return false
		}
		if resource.IntVal != other.IntVal {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource DefaultsStructPartialComplexField) Validate() error {
	return nil
}


