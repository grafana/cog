package time_hint

import (
	"time"
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


// Validate checks any constraint that may be defined for this type
// and returns all violations.
func (resource ObjWithTimeField) Validate() error {
	return nil
}


