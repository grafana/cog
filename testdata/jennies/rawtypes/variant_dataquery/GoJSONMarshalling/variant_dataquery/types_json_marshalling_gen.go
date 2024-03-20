package variant_dataquery
import (
	cogvariants "github.com/grafana/cog/generated/cog/variants"
)

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

