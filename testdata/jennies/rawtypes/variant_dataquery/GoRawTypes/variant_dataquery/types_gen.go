package variant_dataquery

import (
	variants "github.com/grafana/cog/generated/cog/variants"
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

func (resource Query) DataqueryType() string {
	return "prometheus"
}

func VariantConfig() variants.DataqueryConfig {
	return variants.DataqueryConfig{
		Identifier: "prometheus",
	    DataqueryUnmarshaler: func (raw []byte) (variants.Dataquery, error) {
            dataquery := &Query{}

            if err := json.Unmarshal(raw, dataquery); err != nil {
                return nil, err
            }

            return dataquery, nil
       },
	    StrictDataqueryUnmarshaler: func (raw []byte) (variants.Dataquery, error) {
            dataquery := &Query{}

            if err := dataquery.StrictUnmarshalJSON(raw); err != nil {
                return nil, err
            }

            return dataquery, nil
       },
	}
}


func (resource *Query) StrictUnmarshalJSON(raw []byte) error {
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
// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource Query) Validate() error {
	return nil
}


