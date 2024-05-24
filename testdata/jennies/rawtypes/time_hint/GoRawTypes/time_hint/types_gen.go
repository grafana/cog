package time_hint

type ObjTime time.Time

type ObjWithTimeField struct {
	RegisteredAt time.Time `json:"registeredAt"`
}

