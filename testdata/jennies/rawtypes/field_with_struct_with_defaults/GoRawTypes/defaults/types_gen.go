package defaults

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


type DefaultsStructComplexFieldNested struct {
	NestedVal string `json:"nestedVal"`
}

func (resource DefaultsStructComplexFieldNested) Equals(other DefaultsStructComplexFieldNested) bool {
		if resource.NestedVal != other.NestedVal {
			return false
		}

	return true
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


