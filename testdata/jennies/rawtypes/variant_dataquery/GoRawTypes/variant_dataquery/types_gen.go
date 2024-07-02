package variant_dataquery

import (
	cogvariants "github.com/grafana/cog/generated/cog/variants"
)

type Query struct {
	Expr string `json:"expr"`
	Instant *bool `json:"instant,omitempty"`
}
func (resource Query) ImplementsDataqueryVariant() {}


func VariantConfig() cogvariants.DataqueryConfig {
	return cogvariants.DataqueryConfig{
		Identifier: "prometheus",
	    DataqueryUnmarshaler: func (raw []byte) (cogvariants.Dataquery, error) {
            dataquery := Query{}

            if err := json.Unmarshal(raw, &dataquery); err != nil {
                return nil, err
            }

            return dataquery, nil
        },
	}
}


