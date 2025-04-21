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

// NewNestedStruct creates a new NestedStruct object.
func NewNestedStruct() *NestedStruct {
	return &NestedStruct{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `NestedStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *NestedStruct) UnmarshalJSONStrict(raw []byte) error {
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


// Equals tests the equality of two `NestedStruct` objects.
func (resource NestedStruct) Equals(other NestedStruct) bool {
		if resource.StringVal != other.StringVal {
			return false
		}
		if resource.IntVal != other.IntVal {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `NestedStruct` fields for violations and returns them.
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

// NewStruct creates a new Struct object.
func NewStruct() *Struct {
	return &Struct{
		AllFields: NestedStruct{
		StringVal: "hello",
		IntVal: 3,
},
		PartialFields: NestedStruct{
		IntVal: 3,
},
		EmptyFields: *NewNestedStruct(),
		ComplexField: DefaultsStructComplexField{
		Uid: "myUID",
		Nested: map[string]interface {}{"nestedVal":"nested"},
		Array: []string{"hello"},
},
		PartialComplexField: DefaultsStructPartialComplexField{
},
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Struct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Struct) UnmarshalJSONStrict(raw []byte) error {
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
			if err := resource.AllFields.UnmarshalJSONStrict(fields["allFields"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("allFields", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("allFields", errors.New("required field is null"))...)
		
		}
		delete(fields, "allFields")
	
	}
	// Field "partialFields"
	if fields["partialFields"] != nil {
		if string(fields["partialFields"]) != "null" {
			
			resource.PartialFields = NestedStruct{}
			if err := resource.PartialFields.UnmarshalJSONStrict(fields["partialFields"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("partialFields", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("partialFields", errors.New("required field is null"))...)
		
		}
		delete(fields, "partialFields")
	
	}
	// Field "emptyFields"
	if fields["emptyFields"] != nil {
		if string(fields["emptyFields"]) != "null" {
			
			resource.EmptyFields = NestedStruct{}
			if err := resource.EmptyFields.UnmarshalJSONStrict(fields["emptyFields"]); err != nil {
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
			if err := resource.ComplexField.UnmarshalJSONStrict(fields["complexField"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("complexField", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("complexField", errors.New("required field is null"))...)
		
		}
		delete(fields, "complexField")
	
	}
	// Field "partialComplexField"
	if fields["partialComplexField"] != nil {
		if string(fields["partialComplexField"]) != "null" {
			
			resource.PartialComplexField = DefaultsStructPartialComplexField{}
			if err := resource.PartialComplexField.UnmarshalJSONStrict(fields["partialComplexField"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("partialComplexField", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("partialComplexField", errors.New("required field is null"))...)
		
		}
		delete(fields, "partialComplexField")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Struct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `Struct` objects.
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


// Validate checks all the validation constraints that may be defined on `Struct` fields for violations and returns them.
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

// NewDefaultsStructComplexFieldNested creates a new DefaultsStructComplexFieldNested object.
func NewDefaultsStructComplexFieldNested() *DefaultsStructComplexFieldNested {
	return &DefaultsStructComplexFieldNested{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `DefaultsStructComplexFieldNested` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *DefaultsStructComplexFieldNested) UnmarshalJSONStrict(raw []byte) error {
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


// Equals tests the equality of two `DefaultsStructComplexFieldNested` objects.
func (resource DefaultsStructComplexFieldNested) Equals(other DefaultsStructComplexFieldNested) bool {
		if resource.NestedVal != other.NestedVal {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `DefaultsStructComplexFieldNested` fields for violations and returns them.
func (resource DefaultsStructComplexFieldNested) Validate() error {
	return nil
}


type DefaultsStructComplexField struct {
    Uid string `json:"uid"`
    Nested DefaultsStructComplexFieldNested `json:"nested"`
    Array []string `json:"array"`
}

// NewDefaultsStructComplexField creates a new DefaultsStructComplexField object.
func NewDefaultsStructComplexField() *DefaultsStructComplexField {
	return &DefaultsStructComplexField{
		Nested: *NewDefaultsStructComplexFieldNested(),
		Array: []string{},
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `DefaultsStructComplexField` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *DefaultsStructComplexField) UnmarshalJSONStrict(raw []byte) error {
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
			if err := resource.Nested.UnmarshalJSONStrict(fields["nested"]); err != nil {
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


// Equals tests the equality of two `DefaultsStructComplexField` objects.
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


// Validate checks all the validation constraints that may be defined on `DefaultsStructComplexField` fields for violations and returns them.
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

// NewDefaultsStructPartialComplexField creates a new DefaultsStructPartialComplexField object.
func NewDefaultsStructPartialComplexField() *DefaultsStructPartialComplexField {
	return &DefaultsStructPartialComplexField{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `DefaultsStructPartialComplexField` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *DefaultsStructPartialComplexField) UnmarshalJSONStrict(raw []byte) error {
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


// Equals tests the equality of two `DefaultsStructPartialComplexField` objects.
func (resource DefaultsStructPartialComplexField) Equals(other DefaultsStructPartialComplexField) bool {
		if resource.Uid != other.Uid {
			return false
		}
		if resource.IntVal != other.IntVal {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `DefaultsStructPartialComplexField` fields for violations and returns them.
func (resource DefaultsStructPartialComplexField) Validate() error {
	return nil
}


