package dashboard

import (
	"reflect"
	variants "github.com/grafana/cog/generated/cog/variants"
	cog "github.com/grafana/cog/generated/cog"
	"encoding/json"
)

type Dashboard struct {
	Title string `json:"title"`
	Panels []Panel `json:"panels,omitempty"`
}

func (resource Dashboard) Equals(other Dashboard) bool {
		if resource.Title != other.Title {
			return false
		}

		if len(resource.Panels) != len(other.Panels) {
			return false
		}

		for i1 := range resource.Panels {
		if !resource.Panels[i1].Equals(other.Panels[i1]) {
			return false
		}
		}

	return true
}


type DataSourceRef struct {
	Type *string `json:"type,omitempty"`
	Uid *string `json:"uid,omitempty"`
}

func (resource DataSourceRef) Equals(other DataSourceRef) bool {
		if resource.Type == nil && other.Type != nil || resource.Type != nil && other.Type == nil {
			return false
		}

		if resource.Type != nil {
		if *resource.Type != *other.Type {
			return false
		}
		}
		if resource.Uid == nil && other.Uid != nil || resource.Uid != nil && other.Uid == nil {
			return false
		}

		if resource.Uid != nil {
		if *resource.Uid != *other.Uid {
			return false
		}
		}

	return true
}


type FieldConfigSource struct {
	Defaults *FieldConfig `json:"defaults,omitempty"`
}

func (resource FieldConfigSource) Equals(other FieldConfigSource) bool {
		if resource.Defaults == nil && other.Defaults != nil || resource.Defaults != nil && other.Defaults == nil {
			return false
		}

		if resource.Defaults != nil {
		if !resource.Defaults.Equals(*other.Defaults) {
			return false
		}
		}

	return true
}


type FieldConfig struct {
	Unit *string `json:"unit,omitempty"`
	Custom any `json:"custom,omitempty"`
}

func (resource FieldConfig) Equals(other FieldConfig) bool {
		if resource.Unit == nil && other.Unit != nil || resource.Unit != nil && other.Unit == nil {
			return false
		}

		if resource.Unit != nil {
		if *resource.Unit != *other.Unit {
			return false
		}
		}
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.Custom, other.Custom) {
			return false
		}

	return true
}


type Panel struct {
	Title string `json:"title"`
	Type string `json:"type"`
	Datasource *DataSourceRef `json:"datasource,omitempty"`
	Options any `json:"options,omitempty"`
	Targets []variants.Dataquery `json:"targets,omitempty"`
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

func (resource Panel) Equals(other Panel) bool {
		if resource.Title != other.Title {
			return false
		}
		if resource.Type != other.Type {
			return false
		}
		if resource.Datasource == nil && other.Datasource != nil || resource.Datasource != nil && other.Datasource == nil {
			return false
		}

		if resource.Datasource != nil {
		if !resource.Datasource.Equals(*other.Datasource) {
			return false
		}
		}
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.Options, other.Options) {
			return false
		}

		if len(resource.Targets) != len(other.Targets) {
			return false
		}

		for i1 := range resource.Targets {
		if !resource.Targets[i1].Equals(other.Targets[i1]) {
			return false
		}
		}
		if resource.FieldConfig == nil && other.FieldConfig != nil || resource.FieldConfig != nil && other.FieldConfig == nil {
			return false
		}

		if resource.FieldConfig != nil {
		if !resource.FieldConfig.Equals(*other.FieldConfig) {
			return false
		}
		}

	return true
}


