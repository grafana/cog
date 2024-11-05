package time_hint

import (
	"time"
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type ObjTime time.Time

type ObjWithTimeField struct {
    RegisteredAt time.Time `json:"registeredAt"`
}

// NewObjWithTimeField creates a new ObjWithTimeField object.
func NewObjWithTimeField() *ObjWithTimeField {
	return &ObjWithTimeField{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `ObjWithTimeField` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *ObjWithTimeField) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "registeredAt"
	if fields["registeredAt"] != nil {
		if string(fields["registeredAt"]) != "null" {
			if err := json.Unmarshal(fields["registeredAt"], &resource.RegisteredAt); err != nil {
				errs = append(errs, cog.MakeBuildErrors("registeredAt", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("registeredAt", errors.New("required field is null"))...)
		
		}
		delete(fields, "registeredAt")
	} else {errs = append(errs, cog.MakeBuildErrors("registeredAt", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("ObjWithTimeField", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `ObjWithTimeField` objects.
func (resource ObjWithTimeField) Equals(other ObjWithTimeField) bool {
		if resource.RegisteredAt != other.RegisteredAt {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `ObjWithTimeField` fields for violations and returns them.
func (resource ObjWithTimeField) Validate() error {
	return nil
}


