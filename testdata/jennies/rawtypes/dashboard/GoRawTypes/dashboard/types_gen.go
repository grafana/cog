package dashboard

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"strconv"
	"fmt"
	"reflect"
	variants "github.com/grafana/cog/generated/cog/variants"
)

type Dashboard struct {
    Title string `json:"title"`
    Panels []Panel `json:"panels,omitempty"`
}

// NewDashboard creates a new Dashboard object.
func NewDashboard() *Dashboard {
	return &Dashboard{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Dashboard` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Dashboard) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "title"
	if fields["title"] != nil {
		if string(fields["title"]) != "null" {
			if err := json.Unmarshal(fields["title"], &resource.Title); err != nil {
				errs = append(errs, cog.MakeBuildErrors("title", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("title", errors.New("required field is null"))...)
		
		}
		delete(fields, "title")
	} else {errs = append(errs, cog.MakeBuildErrors("title", errors.New("required field is missing from input"))...)
	}
	// Field "panels"
	if fields["panels"] != nil {
		if string(fields["panels"]) != "null" {
			
			partialArray := []json.RawMessage{}
			if err := json.Unmarshal(fields["panels"], &partialArray); err != nil {
				return err
			}

			for i1 := range partialArray {
				var result1 Panel
				
			result1 = Panel{}
			if err := result1.UnmarshalJSONStrict(partialArray[i1]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("panels["+strconv.Itoa(i1)+"]", err)...)
			}
				resource.Panels = append(resource.Panels, result1)
			}
		
		}
		delete(fields, "panels")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Dashboard", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `Dashboard` objects.
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


// Validate checks all the validation constraints that may be defined on `Dashboard` fields for violations and returns them.
func (resource Dashboard) Validate() error {
	var errs cog.BuildErrors

		for i1 := range resource.Panels {
		if err := resource.Panels[i1].Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("panels["+strconv.Itoa(i1)+"]", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type DataSourceRef struct {
    Type *string `json:"type,omitempty"`
    Uid *string `json:"uid,omitempty"`
}

// NewDataSourceRef creates a new DataSourceRef object.
func NewDataSourceRef() *DataSourceRef {
	return &DataSourceRef{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `DataSourceRef` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *DataSourceRef) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "type"
	if fields["type"] != nil {
		if string(fields["type"]) != "null" {
			if err := json.Unmarshal(fields["type"], &resource.Type); err != nil {
				errs = append(errs, cog.MakeBuildErrors("type", err)...)
			}
		
		}
		delete(fields, "type")
	
	}
	// Field "uid"
	if fields["uid"] != nil {
		if string(fields["uid"]) != "null" {
			if err := json.Unmarshal(fields["uid"], &resource.Uid); err != nil {
				errs = append(errs, cog.MakeBuildErrors("uid", err)...)
			}
		
		}
		delete(fields, "uid")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("DataSourceRef", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `DataSourceRef` objects.
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


// Validate checks all the validation constraints that may be defined on `DataSourceRef` fields for violations and returns them.
func (resource DataSourceRef) Validate() error {
	return nil
}


type FieldConfigSource struct {
    Defaults *FieldConfig `json:"defaults,omitempty"`
}

// NewFieldConfigSource creates a new FieldConfigSource object.
func NewFieldConfigSource() *FieldConfigSource {
	return &FieldConfigSource{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `FieldConfigSource` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *FieldConfigSource) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "defaults"
	if fields["defaults"] != nil {
		if string(fields["defaults"]) != "null" {
			
			resource.Defaults = &FieldConfig{}
			if err := resource.Defaults.UnmarshalJSONStrict(fields["defaults"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("defaults", err)...)
			}
		
		}
		delete(fields, "defaults")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("FieldConfigSource", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `FieldConfigSource` objects.
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


// Validate checks all the validation constraints that may be defined on `FieldConfigSource` fields for violations and returns them.
func (resource FieldConfigSource) Validate() error {
	var errs cog.BuildErrors
		if resource.Defaults != nil {
		if err := resource.Defaults.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("defaults", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type FieldConfig struct {
    Unit *string `json:"unit,omitempty"`
    Custom any `json:"custom,omitempty"`
}

// NewFieldConfig creates a new FieldConfig object.
func NewFieldConfig() *FieldConfig {
	return &FieldConfig{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `FieldConfig` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *FieldConfig) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "unit"
	if fields["unit"] != nil {
		if string(fields["unit"]) != "null" {
			if err := json.Unmarshal(fields["unit"], &resource.Unit); err != nil {
				errs = append(errs, cog.MakeBuildErrors("unit", err)...)
			}
		
		}
		delete(fields, "unit")
	
	}
	// Field "custom"
	if fields["custom"] != nil {
		if string(fields["custom"]) != "null" {
			if err := json.Unmarshal(fields["custom"], &resource.Custom); err != nil {
				errs = append(errs, cog.MakeBuildErrors("custom", err)...)
			}
		
		}
		delete(fields, "custom")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("FieldConfig", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `FieldConfig` objects.
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


// Validate checks all the validation constraints that may be defined on `FieldConfig` fields for violations and returns them.
func (resource FieldConfig) Validate() error {
	return nil
}


type Panel struct {
    Title string `json:"title"`
    Type string `json:"type"`
    Datasource *DataSourceRef `json:"datasource,omitempty"`
    Options any `json:"options,omitempty"`
    Targets []variants.Dataquery `json:"targets,omitempty"`
    FieldConfig *FieldConfigSource `json:"fieldConfig,omitempty"`
}

// NewPanel creates a new Panel object.
func NewPanel() *Panel {
	return &Panel{
}
}
// UnmarshalJSON implements a custom JSON unmarshalling logic to decode Panel from JSON.
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
			return fmt.Errorf("error decoding field 'title': %w", err)
		}
	}

	if fields["type"] != nil {
		if err := json.Unmarshal(fields["type"], &resource.Type); err != nil {
			return fmt.Errorf("error decoding field 'type': %w", err)
		}
	}

	if fields["datasource"] != nil {
		if err := json.Unmarshal(fields["datasource"], &resource.Datasource); err != nil {
			return fmt.Errorf("error decoding field 'datasource': %w", err)
		}
	}

	if fields["options"] != nil {
		if err := json.Unmarshal(fields["options"], &resource.Options); err != nil {
			return fmt.Errorf("error decoding field 'options': %w", err)
		}
	}

	if fields["fieldConfig"] != nil {
		if err := json.Unmarshal(fields["fieldConfig"], &resource.FieldConfig); err != nil {
			return fmt.Errorf("error decoding field 'fieldConfig': %w", err)
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

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Panel` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Panel) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "title"
	if fields["title"] != nil {
		if string(fields["title"]) != "null" {
			if err := json.Unmarshal(fields["title"], &resource.Title); err != nil {
				errs = append(errs, cog.MakeBuildErrors("title", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("title", errors.New("required field is null"))...)
		
		}
		delete(fields, "title")
	} else {errs = append(errs, cog.MakeBuildErrors("title", errors.New("required field is missing from input"))...)
	}
	// Field "type"
	if fields["type"] != nil {
		if string(fields["type"]) != "null" {
			if err := json.Unmarshal(fields["type"], &resource.Type); err != nil {
				errs = append(errs, cog.MakeBuildErrors("type", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("type", errors.New("required field is null"))...)
		
		}
		delete(fields, "type")
	} else {errs = append(errs, cog.MakeBuildErrors("type", errors.New("required field is missing from input"))...)
	}
	// Field "datasource"
	if fields["datasource"] != nil {
		if string(fields["datasource"]) != "null" {
			
			resource.Datasource = &DataSourceRef{}
			if err := resource.Datasource.UnmarshalJSONStrict(fields["datasource"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("datasource", err)...)
			}
		
		}
		delete(fields, "datasource")
	
	}
	// Field "options"
	if fields["options"] != nil {
		if string(fields["options"]) != "null" {
			if err := json.Unmarshal(fields["options"], &resource.Options); err != nil {
				errs = append(errs, cog.MakeBuildErrors("options", err)...)
			}
		
		}
		delete(fields, "options")
	
	}
	// Field "targets"
	if fields["targets"] != nil {
		if string(fields["targets"]) != "null" {
			
			partialArray := []json.RawMessage{}
			if err := json.Unmarshal(fields["targets"], &partialArray); err != nil {
				return err
			}

			for i1 := range partialArray {
				var result1 variants.Dataquery
				
			dataquery, err := cog.StrictUnmarshalDataquery(partialArray[i1], "")
			if err != nil {
				errs = append(errs, cog.MakeBuildErrors("targets["+strconv.Itoa(i1)+"]", err)...)
			} else {
				result1 = dataquery
			}
				resource.Targets = append(resource.Targets, result1)
			}
		
		}
		delete(fields, "targets")
	
	}
	// Field "fieldConfig"
	if fields["fieldConfig"] != nil {
		if string(fields["fieldConfig"]) != "null" {
			
			resource.FieldConfig = &FieldConfigSource{}
			if err := resource.FieldConfig.UnmarshalJSONStrict(fields["fieldConfig"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("fieldConfig", err)...)
			}
		
		}
		delete(fields, "fieldConfig")
	
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Panel", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `Panel` objects.
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


// Validate checks all the validation constraints that may be defined on `Panel` fields for violations and returns them.
func (resource Panel) Validate() error {
	var errs cog.BuildErrors
		if resource.Datasource != nil {
		if err := resource.Datasource.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("datasource", err)...)
		}
		}

		for i1 := range resource.Targets {
		if err := resource.Targets[i1].Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("targets["+strconv.Itoa(i1)+"]", err)...)
		}
		}
		if resource.FieldConfig != nil {
		if err := resource.FieldConfig.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("fieldConfig", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


