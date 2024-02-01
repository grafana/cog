package withdashes

type SomeStruct struct {
	FieldAny any `json:"FieldAny"`
}

// Refresh rate or disabled.
type RefreshRate StringOrBool

type StringOrBool struct {
	String *string `json:"String,omitempty"`
	Bool *bool `json:"Bool,omitempty"`
}

