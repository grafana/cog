package constant_reference_discriminator

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type LayoutWithValue = GridLayoutUsingValueOrRowsLayoutUsingValue

// NewLayoutWithValue creates a new LayoutWithValue object.
func NewLayoutWithValue() *LayoutWithValue {
	return NewGridLayoutUsingValueOrRowsLayoutUsingValue()
}
type GridLayoutUsingValue struct {
    Kind string `json:"kind"`
    GridLayoutProperty string `json:"gridLayoutProperty"`
}

// NewGridLayoutUsingValue creates a new GridLayoutUsingValue object.
func NewGridLayoutUsingValue() *GridLayoutUsingValue {
	return &GridLayoutUsingValue{
		Kind: GridLayoutKindType,
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


type RowsLayoutUsingValue struct {
    Kind string `json:"kind"`
    RowsLayoutProperty string `json:"rowsLayoutProperty"`
}

// NewRowsLayoutUsingValue creates a new RowsLayoutUsingValue object.
func NewRowsLayoutUsingValue() *RowsLayoutUsingValue {
	return &RowsLayoutUsingValue{
		Kind: RowsLayoutKindType,
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


type LayoutWithoutValue = GridLayoutWithoutValueOrRowsLayoutWithoutValue

// NewLayoutWithoutValue creates a new LayoutWithoutValue object.
func NewLayoutWithoutValue() *LayoutWithoutValue {
	return NewGridLayoutWithoutValueOrRowsLayoutWithoutValue()
}
type GridLayoutWithoutValue struct {
    Kind string `json:"kind"`
    GridLayoutProperty string `json:"gridLayoutProperty"`
}

// NewGridLayoutWithoutValue creates a new GridLayoutWithoutValue object.
func NewGridLayoutWithoutValue() *GridLayoutWithoutValue {
	return &GridLayoutWithoutValue{
		Kind: GridLayoutKindType,
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
		Kind: RowsLayoutKindType,
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


const GridLayoutKindType = "GridLayout"

const RowsLayoutKindType = "RowsLayout"

type GridLayoutUsingValueOrRowsLayoutUsingValue struct {
    GridLayoutUsingValue *GridLayoutUsingValue `json:"GridLayoutUsingValue,omitempty"`
    RowsLayoutUsingValue *RowsLayoutUsingValue `json:"RowsLayoutUsingValue,omitempty"`
}

// NewGridLayoutUsingValueOrRowsLayoutUsingValue creates a new GridLayoutUsingValueOrRowsLayoutUsingValue object.
func NewGridLayoutUsingValueOrRowsLayoutUsingValue() *GridLayoutUsingValueOrRowsLayoutUsingValue {
	return &GridLayoutUsingValueOrRowsLayoutUsingValue{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `GridLayoutUsingValueOrRowsLayoutUsingValue` as JSON.
func (resource GridLayoutUsingValueOrRowsLayoutUsingValue) MarshalJSON() ([]byte, error) {
	if resource.GridLayoutUsingValue != nil {
		return json.Marshal(resource.GridLayoutUsingValue)
	}
	if resource.RowsLayoutUsingValue != nil {
		return json.Marshal(resource.RowsLayoutUsingValue)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `GridLayoutUsingValueOrRowsLayoutUsingValue` from JSON.
func (resource *GridLayoutUsingValueOrRowsLayoutUsingValue) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["kind"]
	if !found {
		return nil
	}

	switch discriminator {
	case "GridLayout":
		var gridLayoutUsingValue GridLayoutUsingValue
		if err := json.Unmarshal(raw, &gridLayoutUsingValue); err != nil {
			return err
		}

		resource.GridLayoutUsingValue = &gridLayoutUsingValue
		return nil
	case "RowsLayout":
		var rowsLayoutUsingValue RowsLayoutUsingValue
		if err := json.Unmarshal(raw, &rowsLayoutUsingValue); err != nil {
			return err
		}

		resource.RowsLayoutUsingValue = &rowsLayoutUsingValue
		return nil
	}

	return nil
}


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `GridLayoutUsingValueOrRowsLayoutUsingValue` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *GridLayoutUsingValueOrRowsLayoutUsingValue) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["kind"]
	if !found {
		return fmt.Errorf("discriminator field 'kind' not found in payload")
	}

	switch discriminator {
		case "GridLayout":
		gridLayoutUsingValue := &GridLayoutUsingValue{}
		if err := gridLayoutUsingValue.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.GridLayoutUsingValue = gridLayoutUsingValue
		return nil
		case "RowsLayout":
		rowsLayoutUsingValue := &RowsLayoutUsingValue{}
		if err := rowsLayoutUsingValue.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.RowsLayoutUsingValue = rowsLayoutUsingValue
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `kind = %v`", discriminator)
}

// Equals tests the equality of two `GridLayoutUsingValueOrRowsLayoutUsingValue` objects.
func (resource GridLayoutUsingValueOrRowsLayoutUsingValue) Equals(other GridLayoutUsingValueOrRowsLayoutUsingValue) bool {
		if resource.GridLayoutUsingValue == nil && other.GridLayoutUsingValue != nil || resource.GridLayoutUsingValue != nil && other.GridLayoutUsingValue == nil {
			return false
		}

		if resource.GridLayoutUsingValue != nil {
		if !resource.GridLayoutUsingValue.Equals(*other.GridLayoutUsingValue) {
			return false
		}
		}
		if resource.RowsLayoutUsingValue == nil && other.RowsLayoutUsingValue != nil || resource.RowsLayoutUsingValue != nil && other.RowsLayoutUsingValue == nil {
			return false
		}

		if resource.RowsLayoutUsingValue != nil {
		if !resource.RowsLayoutUsingValue.Equals(*other.RowsLayoutUsingValue) {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `GridLayoutUsingValueOrRowsLayoutUsingValue` fields for violations and returns them.
func (resource GridLayoutUsingValueOrRowsLayoutUsingValue) Validate() error {
	var errs cog.BuildErrors
		if resource.GridLayoutUsingValue != nil {
		if err := resource.GridLayoutUsingValue.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("GridLayoutUsingValue", err)...)
		}
		}
		if resource.RowsLayoutUsingValue != nil {
		if err := resource.RowsLayoutUsingValue.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("RowsLayoutUsingValue", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type GridLayoutWithoutValueOrRowsLayoutWithoutValue struct {
    GridLayoutWithoutValue *GridLayoutWithoutValue `json:"GridLayoutWithoutValue,omitempty"`
    RowsLayoutWithoutValue *RowsLayoutWithoutValue `json:"RowsLayoutWithoutValue,omitempty"`
}

// NewGridLayoutWithoutValueOrRowsLayoutWithoutValue creates a new GridLayoutWithoutValueOrRowsLayoutWithoutValue object.
func NewGridLayoutWithoutValueOrRowsLayoutWithoutValue() *GridLayoutWithoutValueOrRowsLayoutWithoutValue {
	return &GridLayoutWithoutValueOrRowsLayoutWithoutValue{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `GridLayoutWithoutValueOrRowsLayoutWithoutValue` as JSON.
func (resource GridLayoutWithoutValueOrRowsLayoutWithoutValue) MarshalJSON() ([]byte, error) {
	if resource.GridLayoutWithoutValue != nil {
		return json.Marshal(resource.GridLayoutWithoutValue)
	}
	if resource.RowsLayoutWithoutValue != nil {
		return json.Marshal(resource.RowsLayoutWithoutValue)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `GridLayoutWithoutValueOrRowsLayoutWithoutValue` from JSON.
func (resource *GridLayoutWithoutValueOrRowsLayoutWithoutValue) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["kind"]
	if !found {
		return nil
	}

	switch discriminator {
	case "GridLayout":
		var gridLayoutWithoutValue GridLayoutWithoutValue
		if err := json.Unmarshal(raw, &gridLayoutWithoutValue); err != nil {
			return err
		}

		resource.GridLayoutWithoutValue = &gridLayoutWithoutValue
		return nil
	case "RowsLayout":
		var rowsLayoutWithoutValue RowsLayoutWithoutValue
		if err := json.Unmarshal(raw, &rowsLayoutWithoutValue); err != nil {
			return err
		}

		resource.RowsLayoutWithoutValue = &rowsLayoutWithoutValue
		return nil
	}

	return nil
}


// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `GridLayoutWithoutValueOrRowsLayoutWithoutValue` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *GridLayoutWithoutValueOrRowsLayoutWithoutValue) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	// FIXME: this is wasteful, we need to find a more efficient way to unmarshal this.
	parsedAsMap := make(map[string]any)
	if err := json.Unmarshal(raw, &parsedAsMap); err != nil {
		return err
	}

	discriminator, found := parsedAsMap["kind"]
	if !found {
		return fmt.Errorf("discriminator field 'kind' not found in payload")
	}

	switch discriminator {
		case "GridLayout":
		gridLayoutWithoutValue := &GridLayoutWithoutValue{}
		if err := gridLayoutWithoutValue.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.GridLayoutWithoutValue = gridLayoutWithoutValue
		return nil
		case "RowsLayout":
		rowsLayoutWithoutValue := &RowsLayoutWithoutValue{}
		if err := rowsLayoutWithoutValue.UnmarshalJSONStrict(raw); err != nil {
			return err
		}

		resource.RowsLayoutWithoutValue = rowsLayoutWithoutValue
		return nil
	}

	return fmt.Errorf("could not unmarshal resource with `kind = %v`", discriminator)
}

// Equals tests the equality of two `GridLayoutWithoutValueOrRowsLayoutWithoutValue` objects.
func (resource GridLayoutWithoutValueOrRowsLayoutWithoutValue) Equals(other GridLayoutWithoutValueOrRowsLayoutWithoutValue) bool {
		if resource.GridLayoutWithoutValue == nil && other.GridLayoutWithoutValue != nil || resource.GridLayoutWithoutValue != nil && other.GridLayoutWithoutValue == nil {
			return false
		}

		if resource.GridLayoutWithoutValue != nil {
		if !resource.GridLayoutWithoutValue.Equals(*other.GridLayoutWithoutValue) {
			return false
		}
		}
		if resource.RowsLayoutWithoutValue == nil && other.RowsLayoutWithoutValue != nil || resource.RowsLayoutWithoutValue != nil && other.RowsLayoutWithoutValue == nil {
			return false
		}

		if resource.RowsLayoutWithoutValue != nil {
		if !resource.RowsLayoutWithoutValue.Equals(*other.RowsLayoutWithoutValue) {
			return false
		}
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `GridLayoutWithoutValueOrRowsLayoutWithoutValue` fields for violations and returns them.
func (resource GridLayoutWithoutValueOrRowsLayoutWithoutValue) Validate() error {
	var errs cog.BuildErrors
		if resource.GridLayoutWithoutValue != nil {
		if err := resource.GridLayoutWithoutValue.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("GridLayoutWithoutValue", err)...)
		}
		}
		if resource.RowsLayoutWithoutValue != nil {
		if err := resource.RowsLayoutWithoutValue.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("RowsLayoutWithoutValue", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}
