import typing


LayoutWithValue: typing.TypeAlias = typing.Union['GridLayoutUsingValue', 'RowsLayoutUsingValue']


class GridLayoutUsingValue:
    kind: typing.Literal["GridLayout"]
    grid_layout_property: str

    def __init__(self, kind: typing.Optional[typing.Literal["GridLayout"]] = None, grid_layout_property: str = ""):
        self.kind = kind if kind is not None else GridLayoutKindType
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
        
        if "kind" in data:
            args["kind"] = data["kind"]
        if "gridLayoutProperty" in data:
            args["grid_layout_property"] = data["gridLayoutProperty"]        

        return cls(**args)


GridLayoutKindType: typing.Literal["GridLayout"] = "GridLayout"


class RowsLayoutUsingValue:
    kind: typing.Literal["RowsLayout"]
    rows_layout_property: str

    def __init__(self, kind: typing.Optional[typing.Literal["RowsLayout"]] = None, rows_layout_property: str = ""):
        self.kind = kind if kind is not None else RowsLayoutKindType
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
        
        if "kind" in data:
            args["kind"] = data["kind"]
        if "rowsLayoutProperty" in data:
            args["rows_layout_property"] = data["rowsLayoutProperty"]        

        return cls(**args)


RowsLayoutKindType: typing.Literal["RowsLayout"] = "RowsLayout"


LayoutWithoutValue: typing.TypeAlias = typing.Union['GridLayoutWithoutValue', 'RowsLayoutWithoutValue']


class GridLayoutWithoutValue:
    kind: typing.Literal["GridLayout"]
    grid_layout_property: str

    def __init__(self, kind: typing.Optional[typing.Literal["GridLayout"]] = None, grid_layout_property: str = ""):
        self.kind = kind if kind is not None else GridLayoutKindType
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
        
        if "kind" in data:
            args["kind"] = data["kind"]
        if "gridLayoutProperty" in data:
            args["grid_layout_property"] = data["gridLayoutProperty"]        

        return cls(**args)


class RowsLayoutWithoutValue:
    kind: typing.Literal["RowsLayout"]
    rows_layout_property: str

    def __init__(self, kind: typing.Optional[typing.Literal["RowsLayout"]] = None, rows_layout_property: str = ""):
        self.kind = kind if kind is not None else RowsLayoutKindType
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
        
        if "kind" in data:
            args["kind"] = data["kind"]
        if "rowsLayoutProperty" in data:
            args["rows_layout_property"] = data["rowsLayoutProperty"]        

        return cls(**args)
