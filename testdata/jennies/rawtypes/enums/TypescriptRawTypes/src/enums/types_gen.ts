// This is a very interesting string enum.
export enum Operator {
	GreaterThan = ">",
	LessThan = "<",
}

export const defaultOperator = (): Operator => (Operator.GreaterThan);

export enum TableSortOrder {
	Asc = "asc",
	Desc = "desc",
}

export const defaultTableSortOrder = (): TableSortOrder => (TableSortOrder.Asc);

export enum LogsSortOrder {
	Asc = "time_asc",
	Desc = "time_desc",
}

export const defaultLogsSortOrder = (): LogsSortOrder => (LogsSortOrder.Asc);

// 0 for no shared crosshair or tooltip (default).
// 1 for shared crosshair.
// 2 for shared crosshair AND shared tooltip.
export enum DashboardCursorSync {
	Off = 0,
	Crosshair = 1,
	Tooltip = 2,
}

export const defaultDashboardCursorSync = (): DashboardCursorSync => (DashboardCursorSync.Off);

