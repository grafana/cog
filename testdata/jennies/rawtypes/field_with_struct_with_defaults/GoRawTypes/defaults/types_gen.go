package defaults

type NestedStruct struct {
	StringVal string `json:"stringVal" yaml:"stringVal"`
	IntVal int64 `json:"intVal" yaml:"intVal"`
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


type Struct struct {
	AllFields NestedStruct `json:"allFields" yaml:"allFields"`
	PartialFields NestedStruct `json:"partialFields" yaml:"partialFields"`
	EmptyFields NestedStruct `json:"emptyFields" yaml:"emptyFields"`
	ComplexField DefaultsStructComplexField `json:"complexField" yaml:"complexField"`
	PartialComplexField DefaultsStructPartialComplexField `json:"partialComplexField" yaml:"partialComplexField"`
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


type DefaultsStructComplexFieldNested struct {
	NestedVal string `json:"nestedVal" yaml:"nestedVal"`
}

func (resource DefaultsStructComplexFieldNested) Equals(other DefaultsStructComplexFieldNested) bool {
		if resource.NestedVal != other.NestedVal {
			return false
		}

	return true
}


type DefaultsStructComplexField struct {
	Uid string `json:"uid" yaml:"uid"`
	Nested DefaultsStructComplexFieldNested `json:"nested" yaml:"nested"`
	Array []string `json:"array" yaml:"array"`
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


type DefaultsStructPartialComplexField struct {
	Uid string `json:"uid" yaml:"uid"`
	IntVal int64 `json:"intVal" yaml:"intVal"`
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


