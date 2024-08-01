package variant_dataquery

import (
	variants "github.com/grafana/cog/generated/cog/variants"
)

type Query struct {
	Expr string `json:"expr"`
	Instant *bool `json:"instant,omitempty"`
}
func (resource Query) ImplementsDataqueryVariant() {}


func VariantConfig() variants.DataqueryConfig {
	return variants.DataqueryConfig{
		Identifier: "prometheus",
	    DataqueryUnmarshaler: func (raw []byte) (variants.Dataquery, error) {
            dataquery := Query{}

            if err := json.Unmarshal(raw, &dataquery); err != nil {
                return nil, err
            }

            return dataquery, nil
       },
	}
}


