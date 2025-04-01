package constant_reference_discriminator

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type LayoutWithValue any

type GridLayoutUsingValue struct {
    Kind string `json:"kind"`
    GridLayoutProperty string `json:"gridLayoutProperty"`
}

// NewGridLayoutUsingValue creates a new GridLayoutUsingValue object.
func NewGridLayoutUsingValue() *GridLayoutUsingValue {
	return &GridLayoutUsingValue{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `GridLayoutUsingValue` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *GridLayoutUsingValue) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "kind"
	if fields["kind"] != nil {
		if string(fields["kind"]) != "null" {
			if err := json.Unmarshal(fields["kind"], &resource.Kind); err != nil {
				errs = append(errs, cog.MakeBuildErrors("kind", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is null"))...)
		
		}
		delete(fields, "kind")
	} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is missing from input"))...)
	}
	// Field "gridLayoutProperty"
	if fields["gridLayoutProperty"] != nil {
		if string(fields["gridLayoutProperty"]) != "null" {
			if err := json.Unmarshal(fields["gridLayoutProperty"], &resource.GridLayoutProperty); err != nil {
				errs = append(errs, cog.MakeBuildErrors("gridLayoutProperty", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("gridLayoutProperty", errors.New("required field is null"))...)
		
		}
		delete(fields, "gridLayoutProperty")
	} else {errs = append(errs, cog.MakeBuildErrors("gridLayoutProperty", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("GridLayoutUsingValue", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `GridLayoutUsingValue` objects.
func (resource GridLayoutUsingValue) Equals(other GridLayoutUsingValue) bool {
		if resource.Kind != other.Kind {
			return false
		}
		if resource.GridLayoutProperty != other.GridLayoutProperty {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `GridLayoutUsingValue` fields for violations and returns them.
func (resource GridLayoutUsingValue) Validate() error {
	return nil
}


const GridLayoutKindType = "GridLayout"

type RowsLayoutUsingValue struct {
    Kind string `json:"kind"`
    RowsLayoutProperty string `json:"rowsLayoutProperty"`
}

// NewRowsLayoutUsingValue creates a new RowsLayoutUsingValue object.
func NewRowsLayoutUsingValue() *RowsLayoutUsingValue {
	return &RowsLayoutUsingValue{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `RowsLayoutUsingValue` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *RowsLayoutUsingValue) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "kind"
	if fields["kind"] != nil {
		if string(fields["kind"]) != "null" {
			if err := json.Unmarshal(fields["kind"], &resource.Kind); err != nil {
				errs = append(errs, cog.MakeBuildErrors("kind", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is null"))...)
		
		}
		delete(fields, "kind")
	} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is missing from input"))...)
	}
	// Field "rowsLayoutProperty"
	if fields["rowsLayoutProperty"] != nil {
		if string(fields["rowsLayoutProperty"]) != "null" {
			if err := json.Unmarshal(fields["rowsLayoutProperty"], &resource.RowsLayoutProperty); err != nil {
				errs = append(errs, cog.MakeBuildErrors("rowsLayoutProperty", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("rowsLayoutProperty", errors.New("required field is null"))...)
		
		}
		delete(fields, "rowsLayoutProperty")
	} else {errs = append(errs, cog.MakeBuildErrors("rowsLayoutProperty", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("RowsLayoutUsingValue", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `RowsLayoutUsingValue` objects.
func (resource RowsLayoutUsingValue) Equals(other RowsLayoutUsingValue) bool {
		if resource.Kind != other.Kind {
			return false
		}
		if resource.RowsLayoutProperty != other.RowsLayoutProperty {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `RowsLayoutUsingValue` fields for violations and returns them.
func (resource RowsLayoutUsingValue) Validate() error {
	return nil
}


const RowsLayoutKindType = "RowsLayout"

type LayoutWithoutValue any

type GridLayoutWithoutValue struct {
    Kind string `json:"kind"`
    GridLayoutProperty string `json:"gridLayoutProperty"`
}

// NewGridLayoutWithoutValue creates a new GridLayoutWithoutValue object.
func NewGridLayoutWithoutValue() *GridLayoutWithoutValue {
	return &GridLayoutWithoutValue{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `GridLayoutWithoutValue` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *GridLayoutWithoutValue) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "kind"
	if fields["kind"] != nil {
		if string(fields["kind"]) != "null" {
			if err := json.Unmarshal(fields["kind"], &resource.Kind); err != nil {
				errs = append(errs, cog.MakeBuildErrors("kind", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is null"))...)
		
		}
		delete(fields, "kind")
	} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is missing from input"))...)
	}
	// Field "gridLayoutProperty"
	if fields["gridLayoutProperty"] != nil {
		if string(fields["gridLayoutProperty"]) != "null" {
			if err := json.Unmarshal(fields["gridLayoutProperty"], &resource.GridLayoutProperty); err != nil {
				errs = append(errs, cog.MakeBuildErrors("gridLayoutProperty", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("gridLayoutProperty", errors.New("required field is null"))...)
		
		}
		delete(fields, "gridLayoutProperty")
	} else {errs = append(errs, cog.MakeBuildErrors("gridLayoutProperty", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("GridLayoutWithoutValue", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `GridLayoutWithoutValue` objects.
func (resource GridLayoutWithoutValue) Equals(other GridLayoutWithoutValue) bool {
		if resource.Kind != other.Kind {
			return false
		}
		if resource.GridLayoutProperty != other.GridLayoutProperty {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `GridLayoutWithoutValue` fields for violations and returns them.
func (resource GridLayoutWithoutValue) Validate() error {
	return nil
}


type RowsLayoutWithoutValue struct {
    Kind string `json:"kind"`
    RowsLayoutProperty string `json:"rowsLayoutProperty"`
}

// NewRowsLayoutWithoutValue creates a new RowsLayoutWithoutValue object.
func NewRowsLayoutWithoutValue() *RowsLayoutWithoutValue {
	return &RowsLayoutWithoutValue{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `RowsLayoutWithoutValue` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *RowsLayoutWithoutValue) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "kind"
	if fields["kind"] != nil {
		if string(fields["kind"]) != "null" {
			if err := json.Unmarshal(fields["kind"], &resource.Kind); err != nil {
				errs = append(errs, cog.MakeBuildErrors("kind", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is null"))...)
		
		}
		delete(fields, "kind")
	} else {errs = append(errs, cog.MakeBuildErrors("kind", errors.New("required field is missing from input"))...)
	}
	// Field "rowsLayoutProperty"
	if fields["rowsLayoutProperty"] != nil {
		if string(fields["rowsLayoutProperty"]) != "null" {
			if err := json.Unmarshal(fields["rowsLayoutProperty"], &resource.RowsLayoutProperty); err != nil {
				errs = append(errs, cog.MakeBuildErrors("rowsLayoutProperty", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("rowsLayoutProperty", errors.New("required field is null"))...)
		
		}
		delete(fields, "rowsLayoutProperty")
	} else {errs = append(errs, cog.MakeBuildErrors("rowsLayoutProperty", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("RowsLayoutWithoutValue", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `RowsLayoutWithoutValue` objects.
func (resource RowsLayoutWithoutValue) Equals(other RowsLayoutWithoutValue) bool {
		if resource.Kind != other.Kind {
			return false
		}
		if resource.RowsLayoutProperty != other.RowsLayoutProperty {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `RowsLayoutWithoutValue` fields for violations and returns them.
func (resource RowsLayoutWithoutValue) Validate() error {
	return nil
}
