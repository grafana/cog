package default_disjunction_value

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type DisjunctionClasses = ValueAOrValueBOrValueC

// NewDisjunctionClasses creates a new DisjunctionClasses object.
func NewDisjunctionClasses() *DisjunctionClasses {
	return NewValueAOrValueBOrValueC()
}
type ValueA struct {
    Type string `json:"type"`
    AnArray []string `json:"anArray"`
    OtherRef ValueB `json:"otherRef"`
}

// NewValueA creates a new ValueA object.
func NewValueA() *ValueA {
	return &ValueA{
		Type: "A",
		OtherRef: *NewValueB(),
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `ValueA` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *ValueA) UnmarshalJSONStrict(raw []byte) error {
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
	// Field "anArray"
	if fields["anArray"] != nil {
		if string(fields["anArray"]) != "null" {
			
			if err := json.Unmarshal(fields["anArray"], &resource.AnArray); err != nil {
				errs = append(errs, cog.MakeBuildErrors("anArray", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("anArray", errors.New("required field is null"))...)
		
		}
		delete(fields, "anArray")
	} else {errs = append(errs, cog.MakeBuildErrors("anArray", errors.New("required field is missing from input"))...)
	}
	// Field "otherRef"
	if fields["otherRef"] != nil {
		if string(fields["otherRef"]) != "null" {
			
			resource.OtherRef = ValueB{}
			if err := resource.OtherRef.UnmarshalJSONStrict(fields["otherRef"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("otherRef", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("otherRef", errors.New("required field is null"))...)
		
		}
		delete(fields, "otherRef")
	} else {errs = append(errs, cog.MakeBuildErrors("otherRef", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("ValueA", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `ValueA` objects.
func (resource ValueA) Equals(other ValueA) bool {
		if resource.Type != other.Type {
			return false
		}

		if len(resource.AnArray) != len(other.AnArray) {
			return false
		}

		for i1 := range resource.AnArray {
		if resource.AnArray[i1] != other.AnArray[i1] {
			return false
		}
		}
		if !resource.OtherRef.Equals(other.OtherRef) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `ValueA` fields for violations and returns them.
func (resource ValueA) Validate() error {
	var errs cog.BuildErrors
		if err := resource.OtherRef.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("otherRef", err)...)
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type ValueB struct {
    Type string `json:"type"`
    AMap map[string]int64 `json:"aMap"`
    Def Int64OrStringOrBool `json:"def"`
}

// NewValueB creates a new ValueB object.
func NewValueB() *ValueB {
	return &ValueB{
		Type: "B",
		Def: *NewInt64OrStringOrBool(),
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `ValueB` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *ValueB) UnmarshalJSONStrict(raw []byte) error {
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
	// Field "aMap"
	if fields["aMap"] != nil {
		if string(fields["aMap"]) != "null" {
			
			if err := json.Unmarshal(fields["aMap"], &resource.AMap); err != nil {
				errs = append(errs, cog.MakeBuildErrors("aMap", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("aMap", errors.New("required field is null"))...)
		
		}
		delete(fields, "aMap")
	} else {errs = append(errs, cog.MakeBuildErrors("aMap", errors.New("required field is missing from input"))...)
	}
	// Field "def"
	if fields["def"] != nil {
		if string(fields["def"]) != "null" {
			
			resource.Def = Int64OrStringOrBool{}
			if err := resource.Def.UnmarshalJSONStrict(fields["def"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("def", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("def", errors.New("required field is null"))...)
		
		}
		delete(fields, "def")
	} else {errs = append(errs, cog.MakeBuildErrors("def", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("ValueB", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `ValueB` objects.
func (resource ValueB) Equals(other ValueB) bool {
		if resource.Type != other.Type {
			return false
		}

		if len(resource.AMap) != len(other.AMap) {
			return false
		}

		for key1 := range resource.AMap {
		if resource.AMap[key1] != other.AMap[key1] {
			return false
		}
		}
		if !resource.Def.Equals(other.Def) {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `ValueB` fields for violations and returns them.
func (resource ValueB) Validate() error {
	var errs cog.BuildErrors
		if err := resource.Def.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("def", err)...)
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type ValueC struct {
    Type string `json:"type"`
    Other float32 `json:"other"`
}

// NewValueC creates a new ValueC object.
func NewValueC() *ValueC {
	return &ValueC{
		Type: "C",
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `ValueC` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *ValueC) UnmarshalJSONStrict(raw []byte) error {
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
	// Field "other"
	if fields["other"] != nil {
		if string(fields["other"]) != "null" {
			if err := json.Unmarshal(fields["other"], &resource.Other); err != nil {
				errs = append(errs, cog.MakeBuildErrors("other", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("other", errors.New("required field is null"))...)
		
		}
		delete(fields, "other")
	} else {errs = append(errs, cog.MakeBuildErrors("other", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("ValueC", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `ValueC` objects.
func (resource ValueC) Equals(other ValueC) bool {
		if resource.Type != other.Type {
			return false
		}
		if resource.Other != other.Other {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `ValueC` fields for violations and returns them.
func (resource ValueC) Validate() error {
	return nil
}


type DisjunctionConstants = StringOrInt64OrBool

// NewDisjunctionConstants creates a new DisjunctionConstants object.
func NewDisjunctionConstants() *DisjunctionConstants {
	return NewStringOrInt64OrBool()
}
type ValueAOrValueBOrValueC struct {
    ValueA *ValueA `json:"ValueA,omitempty"`
    ValueB *ValueB `json:"ValueB,omitempty"`
    ValueC *ValueC `json:"ValueC,omitempty"`
}

// NewValueAOrValueBOrValueC creates a new ValueAOrValueBOrValueC object.
func NewValueAOrValueBOrValueC() *ValueAOrValueBOrValueC {
	return &ValueAOrValueBOrValueC{
		ValueA: NewValueA(),
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `ValueAOrValueBOrValueC` as JSON.
func (resource ValueAOrValueBOrValueC) MarshalJSON() ([]byte, error) {
	if resource.ValueA != nil {
		return json.Marshal(resource.ValueA)
	}
	if resource.ValueB != nil {
		return json.Marshal(resource.ValueB)
	}
	if resource.ValueC != nil {
		return json.Marshal(resource.ValueC)
	}
	return nil, fmt.Errorf("no value for disjunction of refs")
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `ValueAOrValueBOrValueC` from JSON.
func (resource *ValueAOrValueBOrValueC) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["type"]
	if !found {
		return errors.New("discriminator field 'type' not found in payload")
	}

	switch discriminator {
	case "A":
		var valueA ValueA
		if err := json.Unmarshal(raw, &valueA); err != nil {
			return err
		}

		resource.ValueA = &valueA
		return nil
	case "B":
		var valueB ValueB
		if err := json.Unmarshal(raw, &valueB); err != nil {
			return err
		}

		resource.ValueB = &valueB
		return nil
	case "C":
		var valueC ValueC
		if err := json.Unmarshal(raw, &valueC); err != nil {
			return err
		}

		resource.ValueC = &valueC
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `type = %v`", discriminator)
}


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `ValueAOrValueBOrValueC` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *ValueAOrValueBOrValueC) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["type"]
	if !found {
		return fmt.Errorf("discriminator field 'type' not found in payload")
	}

	switch discriminator {
		case "A":
		valueA := &ValueA{}
		if err := valueA.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.ValueA = valueA
		return nil
		case "B":
		valueB := &ValueB{}
		if err := valueB.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.ValueB = valueB
		return nil
		case "C":
		valueC := &ValueC{}
		if err := valueC.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.ValueC = valueC
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `type = %v`", discriminator)
}

// Equals tests the equality of two `ValueAOrValueBOrValueC` objects.
func (resource ValueAOrValueBOrValueC) Equals(other ValueAOrValueBOrValueC) bool {
		if resource.ValueA == nil && other.ValueA != nil || resource.ValueA != nil && other.ValueA == nil {
			return false
		}

		if resource.ValueA != nil {
		if !resource.ValueA.Equals(*other.ValueA) {
			return false
		}
		}
		if resource.ValueB == nil && other.ValueB != nil || resource.ValueB != nil && other.ValueB == nil {
			return false
		}

		if resource.ValueB != nil {
		if !resource.ValueB.Equals(*other.ValueB) {
			return false
		}
		}
		if resource.ValueC == nil && other.ValueC != nil || resource.ValueC != nil && other.ValueC == nil {
			return false
		}

		if resource.ValueC != nil {
		if !resource.ValueC.Equals(*other.ValueC) {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `ValueAOrValueBOrValueC` fields for violations and returns them.
func (resource ValueAOrValueBOrValueC) Validate() error {
	var errs cog.BuildErrors
		if resource.ValueA != nil {
		if err := resource.ValueA.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("ValueA", err)...)
		}
		}
		if resource.ValueB != nil {
		if err := resource.ValueB.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("ValueB", err)...)
		}
		}
		if resource.ValueC != nil {
		if err := resource.ValueC.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("ValueC", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type Int64OrStringOrBool struct {
    Int64 *int64 `json:"Int64,omitempty"`
    String *string `json:"String,omitempty"`
    Bool *bool `json:"Bool,omitempty"`
}

// NewInt64OrStringOrBool creates a new Int64OrStringOrBool object.
func NewInt64OrStringOrBool() *Int64OrStringOrBool {
	return &Int64OrStringOrBool{
		Int64: (func (input int64) *int64 { return &input })(1),
		String: (func (input string) *string { return &input })("a"),
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `Int64OrStringOrBool` as JSON.
func (resource Int64OrStringOrBool) MarshalJSON() ([]byte, error) {
	if resource.Int64 != nil {
		return json.Marshal(resource.Int64)
	}

	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}

	return nil, fmt.Errorf("no value for disjunction of scalars")
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `Int64OrStringOrBool` from JSON.
func (resource *Int64OrStringOrBool) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// Int64
	var Int64 int64
	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
		resource.Int64 = nil
	} else {
		resource.Int64 = &Int64
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


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Int64OrStringOrBool` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Int64OrStringOrBool) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// Int64
	var Int64 int64

	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
	} else {
		resource.Int64 = &Int64
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

	// Bool
	var Bool bool

	if err := json.Unmarshal(raw, &Bool); err != nil {
		errList = append(errList, err)
	} else {
		resource.Bool = &Bool
		return nil
	}


	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("Int64OrStringOrBool", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `Int64OrStringOrBool` objects.
func (resource Int64OrStringOrBool) Equals(other Int64OrStringOrBool) bool {
		if resource.Int64 == nil && other.Int64 != nil || resource.Int64 != nil && other.Int64 == nil {
			return false
		}

		if resource.Int64 != nil {
		if *resource.Int64 != *other.Int64 {
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
		if resource.Bool == nil && other.Bool != nil || resource.Bool != nil && other.Bool == nil {
			return false
		}

		if resource.Bool != nil {
		if *resource.Bool != *other.Bool {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Int64OrStringOrBool` fields for violations and returns them.
func (resource Int64OrStringOrBool) Validate() error {
	return nil
}


type StringOrInt64OrBool struct {
    String *string `json:"String,omitempty"`
    Int64 *int64 `json:"Int64,omitempty"`
    Bool *bool `json:"Bool,omitempty"`
}

// NewStringOrInt64OrBool creates a new StringOrInt64OrBool object.
func NewStringOrInt64OrBool() *StringOrInt64OrBool {
	return &StringOrInt64OrBool{
		String: (func (input string) *string { return &input })("abc"),
		Int64: (func (input int64) *int64 { return &input })(1),
		Bool: (func (input bool) *bool { return &input })(true),
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `StringOrInt64OrBool` as JSON.
func (resource StringOrInt64OrBool) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Int64 != nil {
		return json.Marshal(resource.Int64)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}

	return nil, fmt.Errorf("no value for disjunction of scalars")
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `StringOrInt64OrBool` from JSON.
func (resource *StringOrInt64OrBool) UnmarshalJSON(raw []byte) error {
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

	// Int64
	var Int64 int64
	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
		resource.Int64 = nil
	} else {
		resource.Int64 = &Int64
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


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `StringOrInt64OrBool` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *StringOrInt64OrBool) UnmarshalJSONStrict(raw []byte) error {
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

	// Int64
	var Int64 int64

	if err := json.Unmarshal(raw, &Int64); err != nil {
		errList = append(errList, err)
	} else {
		resource.Int64 = &Int64
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


	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("StringOrInt64OrBool", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `StringOrInt64OrBool` objects.
func (resource StringOrInt64OrBool) Equals(other StringOrInt64OrBool) bool {
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
		if resource.Bool == nil && other.Bool != nil || resource.Bool != nil && other.Bool == nil {
			return false
		}

		if resource.Bool != nil {
		if *resource.Bool != *other.Bool {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `StringOrInt64OrBool` fields for violations and returns them.
func (resource StringOrInt64OrBool) Validate() error {
	return nil
}
