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

func (resource *ObjWithTimeField) StrictUnmarshalJSON(raw []byte) error {
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


func (resource ObjWithTimeField) Equals(other ObjWithTimeField) bool {
		if resource.RegisteredAt != other.RegisteredAt {
			return false
		}

	return true
}


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource ObjWithTimeField) Validate() error {
	return nil
}


