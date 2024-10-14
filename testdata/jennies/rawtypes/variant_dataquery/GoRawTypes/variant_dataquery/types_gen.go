package variant_dataquery

import (
	variants "github.com/grafana/cog/generated/cog/variants"
	"encoding/json"
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
	}
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


