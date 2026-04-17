import enum
import typing


class StringEnum(enum.StrEnum):
    A = "a"
    B = "b"
    C = "c"


class StringEnumWithDefault(enum.StrEnum):
    A = "a"
    B = "b"
    C = "c"


class SomeStruct:
    data: dict['StringEnum', str]

    def __init__(self, data: typing.Optional[dict['StringEnum', str]] = None) -> None:
        self.data = data if data is not None else {}

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "data": self.data,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "data" in data:
            args["data"] = data["data"]        

        return cls(**args)


class SomeStructWithDefaultEnum:
    data: dict['StringEnumWithDefault', str]

    def __init__(self, data: typing.Optional[dict['StringEnumWithDefault', str]] = None) -> None:
        self.data = data if data is not None else {}

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "data": self.data,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "data" in data:
            args["data"] = data["data"]        

        return cls(**args)



