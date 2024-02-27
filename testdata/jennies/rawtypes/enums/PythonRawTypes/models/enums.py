import enum


class Operator(enum.StrEnum):
    """
    This is a very interesting string enum.
    """

    GREATER_THAN = ">"
    LESS_THAN = "<"


class TableSortOrder(enum.StrEnum):
    ASC = "asc"
    DESC = "desc"


class LogsSortOrder(enum.StrEnum):
    ASC = "time_asc"
    DESC = "time_desc"


class DashboardCursorSync(enum.IntEnum):
    """
    0 for no shared crosshair or tooltip (default).
    1 for shared crosshair.
    2 for shared crosshair AND shared tooltip.
    """

    OFF = 0
    CROSSHAIR = 1
    TOOLTIP = 2



