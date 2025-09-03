// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
//     GoRawTypes

package equality

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	cog "github.com/grafana/cog/testdata/generated/cog"
)

// Modified by compiler pass 'PrefixEnumValues'
type Direction string

const (
	DirectionTop    Direction = "top"
	DirectionBottom Direction = "bottom"
	DirectionLeft   Direction = "left"
	DirectionRight  Direction = "right"
)

type Variable struct {
	Name string `json:"name"`
}

// NewVariable creates a new Variable object.
func NewVariable() *Variable {
	return &Variable{}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Variable` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Variable) UnmarshalJSONStrict(raw []byte) error {
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
	} else {
		errs = append(errs, cog.MakeBuildErrors("name", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Variable", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `Variable` objects.
func (resource Variable) Equals(other Variable) bool {
	if resource.Name != other.Name {
		return false
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `Variable` fields for violations and returns them.
func (resource Variable) Validate() error {
	return nil
}

type Container struct {
	StringField string    `json:"stringField"`
	IntField    int64     `json:"intField"`
	EnumField   Direction `json:"enumField"`
	RefField    Variable  `json:"refField"`
}

// NewContainer creates a new Container object.
func NewContainer() *Container {
	return &Container{
		RefField: *NewVariable(),
	}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Container` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Container) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "stringField"
	if fields["stringField"] != nil {
		if string(fields["stringField"]) != "null" {
			if err := json.Unmarshal(fields["stringField"], &resource.StringField); err != nil {
				errs = append(errs, cog.MakeBuildErrors("stringField", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("stringField", errors.New("required field is null"))...)

		}
		delete(fields, "stringField")
	} else {
		errs = append(errs, cog.MakeBuildErrors("stringField", errors.New("required field is missing from input"))...)
	}
	// Field "intField"
	if fields["intField"] != nil {
		if string(fields["intField"]) != "null" {
			if err := json.Unmarshal(fields["intField"], &resource.IntField); err != nil {
				errs = append(errs, cog.MakeBuildErrors("intField", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("intField", errors.New("required field is null"))...)

		}
		delete(fields, "intField")
	} else {
		errs = append(errs, cog.MakeBuildErrors("intField", errors.New("required field is missing from input"))...)
	}
	// Field "enumField"
	if fields["enumField"] != nil {
		if string(fields["enumField"]) != "null" {
			if err := json.Unmarshal(fields["enumField"], &resource.EnumField); err != nil {
				errs = append(errs, cog.MakeBuildErrors("enumField", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("enumField", errors.New("required field is null"))...)

		}
		delete(fields, "enumField")
	} else {
		errs = append(errs, cog.MakeBuildErrors("enumField", errors.New("required field is missing from input"))...)
	}
	// Field "refField"
	if fields["refField"] != nil {
		if string(fields["refField"]) != "null" {

			resource.RefField = Variable{}
			if err := resource.RefField.UnmarshalJSONStrict(fields["refField"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("refField", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("refField", errors.New("required field is null"))...)

		}
		delete(fields, "refField")
	} else {
		errs = append(errs, cog.MakeBuildErrors("refField", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Container", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `Container` objects.
func (resource Container) Equals(other Container) bool {
	if resource.StringField != other.StringField {
		return false
	}
	if resource.IntField != other.IntField {
		return false
	}
	if resource.EnumField != other.EnumField {
		return false
	}
	if !resource.RefField.Equals(other.RefField) {
		return false
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `Container` fields for violations and returns them.
func (resource Container) Validate() error {
	var errs cog.BuildErrors
	if err := resource.RefField.Validate(); err != nil {
		errs = append(errs, cog.MakeBuildErrors("refField", err)...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type Optionals struct {
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	StringField *string `json:"stringField,omitempty"`
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	EnumField *Direction `json:"enumField,omitempty"`
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	RefField *Variable `json:"refField,omitempty"`
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	ByteField []byte `json:"byteField,omitempty"`
}

// NewOptionals creates a new Optionals object.
func NewOptionals() *Optionals {
	return &Optionals{}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Optionals` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Optionals) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "stringField"
	if fields["stringField"] != nil {
		if string(fields["stringField"]) != "null" {
			if err := json.Unmarshal(fields["stringField"], &resource.StringField); err != nil {
				errs = append(errs, cog.MakeBuildErrors("stringField", err)...)
			}

		}
		delete(fields, "stringField")

	}
	// Field "enumField"
	if fields["enumField"] != nil {
		if string(fields["enumField"]) != "null" {
			if err := json.Unmarshal(fields["enumField"], &resource.EnumField); err != nil {
				errs = append(errs, cog.MakeBuildErrors("enumField", err)...)
			}

		}
		delete(fields, "enumField")

	}
	// Field "refField"
	if fields["refField"] != nil {
		if string(fields["refField"]) != "null" {

			resource.RefField = &Variable{}
			if err := resource.RefField.UnmarshalJSONStrict(fields["refField"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("refField", err)...)
			}

		}
		delete(fields, "refField")

	}
	// Field "byteField"
	if fields["byteField"] != nil {
		if string(fields["byteField"]) != "null" {
			if err := json.Unmarshal(fields["byteField"], &resource.ByteField); err != nil {
				errs = append(errs, cog.MakeBuildErrors("byteField", err)...)
			}

		}
		delete(fields, "byteField")

	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Optionals", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `Optionals` objects.
func (resource Optionals) Equals(other Optionals) bool {
	if resource.StringField == nil && other.StringField != nil || resource.StringField != nil && other.StringField == nil {
		return false
	}

	if resource.StringField != nil {
		if *resource.StringField != *other.StringField {
			return false
		}
	}
	if resource.EnumField == nil && other.EnumField != nil || resource.EnumField != nil && other.EnumField == nil {
		return false
	}

	if resource.EnumField != nil {
		if *resource.EnumField != *other.EnumField {
			return false
		}
	}
	if resource.RefField == nil && other.RefField != nil || resource.RefField != nil && other.RefField == nil {
		return false
	}

	if resource.RefField != nil {
		if !resource.RefField.Equals(*other.RefField) {
			return false
		}
	}
	if resource.ByteField == nil && other.ByteField != nil || resource.ByteField != nil && other.ByteField == nil {
		return false
	}

	if resource.ByteField != nil {
		if !bytes.Equal(resource.ByteField, other.ByteField) {
			return false
		}
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `Optionals` fields for violations and returns them.
func (resource Optionals) Validate() error {
	var errs cog.BuildErrors
	if resource.RefField != nil {
		if err := resource.RefField.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("refField", err)...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type Arrays struct {
	Ints             []int64                          `json:"ints"`
	Strings          []string                         `json:"strings"`
	ArrayOfArray     [][]string                       `json:"arrayOfArray"`
	Refs             []Variable                       `json:"refs"`
	AnonymousStructs []EqualityArraysAnonymousStructs `json:"anonymousStructs"`
	ArrayOfAny       []any                            `json:"arrayOfAny"`
}

// NewArrays creates a new Arrays object.
func NewArrays() *Arrays {
	return &Arrays{
		Ints:             []int64{},
		Strings:          []string{},
		ArrayOfArray:     [][]string{},
		Refs:             []Variable{},
		AnonymousStructs: []EqualityArraysAnonymousStructs{},
		ArrayOfAny:       []any{},
	}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Arrays` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Arrays) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "ints"
	if fields["ints"] != nil {
		if string(fields["ints"]) != "null" {

			if err := json.Unmarshal(fields["ints"], &resource.Ints); err != nil {
				errs = append(errs, cog.MakeBuildErrors("ints", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("ints", errors.New("required field is null"))...)

		}
		delete(fields, "ints")
	} else {
		errs = append(errs, cog.MakeBuildErrors("ints", errors.New("required field is missing from input"))...)
	}
	// Field "strings"
	if fields["strings"] != nil {
		if string(fields["strings"]) != "null" {

			if err := json.Unmarshal(fields["strings"], &resource.Strings); err != nil {
				errs = append(errs, cog.MakeBuildErrors("strings", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("strings", errors.New("required field is null"))...)

		}
		delete(fields, "strings")
	} else {
		errs = append(errs, cog.MakeBuildErrors("strings", errors.New("required field is missing from input"))...)
	}
	// Field "arrayOfArray"
	if fields["arrayOfArray"] != nil {
		if string(fields["arrayOfArray"]) != "null" {

			if err := json.Unmarshal(fields["arrayOfArray"], &resource.ArrayOfArray); err != nil {
				errs = append(errs, cog.MakeBuildErrors("arrayOfArray", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("arrayOfArray", errors.New("required field is null"))...)

		}
		delete(fields, "arrayOfArray")
	} else {
		errs = append(errs, cog.MakeBuildErrors("arrayOfArray", errors.New("required field is missing from input"))...)
	}
	// Field "refs"
	if fields["refs"] != nil {
		if string(fields["refs"]) != "null" {

			partialArray := []json.RawMessage{}
			if err := json.Unmarshal(fields["refs"], &partialArray); err != nil {
				return err
			}

			for i1 := range partialArray {
				var result1 Variable

				result1 = Variable{}
				if err := result1.UnmarshalJSONStrict(partialArray[i1]); err != nil {
					errs = append(errs, cog.MakeBuildErrors("refs["+strconv.Itoa(i1)+"]", err)...)
				}
				resource.Refs = append(resource.Refs, result1)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("refs", errors.New("required field is null"))...)

		}
		delete(fields, "refs")
	} else {
		errs = append(errs, cog.MakeBuildErrors("refs", errors.New("required field is missing from input"))...)
	}
	// Field "anonymousStructs"
	if fields["anonymousStructs"] != nil {
		if string(fields["anonymousStructs"]) != "null" {

			partialArray := []json.RawMessage{}
			if err := json.Unmarshal(fields["anonymousStructs"], &partialArray); err != nil {
				return err
			}

			for i1 := range partialArray {
				var result1 EqualityArraysAnonymousStructs

				result1 = EqualityArraysAnonymousStructs{}
				if err := result1.UnmarshalJSONStrict(partialArray[i1]); err != nil {
					errs = append(errs, cog.MakeBuildErrors("anonymousStructs["+strconv.Itoa(i1)+"]", err)...)
				}
				resource.AnonymousStructs = append(resource.AnonymousStructs, result1)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("anonymousStructs", errors.New("required field is null"))...)

		}
		delete(fields, "anonymousStructs")
	} else {
		errs = append(errs, cog.MakeBuildErrors("anonymousStructs", errors.New("required field is missing from input"))...)
	}
	// Field "arrayOfAny"
	if fields["arrayOfAny"] != nil {
		if string(fields["arrayOfAny"]) != "null" {

			if err := json.Unmarshal(fields["arrayOfAny"], &resource.ArrayOfAny); err != nil {
				errs = append(errs, cog.MakeBuildErrors("arrayOfAny", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("arrayOfAny", errors.New("required field is null"))...)

		}
		delete(fields, "arrayOfAny")
	} else {
		errs = append(errs, cog.MakeBuildErrors("arrayOfAny", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Arrays", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `Arrays` objects.
func (resource Arrays) Equals(other Arrays) bool {

	if len(resource.Ints) != len(other.Ints) {
		return false
	}

	for i1 := range resource.Ints {
		if resource.Ints[i1] != other.Ints[i1] {
			return false
		}
	}

	if len(resource.Strings) != len(other.Strings) {
		return false
	}

	for i1 := range resource.Strings {
		if resource.Strings[i1] != other.Strings[i1] {
			return false
		}
	}

	if len(resource.ArrayOfArray) != len(other.ArrayOfArray) {
		return false
	}

	for i1 := range resource.ArrayOfArray {

		if len(resource.ArrayOfArray[i1]) != len(other.ArrayOfArray[i1]) {
			return false
		}

		for i2 := range resource.ArrayOfArray[i1] {
			if resource.ArrayOfArray[i1][i2] != other.ArrayOfArray[i1][i2] {
				return false
			}
		}
	}

	if len(resource.Refs) != len(other.Refs) {
		return false
	}

	for i1 := range resource.Refs {
		if !resource.Refs[i1].Equals(other.Refs[i1]) {
			return false
		}
	}

	if len(resource.AnonymousStructs) != len(other.AnonymousStructs) {
		return false
	}

	for i1 := range resource.AnonymousStructs {
		if !resource.AnonymousStructs[i1].Equals(other.AnonymousStructs[i1]) {
			return false
		}
	}

	if len(resource.ArrayOfAny) != len(other.ArrayOfAny) {
		return false
	}

	for i1 := range resource.ArrayOfAny {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.ArrayOfAny[i1], other.ArrayOfAny[i1]) {
			return false
		}
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `Arrays` fields for violations and returns them.
func (resource Arrays) Validate() error {
	var errs cog.BuildErrors

	for i1 := range resource.Refs {
		if err := resource.Refs[i1].Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("refs["+strconv.Itoa(i1)+"]", err)...)
		}
	}

	for i1 := range resource.AnonymousStructs {
		if err := resource.AnonymousStructs[i1].Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("anonymousStructs["+strconv.Itoa(i1)+"]", err)...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

type Maps struct {
	Ints             map[string]int64                        `json:"ints"`
	Strings          map[string]string                       `json:"strings"`
	Refs             map[string]Variable                     `json:"refs"`
	AnonymousStructs map[string]EqualityMapsAnonymousStructs `json:"anonymousStructs"`
	StringToAny      map[string]any                          `json:"stringToAny"`
}

// NewMaps creates a new Maps object.
func NewMaps() *Maps {
	return &Maps{
		Ints:             map[string]int64{},
		Strings:          map[string]string{},
		Refs:             map[string]Variable{},
		AnonymousStructs: map[string]EqualityMapsAnonymousStructs{},
		StringToAny:      map[string]any{},
	}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Maps` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Maps) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "ints"
	if fields["ints"] != nil {
		if string(fields["ints"]) != "null" {

			if err := json.Unmarshal(fields["ints"], &resource.Ints); err != nil {
				errs = append(errs, cog.MakeBuildErrors("ints", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("ints", errors.New("required field is null"))...)

		}
		delete(fields, "ints")
	} else {
		errs = append(errs, cog.MakeBuildErrors("ints", errors.New("required field is missing from input"))...)
	}
	// Field "strings"
	if fields["strings"] != nil {
		if string(fields["strings"]) != "null" {

			if err := json.Unmarshal(fields["strings"], &resource.Strings); err != nil {
				errs = append(errs, cog.MakeBuildErrors("strings", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("strings", errors.New("required field is null"))...)

		}
		delete(fields, "strings")
	} else {
		errs = append(errs, cog.MakeBuildErrors("strings", errors.New("required field is missing from input"))...)
	}
	// Field "refs"
	if fields["refs"] != nil {
		if string(fields["refs"]) != "null" {

			partialMap := make(map[string]json.RawMessage)
			if err := json.Unmarshal(fields["refs"], &partialMap); err != nil {
				return err
			}
			parsedMap1 := make(map[string]Variable, len(partialMap))
			for key1 := range partialMap {
				var result1 Variable

				result1 = Variable{}
				if err := result1.UnmarshalJSONStrict(partialMap[key1]); err != nil {
					errs = append(errs, cog.MakeBuildErrors("refs["+key1+"]", err)...)
				}
				parsedMap1[key1] = result1
			}
			resource.Refs = parsedMap1
		} else {
			errs = append(errs, cog.MakeBuildErrors("refs", errors.New("required field is null"))...)

		}
		delete(fields, "refs")
	} else {
		errs = append(errs, cog.MakeBuildErrors("refs", errors.New("required field is missing from input"))...)
	}
	// Field "anonymousStructs"
	if fields["anonymousStructs"] != nil {
		if string(fields["anonymousStructs"]) != "null" {

			partialMap := make(map[string]json.RawMessage)
			if err := json.Unmarshal(fields["anonymousStructs"], &partialMap); err != nil {
				return err
			}
			parsedMap1 := make(map[string]EqualityMapsAnonymousStructs, len(partialMap))
			for key1 := range partialMap {
				var result1 EqualityMapsAnonymousStructs

				result1 = EqualityMapsAnonymousStructs{}
				if err := result1.UnmarshalJSONStrict(partialMap[key1]); err != nil {
					errs = append(errs, cog.MakeBuildErrors("anonymousStructs["+key1+"]", err)...)
				}
				parsedMap1[key1] = result1
			}
			resource.AnonymousStructs = parsedMap1
		} else {
			errs = append(errs, cog.MakeBuildErrors("anonymousStructs", errors.New("required field is null"))...)

		}
		delete(fields, "anonymousStructs")
	} else {
		errs = append(errs, cog.MakeBuildErrors("anonymousStructs", errors.New("required field is missing from input"))...)
	}
	// Field "stringToAny"
	if fields["stringToAny"] != nil {
		if string(fields["stringToAny"]) != "null" {

			if err := json.Unmarshal(fields["stringToAny"], &resource.StringToAny); err != nil {
				errs = append(errs, cog.MakeBuildErrors("stringToAny", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("stringToAny", errors.New("required field is null"))...)

		}
		delete(fields, "stringToAny")
	} else {
		errs = append(errs, cog.MakeBuildErrors("stringToAny", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Maps", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `Maps` objects.
func (resource Maps) Equals(other Maps) bool {

	if len(resource.Ints) != len(other.Ints) {
		return false
	}

	for key1 := range resource.Ints {
		if resource.Ints[key1] != other.Ints[key1] {
			return false
		}
	}

	if len(resource.Strings) != len(other.Strings) {
		return false
	}

	for key1 := range resource.Strings {
		if resource.Strings[key1] != other.Strings[key1] {
			return false
		}
	}

	if len(resource.Refs) != len(other.Refs) {
		return false
	}

	for key1 := range resource.Refs {
		if !resource.Refs[key1].Equals(other.Refs[key1]) {
			return false
		}
	}

	if len(resource.AnonymousStructs) != len(other.AnonymousStructs) {
		return false
	}

	for key1 := range resource.AnonymousStructs {
		if !resource.AnonymousStructs[key1].Equals(other.AnonymousStructs[key1]) {
			return false
		}
	}

	if len(resource.StringToAny) != len(other.StringToAny) {
		return false
	}

	for key1 := range resource.StringToAny {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.StringToAny[key1], other.StringToAny[key1]) {
			return false
		}
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `Maps` fields for violations and returns them.
func (resource Maps) Validate() error {
	var errs cog.BuildErrors

	for key1 := range resource.Refs {
		if err := resource.Refs[key1].Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("refs["+key1+"]", err)...)
		}
	}

	for key1 := range resource.AnonymousStructs {
		if err := resource.AnonymousStructs[key1].Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("anonymousStructs["+key1+"]", err)...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Modified by compiler pass 'AnonymousStructsToNamed'
type EqualityArraysAnonymousStructs struct {
	Inner string `json:"inner"`
}

// NewEqualityArraysAnonymousStructs creates a new EqualityArraysAnonymousStructs object.
func NewEqualityArraysAnonymousStructs() *EqualityArraysAnonymousStructs {
	return &EqualityArraysAnonymousStructs{}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `EqualityArraysAnonymousStructs` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *EqualityArraysAnonymousStructs) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "inner"
	if fields["inner"] != nil {
		if string(fields["inner"]) != "null" {
			if err := json.Unmarshal(fields["inner"], &resource.Inner); err != nil {
				errs = append(errs, cog.MakeBuildErrors("inner", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("inner", errors.New("required field is null"))...)

		}
		delete(fields, "inner")
	} else {
		errs = append(errs, cog.MakeBuildErrors("inner", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("EqualityArraysAnonymousStructs", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `EqualityArraysAnonymousStructs` objects.
func (resource EqualityArraysAnonymousStructs) Equals(other EqualityArraysAnonymousStructs) bool {
	if resource.Inner != other.Inner {
		return false
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `EqualityArraysAnonymousStructs` fields for violations and returns them.
func (resource EqualityArraysAnonymousStructs) Validate() error {
	return nil
}

// Modified by compiler pass 'AnonymousStructsToNamed'
type EqualityMapsAnonymousStructs struct {
	Inner string `json:"inner"`
}

// NewEqualityMapsAnonymousStructs creates a new EqualityMapsAnonymousStructs object.
func NewEqualityMapsAnonymousStructs() *EqualityMapsAnonymousStructs {
	return &EqualityMapsAnonymousStructs{}
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `EqualityMapsAnonymousStructs` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *EqualityMapsAnonymousStructs) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "inner"
	if fields["inner"] != nil {
		if string(fields["inner"]) != "null" {
			if err := json.Unmarshal(fields["inner"], &resource.Inner); err != nil {
				errs = append(errs, cog.MakeBuildErrors("inner", err)...)
			}
		} else {
			errs = append(errs, cog.MakeBuildErrors("inner", errors.New("required field is null"))...)

		}
		delete(fields, "inner")
	} else {
		errs = append(errs, cog.MakeBuildErrors("inner", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("EqualityMapsAnonymousStructs", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `EqualityMapsAnonymousStructs` objects.
func (resource EqualityMapsAnonymousStructs) Equals(other EqualityMapsAnonymousStructs) bool {
	if resource.Inner != other.Inner {
		return false
	}

	return true
}

// Validate checks all the validation constraints that may be defined on `EqualityMapsAnonymousStructs` fields for violations and returns them.
func (resource EqualityMapsAnonymousStructs) Validate() error {
	return nil
}
