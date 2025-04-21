// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
//     GoRawTypes

package defaults

import (
	"encoding/json"
	"errors"
	"fmt"

	cog "github.com/grafana/cog/testdata/generated/cog"
)

type VariableOption struct {
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	Selected *BoolOrString/* DisjunctionToType[disjunction → ref] */ `json:"selected,omitempty"`
	Text     StringOrArrayOfString/* DisjunctionToType[disjunction → ref] */ `json:"text"`
	Value    StringOrArrayOfString/* DisjunctionToType[disjunction → ref] */ `json:"value"`
}

// NewVariableOption creates a new VariableOption object.
func NewVariableOption() *VariableOption {
	return &VariableOption{
		Text:  *NewStringOrArrayOfString(),
		Value: *NewStringOrArrayOfString(),
	}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `VariableOption` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *VariableOption) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "selected"
	if fields["selected"] != nil {
		if string(fields["selected"]) != "null" {

			resource.Selected = &BoolOrString{}
			if err := resource.Selected.UnmarshalJSONStrict(fields["selected"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("selected", err)...)
			}

		}
		delete(fields, "selected")

	}
	// Field "text"
	if fields["text"] != nil {
		if string(fields["text"]) != "null" {

			resource.Text = StringOrArrayOfString{}
			if err := resource.Text.UnmarshalJSONStrict(fields["text"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("text", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("text", errors.New("required field is null"))...)

		}
		delete(fields, "text")
	} else {
		errs = append(errs, cog.MakeBuildErrors("text", errors.New("required field is missing from input"))...)
	}
	// Field "value"
	if fields["value"] != nil {
		if string(fields["value"]) != "null" {

			resource.Value = StringOrArrayOfString{}
			if err := resource.Value.UnmarshalJSONStrict(fields["value"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("value", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("value", errors.New("required field is null"))...)

		}
		delete(fields, "value")
	} else {
		errs = append(errs, cog.MakeBuildErrors("value", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("VariableOption", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `VariableOption` objects.
func (resource VariableOption) Equals(other VariableOption) bool {
	if resource.Selected == nil && other.Selected != nil || resource.Selected != nil && other.Selected == nil {
		return false
	}

	if resource.Selected != nil {
		if !resource.Selected.Equals(*other.Selected) {
			return false
		}
	}
	if !resource.Text.Equals(other.Text) {
		return false
	}
	if !resource.Value.Equals(other.Value) {
		return false
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `VariableOption` fields for violations and returns them.
func (resource VariableOption) Validate() error {
	var errs cog.BuildErrors
	if resource.Selected != nil {
		if err := resource.Selected.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("selected", err)...)
		}
	}
	if err := resource.Text.Validate(); err != nil {
		errs = append(errs, cog.MakeBuildErrors("text", err)...)
	}
	if err := resource.Value.Validate(); err != nil {
		errs = append(errs, cog.MakeBuildErrors("value", err)...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type TextVariable struct {
	Name        string         `json:"name"`
	Current     VariableOption `json:"current"`
	SkipUrlSync bool           `json:"skipUrlSync"`
}

// NewTextVariable creates a new TextVariable object.
func NewTextVariable() *TextVariable {
	return &TextVariable{
		Name: "",
		Current: VariableOption{
			Selected: &BoolOrString{
				String: (func(input string) *string { return &input })("maybe"),
			},
			Text: StringOrArrayOfString{
				String: (func(input string) *string { return &input })(""),
			},
			Value: StringOrArrayOfString{
				ArrayOfString: []string{"val"},
			},
		},
		SkipUrlSync: false,
	}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `TextVariable` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *TextVariable) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "name"
	if fields["name"] != nil {
		if string(fields["name"]) != "null" {
			if err := json.Unmarshal(fields["name"], &resource.Name); err != nil {
				errs = append(errs, cog.MakeBuildErrors("name", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("name", errors.New("required field is null"))...)

		}
		delete(fields, "name")

	}
	// Field "current"
	if fields["current"] != nil {
		if string(fields["current"]) != "null" {

			resource.Current = VariableOption{}
			if err := resource.Current.UnmarshalJSONStrict(fields["current"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("current", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("current", errors.New("required field is null"))...)

		}
		delete(fields, "current")

	}
	// Field "skipUrlSync"
	if fields["skipUrlSync"] != nil {
		if string(fields["skipUrlSync"]) != "null" {
			if err := json.Unmarshal(fields["skipUrlSync"], &resource.SkipUrlSync); err != nil {
				errs = append(errs, cog.MakeBuildErrors("skipUrlSync", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("skipUrlSync", errors.New("required field is null"))...)

		}
		delete(fields, "skipUrlSync")

	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("TextVariable", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `TextVariable` objects.
func (resource TextVariable) Equals(other TextVariable) bool {
	if resource.Name != other.Name {
		return false
	}
	if !resource.Current.Equals(other.Current) {
		return false
	}
	if resource.SkipUrlSync != other.SkipUrlSync {
		return false
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `TextVariable` fields for violations and returns them.
func (resource TextVariable) Validate() error {
	var errs cog.BuildErrors
	if err := resource.Current.Validate(); err != nil {
		errs = append(errs, cog.MakeBuildErrors("current", err)...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Modified by compiler pass 'DisjunctionToType[created]'
type BoolOrString struct {
	Bool   *bool   `json:"Bool,omitempty"`
	String *string `json:"String,omitempty"`
}

// NewBoolOrString creates a new BoolOrString object.
func NewBoolOrString() *BoolOrString {
	return &BoolOrString{}
}

// MarshalJSON implements a custom JSON marshalling logic to encode `BoolOrString` as JSON.
func (resource BoolOrString) MarshalJSON() ([]byte, error) {
	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}

	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `BoolOrString` from JSON.
func (resource *BoolOrString) UnmarshalJSON(raw []byte) error {
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

	// String
	var String string
	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
		resource.String = nil
	} else {
		resource.String = &String
		return nil
	}

	return errors.Join(errList...)
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `BoolOrString` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *BoolOrString) UnmarshalJSONStrict(raw []byte) error {
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

	// String
	var String string

	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
	} else {
		resource.String = &String
		return nil
	}

	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("BoolOrString", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `BoolOrString` objects.
func (resource BoolOrString) Equals(other BoolOrString) bool {
	if resource.Bool == nil && other.Bool != nil || resource.Bool != nil && other.Bool == nil {
		return false
	}

	if resource.Bool != nil {
		if *resource.Bool != *other.Bool {
			return false
		}
	}
	if resource.String == nil && other.String != nil || resource.String != nil && other.String == nil {
		return false
	}

	if resource.String != nil {
		if *resource.String != *other.String {
			return false
		}
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `BoolOrString` fields for violations and returns them.
func (resource BoolOrString) Validate() error {
	return nil
}

// Modified by compiler pass 'DisjunctionToType[created]'
type StringOrArrayOfString struct {
	String        *string  `json:"String,omitempty"`
	ArrayOfString []string `json:"ArrayOfString,omitempty"`
}

// NewStringOrArrayOfString creates a new StringOrArrayOfString object.
func NewStringOrArrayOfString() *StringOrArrayOfString {
	return &StringOrArrayOfString{}
}

// MarshalJSON implements a custom JSON marshalling logic to encode `StringOrArrayOfString` as JSON.
func (resource StringOrArrayOfString) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.ArrayOfString != nil {
		return json.Marshal(resource.ArrayOfString)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StringOrArrayOfString` from JSON.
func (resource *StringOrArrayOfString) UnmarshalJSON(raw []byte) error {
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

	// ArrayOfString
	var ArrayOfString []string
	if err := json.Unmarshal(raw, &ArrayOfString); err != nil {
		errList = append(errList, err)
		resource.ArrayOfString = nil
	} else {
		resource.ArrayOfString = ArrayOfString
		return nil
	}

	return errors.Join(errList...)
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StringOrArrayOfString` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StringOrArrayOfString) UnmarshalJSONStrict(raw []byte) error {
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

	// ArrayOfString
	var ArrayOfString []string

	if err := json.Unmarshal(raw, &ArrayOfString); err != nil {
		errList = append(errList, err)
	} else {
		resource.ArrayOfString = ArrayOfString
		return nil
	}

	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StringOrArrayOfString", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StringOrArrayOfString` objects.
func (resource StringOrArrayOfString) Equals(other StringOrArrayOfString) bool {
	if resource.String == nil && other.String != nil || resource.String != nil && other.String == nil {
		return false
	}

	if resource.String != nil {
		if *resource.String != *other.String {
			return false
		}
	}

	if len(resource.ArrayOfString) != len(other.ArrayOfString) {
		return false
	}

	for i1 := range resource.ArrayOfString {
		if resource.ArrayOfString[i1] != other.ArrayOfString[i1] {
			return false
		}
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `StringOrArrayOfString` fields for violations and returns them.
func (resource StringOrArrayOfString) Validate() error {
	return nil
}
