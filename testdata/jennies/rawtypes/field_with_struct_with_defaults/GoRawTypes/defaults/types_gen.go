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
	ComplexField struct {
	Uid string `json:"uid"`
	Nested struct {
	NestedVal string `json:"nestedVal"`
} `json:"nested"`
	Array []string `json:"array"`
} `json:"complexField"`
	PartialComplexField struct {
	Uid string `json:"uid"`
	IntVal int64 `json:"intVal"`
} `json:"partialComplexField"`
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
		if resource.ComplexField.Uid != other.ComplexField.Uid {
			return false
		}
		if resource.ComplexField.Nested.NestedVal != other.ComplexField.Nested.NestedVal {
			return false
		}

		if len(resource.ComplexField.Array) != len(other.ComplexField.Array) {
			return false
		}

		for i1 := range resource.ComplexField.Array {
		if resource.ComplexField.Array[i1] != other.ComplexField.Array[i1] {
			return false
		}
		}
		if resource.PartialComplexField.Uid != other.PartialComplexField.Uid {
			return false
		}
		if resource.PartialComplexField.IntVal != other.PartialComplexField.IntVal {
			return false
		}

	return true
}


