package open_struct

import (
	"encoding/json"
	"errors"
	cog "github.com/grafana/cog/generated/cog"
	"fmt"
)

type OpenStruct struct {
    A string `json:"a"`
    B int64 `json:"b"`

    ExtraFields map[string]any `json:"-"`
}

// NewOpenStruct creates a new OpenStruct object.
func NewOpenStruct() *OpenStruct {
	return &OpenStruct{
}
}
// MarshalJSON implements a custom JSON marshalling logic to encode `OpenStruct` as JSON.
func (resource *OpenStruct) MarshalJSON() ([]byte, error) {
	type Alias *OpenStruct
	base, err := json.Marshal(Alias(resource))
	if err != nil {
		return nil, err
	}
	
	var baseMap map[string]any
	if err := json.Unmarshal(base, &baseMap); err != nil {
	    return nil, err
	}
	
	for k, v := range resource.ExtraFields {
		baseMap[k] = v
	}
	
	return json.Marshal(baseMap)
}

// UnmarshalJSON implements a custom JSON unmarshalling logic to decode `OpenStruct` from JSON.
func (resource *OpenStruct) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

    var data map[string]json.RawMessage
    if err := json.Unmarshal(raw, &data); err != nil {
        return err
    }
    // A
    if v, ok := data["a"]; ok {
        if err := json.Unmarshal(v, &resource.A); err != nil {
            return err
        }
        delete(data, "a")
    }
    // B
    if v, ok := data["b"]; ok {
        if err := json.Unmarshal(v, &resource.B); err != nil {
            return err
        }
        delete(data, "b")
    }

    resource.ExtraFields = make(map[string]any)
    for key, value := range data {
		var v any
		if err := json.Unmarshal(value, &v); err != nil {
			return err
		}
        resource.ExtraFields[key] = v
    }
    
	return nil
}

// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `OpenStruct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *OpenStruct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "a"
	if fields["a"] != nil {
		if string(fields["a"]) != "null" {
			if err := json.Unmarshal(fields["a"], &resource.A); err != nil {
				errs = append(errs, cog.MakeBuildErrors("a", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("a", errors.New("required field is null"))...)
		
		}
		delete(fields, "a")
	} else {errs = append(errs, cog.MakeBuildErrors("a", errors.New("required field is missing from input"))...)
	}
	// Field "b"
	if fields["b"] != nil {
		if string(fields["b"]) != "null" {
			if err := json.Unmarshal(fields["b"], &resource.B); err != nil {
				errs = append(errs, cog.MakeBuildErrors("b", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("b", errors.New("required field is null"))...)
		
		}
		delete(fields, "b")
	} else {errs = append(errs, cog.MakeBuildErrors("b", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("OpenStruct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `OpenStruct` objects.
func (resource OpenStruct) Equals(other OpenStruct) bool {
		if resource.A != other.A {
			return false
		}
		if resource.B != other.B {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `OpenStruct` fields for violations and returns them.
func (resource OpenStruct) Validate() error {
	return nil
}
