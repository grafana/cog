package constant_references

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type Enum string
const (
	EnumValueA Enum = "ValueA"
	EnumValueB Enum = "ValueB"
	EnumValueC Enum = "ValueC"
)


type ParentStruct struct {
    MyEnum Enum `json:"myEnum"`
}

// NewParentStruct creates a new ParentStruct object.
func NewParentStruct() *ParentStruct {
	return &ParentStruct{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `ParentStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *ParentStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "myEnum"
	if fields["myEnum"] != nil {
		if string(fields["myEnum"]) != "null" {
			if err := json.Unmarshal(fields["myEnum"], &resource.MyEnum); err != nil {
				errs = append(errs, cog.MakeBuildErrors("myEnum", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is null"))...)
		
		}
		delete(fields, "myEnum")
	} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("ParentStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `ParentStruct` objects.
func (resource ParentStruct) Equals(other ParentStruct) bool {
		if resource.MyEnum != other.MyEnum {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `ParentStruct` fields for violations and returns them.
func (resource ParentStruct) Validate() error {
	return nil
}


type Struct struct {
    MyValue string `json:"myValue"`
    MyEnum Enum `json:"myEnum"`
}

// NewStruct creates a new Struct object.
func NewStruct() *Struct {
	return &Struct{
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
	// Field "myValue"
	if fields["myValue"] != nil {
		if string(fields["myValue"]) != "null" {
			if err := json.Unmarshal(fields["myValue"], &resource.MyValue); err != nil {
				errs = append(errs, cog.MakeBuildErrors("myValue", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("myValue", errors.New("required field is null"))...)
		
		}
		delete(fields, "myValue")
	} else {errs = append(errs, cog.MakeBuildErrors("myValue", errors.New("required field is missing from input"))...)
	}
	// Field "myEnum"
	if fields["myEnum"] != nil {
		if string(fields["myEnum"]) != "null" {
			if err := json.Unmarshal(fields["myEnum"], &resource.MyEnum); err != nil {
				errs = append(errs, cog.MakeBuildErrors("myEnum", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is null"))...)
		
		}
		delete(fields, "myEnum")
	} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is missing from input"))...)
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
		if resource.MyValue != other.MyValue {
			return false
		}
		if resource.MyEnum != other.MyEnum {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Struct` fields for violations and returns them.
func (resource Struct) Validate() error {
	return nil
}


type StructA struct {
    MyEnum Enum `json:"myEnum"`
}

// NewStructA creates a new StructA object.
func NewStructA() *StructA {
	return &StructA{
		MyEnum: EnumValueA,
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StructA` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StructA) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "myEnum"
	if fields["myEnum"] != nil {
		if string(fields["myEnum"]) != "null" {
			if err := json.Unmarshal(fields["myEnum"], &resource.MyEnum); err != nil {
				errs = append(errs, cog.MakeBuildErrors("myEnum", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is null"))...)
		
		}
		delete(fields, "myEnum")
	} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("StructA", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `StructA` objects.
func (resource StructA) Equals(other StructA) bool {
        if resource.MyEnum != other.MyEnum {
            return false
        }

	return true
}


// Validate checks all the validation constraints that may be defined on `StructA` fields for violations and returns them.
func (resource StructA) Validate() error {
	var errs cog.BuildErrors
    if resource.MyEnum != "ValueA" {
        errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("must be ValueA"))...)
    }

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type StructB struct {
    MyEnum Enum `json:"myEnum"`
    MyValue string `json:"myValue"`
}

// NewStructB creates a new StructB object.
func NewStructB() *StructB {
	return &StructB{
		MyEnum: EnumValueB,
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StructB` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StructB) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "myEnum"
	if fields["myEnum"] != nil {
		if string(fields["myEnum"]) != "null" {
			if err := json.Unmarshal(fields["myEnum"], &resource.MyEnum); err != nil {
				errs = append(errs, cog.MakeBuildErrors("myEnum", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is null"))...)
		
		}
		delete(fields, "myEnum")
	} else {errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("required field is missing from input"))...)
	}
	// Field "myValue"
	if fields["myValue"] != nil {
		if string(fields["myValue"]) != "null" {
			if err := json.Unmarshal(fields["myValue"], &resource.MyValue); err != nil {
				errs = append(errs, cog.MakeBuildErrors("myValue", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("myValue", errors.New("required field is null"))...)
		
		}
		delete(fields, "myValue")
	} else {errs = append(errs, cog.MakeBuildErrors("myValue", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("StructB", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `StructB` objects.
func (resource StructB) Equals(other StructB) bool {
        if resource.MyEnum != other.MyEnum {
            return false
        }
		if resource.MyValue != other.MyValue {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StructB` fields for violations and returns them.
func (resource StructB) Validate() error {
	var errs cog.BuildErrors
    if resource.MyEnum != "ValueB" {
        errs = append(errs, cog.MakeBuildErrors("myEnum", errors.New("must be ValueB"))...)
    }

	if len(errs) == 0 {
		return nil
	}

	return errs
}


