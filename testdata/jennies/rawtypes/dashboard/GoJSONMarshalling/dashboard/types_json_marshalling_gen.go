package dashboard
import (
	cog "github.com/grafana/cog/generated/cog"
)

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
			customFieldConfig, err := variantCfg.FieldConfigUnmarshaler(fakeFieldConfigSource.Defaults.Custom)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(fields["fieldConfig"], &resource.FieldConfig); err != nil {
				return err
			}
	
			resource.FieldConfig.Defaults.Custom = customFieldConfig
		} else {
			if err := json.Unmarshal(fields["fieldConfig"], &resource.FieldConfig); err != nil {
				return err
			}
		}
	}

	dataqueryTypeHint := ""
if resource.Datasource != nil && resource.Datasource.Type != nil {
dataqueryTypeHint = *resource.Datasource.Type
}

	targets, err := cog.UnmarshalDataqueryArray(fields["targets"], dataqueryTypeHint)
	if err != nil {
		return err
	}
	resource.Targets = targets

	return nil
}
