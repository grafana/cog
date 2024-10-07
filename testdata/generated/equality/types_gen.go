// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Using jennies:
//     GoRawTypes

package equality

import (
	reflect "reflect"
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

func (resource Variable) Equals(other Variable) bool {
	if resource.Name != other.Name {
		return false
	}

	return true
}

type Container struct {
	StringField string    `json:"stringField"`
	IntField    int64     `json:"intField"`
	EnumField   Direction `json:"enumField"`
	RefField    Variable  `json:"refField"`
}

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

type Optionals struct {
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	StringField *string `json:"stringField,omitempty"`
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	EnumField *Direction `json:"enumField,omitempty"`
	// Modified by compiler pass 'NotRequiredFieldAsNullableType[nullable=true]'
	RefField *Variable `json:"refField,omitempty"`
}

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

	return true
}

type Arrays struct {
	Ints             []int64                          `json:"ints"`
	Strings          []string                         `json:"strings"`
	ArrayOfArray     [][]string                       `json:"arrayOfArray"`
	Refs             []Variable                       `json:"refs"`
	AnonymousStructs []EqualityArraysAnonymousStructs `json:"anonymousStructs"`
	ArrayOfAny       []any                            `json:"arrayOfAny"`
}

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

type Maps struct {
	Ints             map[string]int64                        `json:"ints"`
	Strings          map[string]string                       `json:"strings"`
	Refs             map[string]Variable                     `json:"refs"`
	AnonymousStructs map[string]EqualityMapsAnonymousStructs `json:"anonymousStructs"`
	StringToAny      map[string]any                          `json:"stringToAny"`
}

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

// Modified by compiler pass 'AnonymousStructsToNamed'
type EqualityArraysAnonymousStructs struct {
	Inner string `json:"inner"`
}

func (resource EqualityArraysAnonymousStructs) Equals(other EqualityArraysAnonymousStructs) bool {
	if resource.Inner != other.Inner {
		return false
	}

	return true
}

// Modified by compiler pass 'AnonymousStructsToNamed'
type EqualityMapsAnonymousStructs struct {
	Inner string `json:"inner"`
}

func (resource EqualityMapsAnonymousStructs) Equals(other EqualityMapsAnonymousStructs) bool {
	if resource.Inner != other.Inner {
		return false
	}

	return true
}
