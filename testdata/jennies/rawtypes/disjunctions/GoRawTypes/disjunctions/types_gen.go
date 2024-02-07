package disjunctions

// Refresh rate or disabled.
type RefreshRate = StringOrBool

type StringOrNull *string

type SomeStruct struct {
	Type string `json:"Type"`
	FieldAny any `json:"FieldAny"`
}

type BoolOrRef = BoolOrSomeStruct

type SomeOtherStruct struct {
	Type string `json:"Type"`
	Foo bytes `json:"Foo"`
}

type YetAnotherStruct struct {
	Type string `json:"Type"`
	Bar uint8 `json:"Bar"`
}

type SeveralRefs = SomeStructOrSomeOtherStructOrYetAnotherStruct

type BoolOrSomeStruct struct {
	Bool *bool `json:"Bool,omitempty"`
	SomeStruct *SomeStruct `json:"SomeStruct,omitempty"`
}

type SomeStructOrSomeOtherStructOrYetAnotherStruct struct {
	SomeStruct *SomeStruct `json:"SomeStruct,omitempty"`
	SomeOtherStruct *SomeOtherStruct `json:"SomeOtherStruct,omitempty"`
	YetAnotherStruct *YetAnotherStruct `json:"YetAnotherStruct,omitempty"`
}

type StringOrBool struct {
	String *string `json:"String,omitempty"`
	Bool *bool `json:"Bool,omitempty"`
}

