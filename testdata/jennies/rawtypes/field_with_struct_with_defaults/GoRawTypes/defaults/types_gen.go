package defaults

type NestedStruct struct {
	StringVal string `json:"stringVal"`
	IntVal int64 `json:"intVal"`
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

