package enums

// This is a very interesting string enum.
type Operator string
const (
	OperatorGreaterThan Operator = ">"
	OperatorLessThan Operator = "<"
)


type TableSortOrder string
const (
	TableSortOrderAsc TableSortOrder = "asc"
	TableSortOrderDesc TableSortOrder = "desc"
)


type LogsSortOrder string
const (
	LogsSortOrderAsc LogsSortOrder = "time_asc"
	LogsSortOrderDesc LogsSortOrder = "time_desc"
)


// 0 for no shared crosshair or tooltip (default).
// 1 for shared crosshair.
// 2 for shared crosshair AND shared tooltip.
type DashboardCursorSync int8
const (
	DashboardCursorSyncOff DashboardCursorSync = 0
	DashboardCursorSyncCrosshair DashboardCursorSync = 1
	DashboardCursorSyncTooltip DashboardCursorSync = 2
)


