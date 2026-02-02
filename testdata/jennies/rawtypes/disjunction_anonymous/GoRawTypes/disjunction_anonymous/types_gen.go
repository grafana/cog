package disjunction_anonymous

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"reflect"
	"bytes"
)

type MyStruct struct {
    Scalars StringOrBoolOrFloat64OrInt64 `json:"scalars"`
    SameKind MyStructSameKind `json:"sameKind"`
    Refs any `json:"refs"`
    Mixed StructAOrStringOrInt64 `json:"mixed"`
}

// NewMyStruct creates a new MyStruct object.
func NewMyStruct() *MyStruct {
	return &MyStruct{
		Scalars: *NewStringOrBoolOrFloat64OrInt64(),
		Mixed: *NewStructAOrStringOrInt64(),
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `MyStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *MyStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "scalars"
	if fields["scalars"] != nil {
		if string(fields["scalars"]) != "null" {
			
			resource.Scalars = StringOrBoolOrFloat64OrInt64{}
			if err := resource.Scalars.UnmarshalJSONStrict(fields["scalars"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("scalars", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("scalars", errors.New("required field is null"))...)
		
		}
		delete(fields, "scalars")
	} else {errs = append(errs, cog.MakeBuildErrors("scalars", errors.New("required field is missing from input"))...)
	}
	// Field "sameKind"
	if fields["sameKind"] != nil {
		if string(fields["sameKind"]) != "null" {
			if err := json.Unmarshal(fields["sameKind"], &resource.SameKind); err != nil {
				errs = append(errs, cog.MakeBuildErrors("sameKind", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("sameKind", errors.New("required field is null"))...)
		
		}
		delete(fields, "sameKind")
	} else {errs = append(errs, cog.MakeBuildErrors("sameKind", errors.New("required field is missing from input"))...)
	}
	// Field "refs"
	if fields["refs"] != nil {
		if string(fields["refs"]) != "null" {
			if err := json.Unmarshal(fields["refs"], &resource.Refs); err != nil {
				errs = append(errs, cog.MakeBuildErrors("refs", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("refs", errors.New("required field is null"))...)
		
		}
		delete(fields, "refs")
	} else {errs = append(errs, cog.MakeBuildErrors("refs", errors.New("required field is missing from input"))...)
	}
	// Field "mixed"
	if fields["mixed"] != nil {
		if string(fields["mixed"]) != "null" {
			
			resource.Mixed = StructAOrStringOrInt64{}
			if err := resource.Mixed.UnmarshalJSONStrict(fields["mixed"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("mixed", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("mixed", errors.New("required field is null"))...)
		
		}
		delete(fields, "mixed")
	} else {errs = append(errs, cog.MakeBuildErrors("mixed", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("MyStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `MyStruct` objects.
func (resource MyStruct) Equals(other MyStruct) bool {
		if !resource.Scalars.Equals(other.Scalars) {
			return false
		}
		if resource.SameKind != other.SameKind {
			return false
		}
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.Refs, other.Refs) {
			return false
		}
		if !resource.Mixed.Equals(other.Mixed) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `MyStruct` fields for violations and returns them.
func (resource MyStruct) Validate() error {
	var errs cog.BuildErrors
		if err := resource.Scalars.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("scalars", err)...)
		}
		if err := resource.Mixed.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("mixed", err)...)
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type StructA struct {
    Field string `json:"field"`
}

// NewStructA creates a new StructA object.
func NewStructA() *StructA {
	return &StructA{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StructA` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StructA) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "field"
	if fields["field"] != nil {
		if string(fields["field"]) != "null" {
			if err := json.Unmarshal(fields["field"], &resource.Field); err != nil {
				errs = append(errs, cog.MakeBuildErrors("field", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("field", errors.New("required field is null"))...)
		
		}
		delete(fields, "field")
	} else {errs = append(errs, cog.MakeBuildErrors("field", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("StructA", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `StructA` objects.
func (resource StructA) Equals(other StructA) bool {
		if resource.Field != other.Field {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StructA` fields for violations and returns them.
func (resource StructA) Validate() error {
	return nil
}


type StructB struct {
    Type int64 `json:"type"`
}

// NewStructB creates a new StructB object.
func NewStructB() *StructB {
	return &StructB{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StructB` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StructB) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "type"
	if fields["type"] != nil {
		if string(fields["type"]) != "null" {
			if err := json.Unmarshal(fields["type"], &resource.Type); err != nil {
				errs = append(errs, cog.MakeBuildErrors("type", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("type", errors.New("required field is null"))...)
		
		}
		delete(fields, "type")
	} else {errs = append(errs, cog.MakeBuildErrors("type", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("StructB", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `StructB` objects.
func (resource StructB) Equals(other StructB) bool {
		if resource.Type != other.Type {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StructB` fields for violations and returns them.
func (resource StructB) Validate() error {
	return nil
}


type MyStructSameKind string
const (
	MyStructSameKindA MyStructSameKind = "a"
	MyStructSameKindB MyStructSameKind = "b"
	MyStructSameKindC MyStructSameKind = "c"
)


type StringOrBoolOrFloat64OrInt64 struct {
    String *string `json:"String,omitempty"`
    Bool *bool `json:"Bool,omitempty"`
    Float64 *float64 `json:"Float64,omitempty"`
    Int64 *int64 `json:"Int64,omitempty"`
}

// NewStringOrBoolOrFloat64OrInt64 creates a new StringOrBoolOrFloat64OrInt64 object.
func NewStringOrBoolOrFloat64OrInt64() *StringOrBoolOrFloat64OrInt64 {
	return &StringOrBoolOrFloat64OrInt64{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `StringOrBoolOrFloat64OrInt64` as JSON.
func (resource StringOrBoolOrFloat64OrInt64) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}

	if resource.Float64 != nil {
		return json.Marshal(resource.Float64)
	}

	if resource.Int64 != nil {
		return json.Marshal(resource.Int64)
	}


	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StringOrBoolOrFloat64OrInt64` from JSON.
func (resource *StringOrBoolOrFloat64OrInt64) UnmarshalJSON(raw []byte) error {
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

	// Float64
	var Float64 float64
	if err := json.Unmarshal(raw, &Float64); err != nil {
		errList = append(errList, err)
		resource.Float64 = nil
	} else {
		resource.Float64 = &Float64
		return nil
	}

	// Int64
	var Int64 int64
	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
		resource.Int64 = nil
	} else {
		resource.Int64 = &Int64
		return nil
	}

	return errors.Join(errList...)
}


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StringOrBoolOrFloat64OrInt64` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StringOrBoolOrFloat64OrInt64) UnmarshalJSONStrict(raw []byte) error {
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

	// Float64
	var Float64 float64

	if err := json.Unmarshal(raw, &Float64); err != nil {
		errList = append(errList, err)
	} else {
		resource.Float64 = &Float64
		return nil
	}

	// Int64
	var Int64 int64

	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
	} else {
		resource.Int64 = &Int64
		return nil
	}


	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StringOrBoolOrFloat64OrInt64", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StringOrBoolOrFloat64OrInt64` objects.
func (resource StringOrBoolOrFloat64OrInt64) Equals(other StringOrBoolOrFloat64OrInt64) bool {
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
		if resource.Float64 == nil && other.Float64 != nil || resource.Float64 != nil && other.Float64 == nil {
			return false
		}

		if resource.Float64 != nil {
		if *resource.Float64 != *other.Float64 {
			return false
		}
		}
		if resource.Int64 == nil && other.Int64 != nil || resource.Int64 != nil && other.Int64 == nil {
			return false
		}

		if resource.Int64 != nil {
		if *resource.Int64 != *other.Int64 {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StringOrBoolOrFloat64OrInt64` fields for violations and returns them.
func (resource StringOrBoolOrFloat64OrInt64) Validate() error {
	return nil
}


type StructAOrStringOrInt64 struct {
    StructA *StructA `json:"StructA,omitempty"`
    String *string `json:"String,omitempty"`
    Int64 *int64 `json:"Int64,omitempty"`
}

// NewStructAOrStringOrInt64 creates a new StructAOrStringOrInt64 object.
func NewStructAOrStringOrInt64() *StructAOrStringOrInt64 {
	return &StructAOrStringOrInt64{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `StructAOrStringOrInt64` as JSON.
func (resource StructAOrStringOrInt64) MarshalJSON() ([]byte, error) {
	if resource.StructA != nil {
		return json.Marshal(resource.StructA)
	}
	if resource.String != nil {
		return json.Marshal(resource.String)
	}
	if resource.Int64 != nil {
		return json.Marshal(resource.Int64)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StructAOrStringOrInt64` from JSON.
func (resource *StructAOrStringOrInt64) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// StructA
	var StructA StructA
    structAdec := json.NewDecoder(bytes.NewReader(raw))
    structAdec.DisallowUnknownFields()
    if err := structAdec.Decode(&StructA); err != nil {
        errList = append(errList, err)
        resource.StructA = nil
    } else {
        resource.StructA = &StructA
        return nil
    }

	// String
	var String string
	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
		resource.String = nil
	} else {
		resource.String = &String
		return nil
	}

	// Int64
	var Int64 int64
	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
		resource.Int64 = nil
	} else {
		resource.Int64 = &Int64
		return nil
	}

	return errors.Join(errList...)
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StructAOrStringOrInt64` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StructAOrStringOrInt64) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// StructA
	var StructA StructA
    structAdec := json.NewDecoder(bytes.NewReader(raw))
    structAdec.DisallowUnknownFields()
    if err := structAdec.Decode(&StructA); err != nil {
        errList = append(errList, err)
    } else {
        resource.StructA = &StructA
        return nil
    }

	// String
	var String string
	if err := json.Unmarshal(raw, &String); err != nil {
		errList = append(errList, err)
	} else {
		resource.String = &String
		return nil
	}

	// Int64
	var Int64 int64
	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
	} else {
		resource.Int64 = &Int64
		return nil
	}

	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StructAOrStringOrInt64", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StructAOrStringOrInt64` objects.
func (resource StructAOrStringOrInt64) Equals(other StructAOrStringOrInt64) bool {
		if resource.StructA == nil && other.StructA != nil || resource.StructA != nil && other.StructA == nil {
			return false
		}

		if resource.StructA != nil {
		if !resource.StructA.Equals(*other.StructA) {
			return false
		}
		}
		if resource.String == nil && other.String != nil || resource.String != nil && other.String == nil {
			return false
		}

		if resource.String != nil {
		if *resource.String != *other.String {
			return false
		}
		}
		if resource.Int64 == nil && other.Int64 != nil || resource.Int64 != nil && other.Int64 == nil {
			return false
		}

		if resource.Int64 != nil {
		if *resource.Int64 != *other.Int64 {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StructAOrStringOrInt64` fields for violations and returns them.
func (resource StructAOrStringOrInt64) Validate() error {
	var errs cog.BuildErrors
		if resource.StructA != nil {
		if err := resource.StructA.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("StructA", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


