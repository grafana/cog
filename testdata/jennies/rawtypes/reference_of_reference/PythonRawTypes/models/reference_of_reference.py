import typing


class MyStruct:
    field: typing.Optional['OtherStruct']

    def __init__(self, field: typing.Optional['OtherStruct'] = None) -> None:
        self.field = field

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
        }
        if self.field is not None:
            payload["field"] = self.field
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "field" in data:
            args["field"] = OtherStruct.from_json(data["field"])        

        return cls(**args)


OtherStruct: typing.TypeAlias = 'AnotherStruct'


class AnotherStruct:
    a: str

    def __init__(self, a: str = "") -> None:
        self.a = a

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "a": self.a,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "a" in data:
            args["a"] = data["a"]        

        return cls(**args)



