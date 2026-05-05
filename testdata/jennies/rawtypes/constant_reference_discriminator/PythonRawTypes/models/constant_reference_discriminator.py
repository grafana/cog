import typing


LayoutWithValue: typing.TypeAlias = typing.Union['GridLayoutUsingValue', 'RowsLayoutUsingValue']


class GridLayoutUsingValue:
    kind: str
    grid_layout_property: str

    def __init__(self, grid_layout_property: str = "") -> None:
        self.kind = GridLayoutKindType
        self.grid_layout_property = grid_layout_property

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "kind": self.kind,
            "gridLayoutProperty": self.grid_layout_property,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "gridLayoutProperty" in data:
            args["grid_layout_property"] = data["gridLayoutProperty"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, GridLayoutUsingValue):
            return False
        if self.kind != other.kind:
            return False
        if self.grid_layout_property != other.grid_layout_property:
            return False
        return True


class RowsLayoutUsingValue:
    kind: str
    rows_layout_property: str

    def __init__(self, rows_layout_property: str = "") -> None:
        self.kind = RowsLayoutKindType
        self.rows_layout_property = rows_layout_property

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "kind": self.kind,
            "rowsLayoutProperty": self.rows_layout_property,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "rowsLayoutProperty" in data:
            args["rows_layout_property"] = data["rowsLayoutProperty"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, RowsLayoutUsingValue):
            return False
        if self.kind != other.kind:
            return False
        if self.rows_layout_property != other.rows_layout_property:
            return False
        return True


LayoutWithoutValue: typing.TypeAlias = typing.Union['GridLayoutWithoutValue', 'RowsLayoutWithoutValue']


class GridLayoutWithoutValue:
    kind: str
    grid_layout_property: str

    def __init__(self, grid_layout_property: str = "") -> None:
        self.kind = GridLayoutKindType
        self.grid_layout_property = grid_layout_property

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "kind": self.kind,
            "gridLayoutProperty": self.grid_layout_property,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "gridLayoutProperty" in data:
            args["grid_layout_property"] = data["gridLayoutProperty"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, GridLayoutWithoutValue):
            return False
        if self.kind != other.kind:
            return False
        if self.grid_layout_property != other.grid_layout_property:
            return False
        return True


class RowsLayoutWithoutValue:
    kind: str
    rows_layout_property: str

    def __init__(self, rows_layout_property: str = "") -> None:
        self.kind = RowsLayoutKindType
        self.rows_layout_property = rows_layout_property

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "kind": self.kind,
            "rowsLayoutProperty": self.rows_layout_property,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "rowsLayoutProperty" in data:
            args["rows_layout_property"] = data["rowsLayoutProperty"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, RowsLayoutWithoutValue):
            return False
        if self.kind != other.kind:
            return False
        if self.rows_layout_property != other.rows_layout_property:
            return False
        return True


GridLayoutKindType: typing.Literal["GridLayout"] = "GridLayout"


RowsLayoutKindType: typing.Literal["RowsLayout"] = "RowsLayout"



