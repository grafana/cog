package variant_panelcfg_only_options

import (
	variants "github.com/grafana/cog/generated/cog/variants"
)

type Options struct {
	Content string `json:"content"`
}

func (resource Options) Equals(other Options) bool {
		if resource.Content != other.Content {
			return false
		}

	return true
}


func VariantConfig() variants.PanelcfgConfig {
	return variants.PanelcfgConfig{
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

