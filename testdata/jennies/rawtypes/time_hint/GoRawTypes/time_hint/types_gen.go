package time_hint

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


