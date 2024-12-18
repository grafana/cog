package variant_dataquery

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type Query struct {
    Expr string `json:"expr"`
    Instant *bool `json:"instant,omitempty"`
}
func (resource Query) ImplementsDataqueryVariant() {}


// NewQuery creates a new Query object.
func NewQuery() *Query {
	return &Query{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Query` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Query) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "expr"
	if fields["expr"] != nil {
		if string(fields["expr"]) != "null" {
			if err := json.Unmarshal(fields["expr"], &resource.Expr); err != nil {
				errs = append(errs, cog.MakeBuildErrors("expr", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("expr", errors.New("required field is null"))...)
		
		}
		delete(fields, "expr")
	} else {errs = append(errs, cog.MakeBuildErrors("expr", errors.New("required field is missing from input"))...)
	}
	// Field "instant"
	if fields["instant"] != nil {
		if string(fields["instant"]) != "null" {
			if err := json.Unmarshal(fields["instant"], &resource.Instant); err != nil {
				errs = append(errs, cog.MakeBuildErrors("instant", err)...)
			}
		
		}
		delete(fields, "instant")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Query", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two dataqueries.
func (resource Query) Equals(otherCandidate variants.Dataquery) bool {
	if otherCandidate == nil {
		return false
	}

	other, ok := otherCandidate.(Query)
	if !ok {
		return false
	}
		if resource.Expr != other.Expr {
			return false
		}
		if resource.Instant == nil && other.Instant != nil || resource.Instant != nil && other.Instant == nil {
			return false
		}

		if resource.Instant != nil {
		if *resource.Instant != *other.Instant {
			return false
		}
		}

	return true
}
// Validate checks all the validation constraints that may be defined on `Query` fields for violations and returns them.
func (resource Query) Validate() error {
	return nil
}


