package dashboard

import (
	cogvariants "github.com/grafana/cog/generated/cog/variants"
	cog "github.com/grafana/cog/generated/cog"
)

type Dashboard struct {
	Title string `json:"title"`
	Panels []Panel `json:"panels,omitempty"`
}

type DataSourceRef struct {
	Type *string `json:"type,omitempty"`
	Uid *string `json:"uid,omitempty"`
}

type FieldConfigSource struct {
	Defaults *FieldConfig `json:"defaults,omitempty"`
}

type FieldConfig struct {
	Unit *string `json:"unit,omitempty"`
	Custom any `json:"custom,omitempty"`
}

type Panel struct {
	Title string `json:"title"`
	Type string `json:"type"`
	Datasource *DataSourceRef `json:"datasource,omitempty"`
	Options any `json:"options,omitempty"`
	Targets []cogvariants.Dataquery `json:"targets,omitempty"`
	FieldConfig *FieldConfigSource `json:"fieldConfig,omitempty"`
}

func (resource *Panel) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}
	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	
	if fields["title"] != nil {
		if err := json.Unmarshal(fields["title"], &resource.Title); err != nil {
			return err
		}
	}

	if fields["type"] != nil {
		if err := json.Unmarshal(fields["type"], &resource.Type); err != nil {
			return err
		}
	}

	if fields["datasource"] != nil {
		if err := json.Unmarshal(fields["datasource"], &resource.Datasource); err != nil {
			return err
		}
	}

	if fields["options"] != nil {
		variantCfg, found := cog.ConfigForPanelcfgVariant(resource.Type)
		if found && variantCfg.OptionsUnmarshaler != nil {
			options, err := variantCfg.OptionsUnmarshaler(fields["options"])
			if err != nil {
				return err
			}
			resource.Options = options
		} else {
			if err := json.Unmarshal(fields["options"], &resource.Options); err != nil {
				return err
			}
		}
	}

	if fields["fieldConfig"] != nil {
		if err := json.Unmarshal(fields["fieldConfig"], &resource.FieldConfig); err != nil {
			return err
		}

		variantCfg, found := cog.ConfigForPanelcfgVariant(resource.Type)
		if found && variantCfg.FieldConfigUnmarshaler != nil {
			fakeFieldConfigSource := struct{
				Defaults struct {
					Custom json.RawMessage `json:"custom"` 
				} `json:"defaults"`
			}{}
			if err := json.Unmarshal(fields["fieldConfig"], &fakeFieldConfigSource); err != nil {
				return err
			}

			if fakeFieldConfigSource.Defaults.Custom != nil {
				customFieldConfig, err := variantCfg.FieldConfigUnmarshaler(fakeFieldConfigSource.Defaults.Custom)
				if err != nil {
					return err
				}

				resource.FieldConfig.Defaults.Custom = customFieldConfig
			}
		}
	}

	dataqueryTypeHint := ""
if resource.Datasource != nil && resource.Datasource.Type != nil {
dataqueryTypeHint = *resource.Datasource.Type
}

	if fields["targets"] != nil {
		targets, err := cog.UnmarshalDataqueryArray(fields["targets"], dataqueryTypeHint)
		if err != nil {
			return err
		}
		resource.Targets = targets
	}

	return nil
}

