package disjunctions_of_scalars_and_refs

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"bytes"
)

type DisjunctionOfScalarsAndRefs = StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB

// NewDisjunctionOfScalarsAndRefs creates a new DisjunctionOfScalarsAndRefs object.
func NewDisjunctionOfScalarsAndRefs() *DisjunctionOfScalarsAndRefs {
	return NewStringOrBoolOrArrayOfStringOrMyRefAOrMyRefB()
}
type MyRefA struct {
    Foo string `json:"foo"`
}

// NewMyRefA creates a new MyRefA object.
func NewMyRefA() *MyRefA {
	return &MyRefA{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `MyRefA` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *MyRefA) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "foo"
	if fields["foo"] != nil {
		if string(fields["foo"]) != "null" {
			if err := json.Unmarshal(fields["foo"], &resource.Foo); err != nil {
				errs = append(errs, cog.MakeBuildErrors("foo", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("foo", errors.New("required field is null"))...)
		
		}
		delete(fields, "foo")
	} else {errs = append(errs, cog.MakeBuildErrors("foo", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("MyRefA", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `MyRefA` objects.
func (resource MyRefA) Equals(other MyRefA) bool {
		if resource.Foo != other.Foo {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `MyRefA` fields for violations and returns them.
func (resource MyRefA) Validate() error {
	return nil
}


type MyRefB struct {
    Bar int64 `json:"bar"`
}

// NewMyRefB creates a new MyRefB object.
func NewMyRefB() *MyRefB {
	return &MyRefB{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `MyRefB` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *MyRefB) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "bar"
	if fields["bar"] != nil {
		if string(fields["bar"]) != "null" {
			if err := json.Unmarshal(fields["bar"], &resource.Bar); err != nil {
				errs = append(errs, cog.MakeBuildErrors("bar", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("bar", errors.New("required field is null"))...)
		
		}
		delete(fields, "bar")
	} else {errs = append(errs, cog.MakeBuildErrors("bar", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("MyRefB", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `MyRefB` objects.
func (resource MyRefB) Equals(other MyRefB) bool {
		if resource.Bar != other.Bar {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `MyRefB` fields for violations and returns them.
func (resource MyRefB) Validate() error {
	return nil
}


type StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB struct {
    String *string `json:"String,omitempty"`
    Bool *bool `json:"Bool,omitempty"`
    ArrayOfString []string `json:"ArrayOfString,omitempty"`
    MyRefA *MyRefA `json:"MyRefA,omitempty"`
    MyRefB *MyRefB `json:"MyRefB,omitempty"`
}

// NewStringOrBoolOrArrayOfStringOrMyRefAOrMyRefB creates a new StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB object.
func NewStringOrBoolOrArrayOfStringOrMyRefAOrMyRefB() *StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB {
	return &StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB{
		String: (func (input string) *string { return &input })("a"),
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB` as JSON.
func (resource StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}
	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}
	if resource.ArrayOfString != nil {
		return json.Marshal(resource.ArrayOfString)
	}
	if resource.MyRefA != nil {
		return json.Marshal(resource.MyRefA)
	}
	if resource.MyRefB != nil {
		return json.Marshal(resource.MyRefB)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB` from JSON.
func (resource *StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB) UnmarshalJSON(raw []byte) error {
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

	// ArrayOfString
	var ArrayOfString []string
	if err := json.Unmarshal(raw, &ArrayOfString); err != nil {
		errList = append(errList, err)
		resource.ArrayOfString = nil
	} else {
		resource.ArrayOfString = ArrayOfString
		return nil
	}

	// MyRefA
	var MyRefA MyRefA
    myRefAdec := json.NewDecoder(bytes.NewReader(raw))
    myRefAdec.DisallowUnknownFields()
    if err := myRefAdec.Decode(&MyRefA); err != nil {
        errList = append(errList, err)
        resource.MyRefA = nil
    } else {
        resource.MyRefA = &MyRefA
        return nil
    }

	// MyRefB
	var MyRefB MyRefB
    myRefBdec := json.NewDecoder(bytes.NewReader(raw))
    myRefBdec.DisallowUnknownFields()
    if err := myRefBdec.Decode(&MyRefB); err != nil {
        errList = append(errList, err)
        resource.MyRefB = nil
    } else {
        resource.MyRefB = &MyRefB
        return nil
    }

	return errors.Join(errList...)
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// String
	var String string
	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
	} else {
		resource.String = &String
		return nil
	}

	// Bool
	var Bool bool
	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
	} else {
		resource.Bool = &Bool
		return nil
	}

	// ArrayOfString
	var ArrayOfString []string
	if err := json.Unmarshal(raw, &ArrayOfString); err != nil {
		errList = append(errList, err)
	} else {
		resource.ArrayOfString = ArrayOfString
		return nil
	}

	// MyRefA
	var MyRefA MyRefA
    myRefAdec := json.NewDecoder(bytes.NewReader(raw))
    myRefAdec.DisallowUnknownFields()
    if err := myRefAdec.Decode(&MyRefA); err != nil {
        errList = append(errList, err)
    } else {
        resource.MyRefA = &MyRefA
        return nil
    }

	// MyRefB
	var MyRefB MyRefB
    myRefBdec := json.NewDecoder(bytes.NewReader(raw))
    myRefBdec.DisallowUnknownFields()
    if err := myRefBdec.Decode(&MyRefB); err != nil {
        errList = append(errList, err)
    } else {
        resource.MyRefB = &MyRefB
        return nil
    }

	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB` objects.
func (resource StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB) Equals(other StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB) bool {
		if resource.String == nil && other.String != nil || resource.String != nil && other.String == nil {
			return false
		}

		if resource.String != nil {
		if *resource.String != *other.String {
			return false
		}
		}
		if resource.Bool == nil && other.Bool != nil || resource.Bool != nil && other.Bool == nil {
			return false
		}

		if resource.Bool != nil {
		if *resource.Bool != *other.Bool {
			return false
		}
		}

		if len(resource.ArrayOfString) != len(other.ArrayOfString) {
			return false
		}

		for i1 := range resource.ArrayOfString {
		if resource.ArrayOfString[i1] != other.ArrayOfString[i1] {
			return false
		}
		}
		if resource.MyRefA == nil && other.MyRefA != nil || resource.MyRefA != nil && other.MyRefA == nil {
			return false
		}

		if resource.MyRefA != nil {
		if !resource.MyRefA.Equals(*other.MyRefA) {
			return false
		}
		}
		if resource.MyRefB == nil && other.MyRefB != nil || resource.MyRefB != nil && other.MyRefB == nil {
			return false
		}

		if resource.MyRefB != nil {
		if !resource.MyRefB.Equals(*other.MyRefB) {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB` fields for violations and returns them.
func (resource StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB) Validate() error {
	var errs cog.BuildErrors
		if resource.MyRefA != nil {
		if err := resource.MyRefA.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("MyRefA", err)...)
		}
		}
		if resource.MyRefB != nil {
		if err := resource.MyRefB.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("MyRefB", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
