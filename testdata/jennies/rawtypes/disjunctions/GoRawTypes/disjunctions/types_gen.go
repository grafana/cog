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
	Foo []byte `json:"Foo"`
}

type YetAnotherStruct struct {
	Type string `json:"Type"`
	Bar uint8 `json:"Bar"`
}

type SeveralRefs = SomeStructOrSomeOtherStructOrYetAnotherStruct

type StringOrBool struct {
	String *string `json:"String,omitempty"`
	Bool *bool `json:"Bool,omitempty"`
}

func (resource StringOrBool) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}

	return nil, fmt.Errorf("no value for disjunction of scalars")
}


func (resource *StringOrBool) UnmarshalJSON(raw []byte) error {
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

	// Bool
	var Bool bool
	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
		resource.Bool = nil
	} else {
		resource.Bool = &Bool
		return nil
	}

	return errors.Join(errList...)
}


type BoolOrSomeStruct struct {
	Bool *bool `json:"Bool,omitempty"`
	SomeStruct *SomeStruct `json:"SomeStruct,omitempty"`
}

type SomeStructOrSomeOtherStructOrYetAnotherStruct struct {
	SomeStruct *SomeStruct `json:"SomeStruct,omitempty"`
	SomeOtherStruct *SomeOtherStruct `json:"SomeOtherStruct,omitempty"`
	YetAnotherStruct *YetAnotherStruct `json:"YetAnotherStruct,omitempty"`
}

func (resource SomeStructOrSomeOtherStructOrYetAnotherStruct) MarshalJSON() ([]byte, error) {
	if resource.SomeStruct != nil {
		return json.Marshal(resource.SomeStruct)
	}
	if resource.SomeOtherStruct != nil {
		return json.Marshal(resource.SomeOtherStruct)
	}
	if resource.YetAnotherStruct != nil {
		return json.Marshal(resource.YetAnotherStruct)
	}

	return nil, fmt.Errorf("no value for disjunction of refs")
}

func (resource *SomeStructOrSomeOtherStructOrYetAnotherStruct) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["Type"]
	if !found {
		return errors.New("discriminator field 'Type' not found in payload")
	}

	switch discriminator {
	case "some-other-struct":
		var someOtherStruct SomeOtherStruct
		if err := json.Unmarshal(raw, &someOtherStruct); err != nil {
			return err
		}

		resource.SomeOtherStruct = &someOtherStruct
		return nil
	case "some-struct":
		var someStruct SomeStruct
		if err := json.Unmarshal(raw, &someStruct); err != nil {
			return err
		}

		resource.SomeStruct = &someStruct
		return nil
	case "yet-another-struct":
		var yetAnotherStruct YetAnotherStruct
		if err := json.Unmarshal(raw, &yetAnotherStruct); err != nil {
			return err
		}

		resource.YetAnotherStruct = &yetAnotherStruct
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `Type = %v`", discriminator)
}


