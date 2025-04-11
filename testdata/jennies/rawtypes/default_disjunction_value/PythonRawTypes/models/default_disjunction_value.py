import typing


DisjunctionClasses: typing.TypeAlias = typing.Union['ValueA', 'ValueB', 'ValueC']


class ValueA:
    type_val: typing.Literal["A"]
    an_array: list[str]
    other_ref: 'ValueB'

    def __init__(self, an_array: typing.Optional[list[str]] = None, other_ref: typing.Optional['ValueB'] = None):
        self.type_val = "A"
        self.an_array = an_array if an_array is not None else []
        self.other_ref = other_ref if other_ref is not None else ValueB()

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "type": self.type_val,
            "anArray": self.an_array,
            "otherRef": self.other_ref,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "anArray" in data:
            args["an_array"] = data["anArray"]
        if "otherRef" in data:
            args["other_ref"] = ValueB.from_json(data["otherRef"])        

        return cls(**args)


class ValueB:
    type_val: typing.Literal["B"]
    a_map: dict[str, int]
    def_val: typing.Union[typing.Literal[1], typing.Literal["a"], bool]

    def __init__(self, a_map: typing.Optional[dict[str, int]] = None, def_val: typing.Optional[typing.Union[typing.Literal[1], typing.Literal["a"], bool]] = None):
        self.type_val = "B"
        self.a_map = a_map if a_map is not None else {}
        self.def_val = def_val if def_val is not None else 1

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "type": self.type_val,
            "aMap": self.a_map,
            "def": self.def_val,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "aMap" in data:
            args["a_map"] = data["aMap"]
        if "def" in data:
            args["def_val"] = data["def"]        

        return cls(**args)


class ValueC:
    type_val: typing.Literal["C"]
    other: float

    def __init__(self, other: float = 0):
        self.type_val = "C"
        self.other = other

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "type": self.type_val,
            "other": self.other,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "other" in data:
            args["other"] = data["other"]        

        return cls(**args)


DisjunctionConstants: typing.TypeAlias = typing.Union[typing.Literal["abc"], typing.Literal[1], typing.Literal[True]]
