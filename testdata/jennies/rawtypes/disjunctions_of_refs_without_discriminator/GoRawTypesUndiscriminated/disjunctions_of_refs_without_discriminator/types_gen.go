package disjunctions_of_refs_without_discriminator

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
	"bytes"
)

type DisjunctionWithoutDiscriminator = TypeAOrTypeB

// NewDisjunctionWithoutDiscriminator creates a new DisjunctionWithoutDiscriminator object.
func NewDisjunctionWithoutDiscriminator() *DisjunctionWithoutDiscriminator {
	return NewTypeAOrTypeB()
}
type TypeA struct {
    FieldA string `json:"fieldA"`
}

// NewTypeA creates a new TypeA object.
func NewTypeA() *TypeA {
	return &TypeA{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `TypeA` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *TypeA) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "fieldA"
	if fields["fieldA"] != nil {
		if string(fields["fieldA"]) != "null" {
			if err := json.Unmarshal(fields["fieldA"], &resource.FieldA); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fieldA", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("fieldA", errors.New("required field is null"))...)
		
		}
		delete(fields, "fieldA")
	} else {errs = append(errs, cog.MakeBuildErrors("fieldA", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("TypeA", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `TypeA` objects.
func (resource TypeA) Equals(other TypeA) bool {
		if resource.FieldA != other.FieldA {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `TypeA` fields for violations and returns them.
func (resource TypeA) Validate() error {
	return nil
}


type TypeB struct {
    FieldB int64 `json:"fieldB"`
}

// NewTypeB creates a new TypeB object.
func NewTypeB() *TypeB {
	return &TypeB{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `TypeB` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *TypeB) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "fieldB"
	if fields["fieldB"] != nil {
		if string(fields["fieldB"]) != "null" {
			if err := json.Unmarshal(fields["fieldB"], &resource.FieldB); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fieldB", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("fieldB", errors.New("required field is null"))...)
		
		}
		delete(fields, "fieldB")
	} else {errs = append(errs, cog.MakeBuildErrors("fieldB", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("TypeB", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `TypeB` objects.
func (resource TypeB) Equals(other TypeB) bool {
		if resource.FieldB != other.FieldB {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `TypeB` fields for violations and returns them.
func (resource TypeB) Validate() error {
	return nil
}


type TypeAOrTypeB struct {
    TypeA *TypeA `json:"TypeA,omitempty"`
    TypeB *TypeB `json:"TypeB,omitempty"`
}

// NewTypeAOrTypeB creates a new TypeAOrTypeB object.
func NewTypeAOrTypeB() *TypeAOrTypeB {
	return &TypeAOrTypeB{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `TypeAOrTypeB` as JSON.
func (resource TypeAOrTypeB) MarshalJSON() ([]byte, error) {
	if resource.TypeA != nil {
		return json.Marshal(resource.TypeA)
	}
	if resource.TypeB != nil {
		return json.Marshal(resource.TypeB)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `TypeAOrTypeB` from JSON.
func (resource *TypeAOrTypeB) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// TypeA
	var TypeA TypeA
	typeAdec := json.NewDecoder(bytes.NewReader(raw))
	typeAdec.DisallowUnknownFields()
	if err := typeAdec.Decode(&TypeA); err != nil {
		errList = append(errList, err)
		resource.TypeA = nil
	} else {
		resource.TypeA = &TypeA
		return nil
	}

	// TypeB
	var TypeB TypeB
	typeBdec := json.NewDecoder(bytes.NewReader(raw))
	typeBdec.DisallowUnknownFields()
	if err := typeBdec.Decode(&TypeB); err != nil {
		errList = append(errList, err)
		resource.TypeB = nil
	} else {
		resource.TypeB = &TypeB
		return nil
	}

	return errors.Join(errList...)
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `TypeAOrTypeB` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *TypeAOrTypeB) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors
	var errList []error

	// TypeA
	var TypeA TypeA
	typeAdec := json.NewDecoder(bytes.NewReader(raw))
	typeAdec.DisallowUnknownFields()
	if err := typeAdec.Decode(&TypeA); err != nil {
		errList = append(errList, err)
	} else {
		resource.TypeA = &TypeA
		return nil
	}

	// TypeB
	var TypeB TypeB
	typeBdec := json.NewDecoder(bytes.NewReader(raw))
	typeBdec.DisallowUnknownFields()
	if err := typeBdec.Decode(&TypeB); err != nil {
		errList = append(errList, err)
	} else {
		resource.TypeB = &TypeB
		return nil
	}

	if len(errList) != 0 {
		errs = append(errs, cog.MakeBuildErrors("TypeAOrTypeB", errors.Join(errList...))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

// Equals tests the equality of two `TypeAOrTypeB` objects.
func (resource TypeAOrTypeB) Equals(other TypeAOrTypeB) bool {
		if resource.TypeA == nil && other.TypeA != nil || resource.TypeA != nil && other.TypeA == nil {
			return false
		}

		if resource.TypeA != nil {
		if !resource.TypeA.Equals(*other.TypeA) {
			return false
		}
		}
		if resource.TypeB == nil && other.TypeB != nil || resource.TypeB != nil && other.TypeB == nil {
			return false
		}

		if resource.TypeB != nil {
		if !resource.TypeB.Equals(*other.TypeB) {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `TypeAOrTypeB` fields for violations and returns them.
func (resource TypeAOrTypeB) Validate() error {
	var errs cog.BuildErrors
		if resource.TypeA != nil {
		if err := resource.TypeA.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("TypeA", err)...)
		}
		}
		if resource.TypeB != nil {
		if err := resource.TypeB.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("TypeB", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


