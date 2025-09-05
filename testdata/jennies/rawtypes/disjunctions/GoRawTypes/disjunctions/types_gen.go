package disjunctions

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"reflect"
	"bytes"
)

// Refresh rate or disabled.
type RefreshRate = StringOrBool

// NewRefreshRate creates a new RefreshRate object.
func NewRefreshRate() *RefreshRate {
	return NewStringOrBool()
}
type StringOrNull *string

type SomeStruct struct {
    Type string `json:"Type"`
    FieldAny any `json:"FieldAny"`
}

// NewSomeStruct creates a new SomeStruct object.
func NewSomeStruct() *SomeStruct {
	return &SomeStruct{
		Type: "some-struct",
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


// Equals tests the equality of two `SomeStruct` objects.
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


// Validate checks all the validation constraints that may be defined on `SomeStruct` fields for violations and returns them.
func (resource SomeStruct) Validate() error {
	return nil
}


type BoolOrRef = BoolOrSomeStruct

// NewBoolOrRef creates a new BoolOrRef object.
func NewBoolOrRef() *BoolOrRef {
	return NewBoolOrSomeStruct()
}
type SomeOtherStruct struct {
    Type string `json:"Type"`
    Foo []byte `json:"Foo"`
}

// NewSomeOtherStruct creates a new SomeOtherStruct object.
func NewSomeOtherStruct() *SomeOtherStruct {
	return &SomeOtherStruct{
		Type: "some-other-struct",
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


// Equals tests the equality of two `SomeOtherStruct` objects.
func (resource SomeOtherStruct) Equals(other SomeOtherStruct) bool {
		if resource.Type != other.Type {
			return false
		}
	    if !bytes.Equal(resource.Foo, other.Foo) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `SomeOtherStruct` fields for violations and returns them.
func (resource SomeOtherStruct) Validate() error {
	return nil
}


type YetAnotherStruct struct {
    Type string `json:"Type"`
    Bar uint8 `json:"Bar"`
}

// NewYetAnotherStruct creates a new YetAnotherStruct object.
func NewYetAnotherStruct() *YetAnotherStruct {
	return &YetAnotherStruct{
		Type: "yet-another-struct",
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `YetAnotherStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
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


// Equals tests the equality of two `YetAnotherStruct` objects.
func (resource YetAnotherStruct) Equals(other YetAnotherStruct) bool {
		if resource.Type != other.Type {
			return false
		}
		if resource.Bar != other.Bar {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `YetAnotherStruct` fields for violations and returns them.
func (resource YetAnotherStruct) Validate() error {
	return nil
}


type SeveralRefs = SomeStructOrSomeOtherStructOrYetAnotherStruct

// NewSeveralRefs creates a new SeveralRefs object.
func NewSeveralRefs() *SeveralRefs {
	return NewSomeStructOrSomeOtherStructOrYetAnotherStruct()
}
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


type BoolOrSomeStruct struct {
    Bool *bool `json:"Bool,omitempty"`
    SomeStruct *SomeStruct `json:"SomeStruct,omitempty"`
}

// NewBoolOrSomeStruct creates a new BoolOrSomeStruct object.
func NewBoolOrSomeStruct() *BoolOrSomeStruct {
	return &BoolOrSomeStruct{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `BoolOrSomeStruct` as JSON.
func (resource BoolOrSomeStruct) MarshalJSON() ([]byte, error) {
	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}
	if resource.SomeStruct != nil {
		return json.Marshal(resource.SomeStruct)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `BoolOrSomeStruct` from JSON.
func (resource *BoolOrSomeStruct) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// Bool
	var Bool bool
	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
		resource.Bool = nil
	} else {
		resource.Bool = &Bool
		return nil
	}

	// SomeStruct
	var SomeStruct SomeStruct
    someStructdec := json.NewDecoder(bytes.NewReader(raw))
    someStructdec.DisallowUnknownFields()
    if err := someStructdec.Decode(&SomeStruct); err != nil {
        errList = append(errList, err)
        resource.SomeStruct = nil
    } else {
        resource.SomeStruct = &SomeStruct
        return nil
    }

	return errors.Join(errList...)
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `BoolOrSomeStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *BoolOrSomeStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// Bool
	var Bool bool
	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
	} else {
		resource.Bool = &Bool
		return nil
	}

	// SomeStruct
	var SomeStruct SomeStruct
    someStructdec := json.NewDecoder(bytes.NewReader(raw))
    someStructdec.DisallowUnknownFields()
    if err := someStructdec.Decode(&SomeStruct); err != nil {
        errList = append(errList, err)
    } else {
        resource.SomeStruct = &SomeStruct
        return nil
    }

	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("BoolOrSomeStruct", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `BoolOrSomeStruct` objects.
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


// Validate checks all the validation constraints that may be defined on `BoolOrSomeStruct` fields for violations and returns them.
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

// NewSomeStructOrSomeOtherStructOrYetAnotherStruct creates a new SomeStructOrSomeOtherStructOrYetAnotherStruct object.
func NewSomeStructOrSomeOtherStructOrYetAnotherStruct() *SomeStructOrSomeOtherStructOrYetAnotherStruct {
	return &SomeStructOrSomeOtherStructOrYetAnotherStruct{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `SomeStructOrSomeOtherStructOrYetAnotherStruct` as JSON.
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

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `SomeStructOrSomeOtherStructOrYetAnotherStruct` from JSON.
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
		return nil
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

	return nil
}


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `SomeStructOrSomeOtherStructOrYetAnotherStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
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

// Equals tests the equality of two `SomeStructOrSomeOtherStructOrYetAnotherStruct` objects.
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


// Validate checks all the validation constraints that may be defined on `SomeStructOrSomeOtherStructOrYetAnotherStruct` fields for violations and returns them.
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


