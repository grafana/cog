package variantpanelcfgonlyoptions
import (
	cogvariants "github.com/grafana/cog/generated/cog/variants"
)

func VariantConfig() cogvariants.PanelcfgConfig {
	return cogvariants.PanelcfgConfig{
		Identifier: "text",
		OptionsUnmarshaler: func (raw []byte) (any, error) {
			options := Options{}

			if err := json.Unmarshal(raw, &options); err != nil {
				return nil, err
			}

			return options, nil
		},
	}
}

