package time_hint

import (
	"time"
	cog "github.com/grafana/cog/generated/cog"
)

type ObjTime time.Time

type ObjWithTimeField struct {
	RegisteredAt time.Time `json:"registeredAt"`
}

func (resource ObjWithTimeField) Equals(other ObjWithTimeField) bool {
		if resource.RegisteredAt != other.RegisteredAt {
			return false
		}

	return true
}


func (resource ObjWithTimeField) Validate() error {
	var errs cog.BuildErrors

	if len(errs) == 0 {
		return nil
	}

	return errs
}


