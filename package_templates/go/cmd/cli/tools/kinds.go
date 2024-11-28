package tools

type Kind = string

const (
	KindDashboard = "dashboard"
	KindPanel     = "panel"
	KindQuery     = "query"
)

func KnownKinds() []Kind {
	return []Kind{
		KindDashboard,
		KindPanel,
		KindQuery,
	}
}
