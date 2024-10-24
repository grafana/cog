package disjunctions

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"reflect"
)

// Refresh rate or disabled.
type RefreshRate = StringOrBool

type StringOrNull *string

type SomeStruct struct {
	Type string `json:"Type"`
	FieldAny any `json:"FieldAny"`
}

func (resource *SomeStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "Type"
	if fields["Type"] != nil {
		if string(fields["Type"]) != "null" {
			if err := json.Unmarshal(fields["Type"], &resource.Type); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Type", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("Type", errors.New("required field is null"))...)
		
		}
		delete(fields, "Type")
	} else {errs = append(errs, cog.MakeBuildErrors("Type", errors.New("required field is missing from input"))...)
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
		errs = append(errs, cog.MakeBuildErrors("SomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource SomeStruct) Equals(other SomeStruct) bool {
		if resource.Type != other.Type {
			return false
		}
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeStruct) Validate() error {
	return nil
}


type BoolOrRef = BoolOrSomeStruct

type SomeOtherStruct struct {
	Type string `json:"Type"`
	Foo []byte `json:"Foo"`
}

func (resource *SomeOtherStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "Type"
	if fields["Type"] != nil {
		if string(fields["Type"]) != "null" {
			if err := json.Unmarshal(fields["Type"], &resource.Type); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Type", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("Type", errors.New("required field is null"))...)
		
		}
		delete(fields, "Type")
	} else {errs = append(errs, cog.MakeBuildErrors("Type", errors.New("required field is missing from input"))...)
	}
	// Field "Foo"
	if fields["Foo"] != nil {
		if string(fields["Foo"]) != "null" {
			if err := json.Unmarshal(fields["Foo"], &resource.Foo); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Foo", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("Foo", errors.New("required field is null"))...)
		
		}
		delete(fields, "Foo")
	} else {errs = append(errs, cog.MakeBuildErrors("Foo", errors.New("required field is missing from input"))...)
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
		if resource.Type != other.Type {
			return false
		}
		if resource.Foo != other.Foo {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeOtherStruct) Validate() error {
	return nil
}


type YetAnotherStruct struct {
	Type string `json:"Type"`
	Bar uint8 `json:"Bar"`
}

func (resource *YetAnotherStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "Type"
	if fields["Type"] != nil {
		if string(fields["Type"]) != "null" {
			if err := json.Unmarshal(fields["Type"], &resource.Type); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Type", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("Type", errors.New("required field is null"))...)
		
		}
		delete(fields, "Type")
	} else {errs = append(errs, cog.MakeBuildErrors("Type", errors.New("required field is missing from input"))...)
	}
	// Field "Bar"
	if fields["Bar"] != nil {
		if string(fields["Bar"]) != "null" {
			if err := json.Unmarshal(fields["Bar"], &resource.Bar); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Bar", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("Bar", errors.New("required field is null"))...)
		
		}
		delete(fields, "Bar")
	} else {errs = append(errs, cog.MakeBuildErrors("Bar", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("YetAnotherStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource YetAnotherStruct) Equals(other YetAnotherStruct) bool {
		if resource.Type != other.Type {
			return false
		}
		if resource.Bar != other.Bar {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource YetAnotherStruct) Validate() error {
	return nil
}


type SeveralRefs = SomeStructOrSomeOtherStructOrYetAnotherStruct

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


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource StringOrBool) Validate() error {
	return nil
}


type BoolOrSomeStruct struct {
	Bool *bool `json:"Bool,omitempty"`
	SomeStruct *SomeStruct `json:"SomeStruct,omitempty"`
}

func (resource *BoolOrSomeStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "Bool"
	if fields["Bool"] != nil {
		if string(fields["Bool"]) != "null" {
			if err := json.Unmarshal(fields["Bool"], &resource.Bool); err != nil {
				errs = append(errs, cog.MakeBuildErrors("Bool", err)...)
			}
		
		}
		delete(fields, "Bool")
	
	}
	// Field "SomeStruct"
	if fields["SomeStruct"] != nil {
		if string(fields["SomeStruct"]) != "null" {
			
			resource.SomeStruct = &SomeStruct{}
			if err := resource.SomeStruct.UnmarshalJSONStrict(fields["SomeStruct"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("SomeStruct", err)...)
			}
		
		}
		delete(fields, "SomeStruct")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("BoolOrSomeStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


func (resource BoolOrSomeStruct) Equals(other BoolOrSomeStruct) bool {
		if resource.Bool == nil && other.Bool != nil || resource.Bool != nil && other.Bool == nil {
			return false
		}

		if resource.Bool != nil {
		if *resource.Bool != *other.Bool {
			return false
		}
		}
		if resource.SomeStruct == nil && other.SomeStruct != nil || resource.SomeStruct != nil && other.SomeStruct == nil {
			return false
		}

		if resource.SomeStruct != nil {
		if !resource.SomeStruct.Equals(*other.SomeStruct) {
			return false
		}
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource BoolOrSomeStruct) Validate() error {
	var errs cog.BuildErrors
		if resource.SomeStruct != nil {
		if err := resource.SomeStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("SomeStruct", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type SomeStructOrSomeOtherStructOrYetAnotherStruct struct {
	SomeStruct *SomeStruct `json:"SomeStruct,omitempty"`
	SomeOtherStruct *SomeOtherStruct `json:"SomeOtherStruct,omitempty"`
	YetAnotherStruct *YetAnotherStruct `json:"YetAnotherStruct,omitempty"`
}

func (resource SomeStructOrSomeOtherStructOrYetAnotherStruct) MarshalJSON() ([]byte, error) {
	if resource.SomeStruct != nil {
		return json.Marshal(resource.SomeStruct)
	}
	if resource.SomeOtherStruct != nil {
		return json.Marshal(resource.SomeOtherStruct)
	}
	if resource.YetAnotherStruct != nil {
		return json.Marshal(resource.YetAnotherStruct)
	}

	return nil, fmt.Errorf("no value for disjunction of refs")
}

func (resource *SomeStructOrSomeOtherStructOrYetAnotherStruct) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["Type"]
	if !found {
		return errors.New("discriminator field 'Type' not found in payload")
	}

	switch discriminator {
	case "some-other-struct":
		var someOtherStruct SomeOtherStruct
		if err := json.Unmarshal(raw, &someOtherStruct); err != nil {
			return err
		}

		resource.SomeOtherStruct = &someOtherStruct
		return nil
	case "some-struct":
		var someStruct SomeStruct
		if err := json.Unmarshal(raw, &someStruct); err != nil {
			return err
		}

		resource.SomeStruct = &someStruct
		return nil
	case "yet-another-struct":
		var yetAnotherStruct YetAnotherStruct
		if err := json.Unmarshal(raw, &yetAnotherStruct); err != nil {
			return err
		}

		resource.YetAnotherStruct = &yetAnotherStruct
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `Type = %v`", discriminator)
}


func (resource *SomeStructOrSomeOtherStructOrYetAnotherStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["Type"]
	if !found {
		return fmt.Errorf("discriminator field 'Type' not found in payload")
	}

	switch discriminator {
		case "some-other-struct":
		someOtherStruct := &SomeOtherStruct{}
		if err := someOtherStruct.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.SomeOtherStruct = someOtherStruct
		return nil
		case "some-struct":
		someStruct := &SomeStruct{}
		if err := someStruct.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.SomeStruct = someStruct
		return nil
		case "yet-another-struct":
		yetAnotherStruct := &YetAnotherStruct{}
		if err := yetAnotherStruct.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.YetAnotherStruct = yetAnotherStruct
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `Type = %v`", discriminator)
}

func (resource SomeStructOrSomeOtherStructOrYetAnotherStruct) Equals(other SomeStructOrSomeOtherStructOrYetAnotherStruct) bool {
		if resource.SomeStruct == nil && other.SomeStruct != nil || resource.SomeStruct != nil && other.SomeStruct == nil {
			return false
		}

		if resource.SomeStruct != nil {
		if !resource.SomeStruct.Equals(*other.SomeStruct) {
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
		if resource.YetAnotherStruct == nil && other.YetAnotherStruct != nil || resource.YetAnotherStruct != nil && other.YetAnotherStruct == nil {
			return false
		}

		if resource.YetAnotherStruct != nil {
		if !resource.YetAnotherStruct.Equals(*other.YetAnotherStruct) {
			return false
		}
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource SomeStructOrSomeOtherStructOrYetAnotherStruct) Validate() error {
	var errs cog.BuildErrors
		if resource.SomeStruct != nil {
		if err := resource.SomeStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("SomeStruct", err)...)
		}
		}
		if resource.SomeOtherStruct != nil {
		if err := resource.SomeOtherStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("SomeOtherStruct", err)...)
		}
		}
		if resource.YetAnotherStruct != nil {
		if err := resource.YetAnotherStruct.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("YetAnotherStruct", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


