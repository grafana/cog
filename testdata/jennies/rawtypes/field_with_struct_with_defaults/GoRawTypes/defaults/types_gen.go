package defaults

import (
	cog "github.com/grafana/cog/generated/cog"
)

type NestedStruct struct {
	StringVal string `json:"stringVal"`
	IntVal int64 `json:"intVal"`
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


func (resource NestedStruct) Validate() error {
	var errs cog.BuildErrors

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type Struct struct {
	AllFields NestedStruct `json:"allFields"`
	PartialFields NestedStruct `json:"partialFields"`
	EmptyFields NestedStruct `json:"emptyFields"`
	ComplexField DefaultsStructComplexField `json:"complexField"`
	PartialComplexField DefaultsStructPartialComplexField `json:"partialComplexField"`
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

func (resource DefaultsStructComplexFieldNested) Equals(other DefaultsStructComplexFieldNested) bool {
		if resource.NestedVal != other.NestedVal {
			return false
		}

	return true
}


func (resource DefaultsStructComplexFieldNested) Validate() error {
	var errs cog.BuildErrors

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type DefaultsStructComplexField struct {
	Uid string `json:"uid"`
	Nested DefaultsStructComplexFieldNested `json:"nested"`
	Array []string `json:"array"`
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

func (resource DefaultsStructPartialComplexField) Equals(other DefaultsStructPartialComplexField) bool {
		if resource.Uid != other.Uid {
			return false
		}
		if resource.IntVal != other.IntVal {
			return false
		}

	return true
}


func (resource DefaultsStructPartialComplexField) Validate() error {
	var errs cog.BuildErrors

	if len(errs) == 0 {
		return nil
	}

	return errs
}


