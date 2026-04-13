import typing


class Struct:
    a: typing.Optional['MyObject']
    b: typing.Optional['MyObject']
    c: typing.Optional[str]
    d: typing.Optional[list[str]]
    e: dict[str, typing.Optional[str]]
    f: typing.Optional['NullableFieldsStructF']
    g: str

    def __init__(self, a: typing.Optional['MyObject'] = None, b: typing.Optional['MyObject'] = None, c: typing.Optional[str] = None, d: typing.Optional[list[str]] = None, e: typing.Optional[dict[str, typing.Optional[str]]] = None, f: typing.Optional['NullableFieldsStructF'] = None) -> None:
        self.a = a
        self.b = b
        self.c = c
        self.d = d
        self.e = e if e is not None else {}
        self.f = f
        self.g = ConstantRef

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "a": self.a,
            "c": self.c,
            "d": self.d,
            "e": self.e,
            "f": self.f,
            "g": self.g,
        }
        if self.b is not None:
            payload["b"] = self.b
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "a" in data:
            args["a"] = MyObject.from_json(data["a"])
        if "b" in data:
            args["b"] = MyObject.from_json(data["b"])
        if "c" in data:
            args["c"] = data["c"]
        if "d" in data:
            args["d"] = data["d"]
        if "e" in data:
            args["e"] = data["e"]
        if "f" in data:
            args["f"] = NullableFieldsStructF.from_json(data["f"])        

        return cls(**args)


class MyObject:
    field: str

    def __init__(self, field: str = "") -> None:
        self.field = field

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "field": self.field,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "field" in data:
            args["field"] = data["field"]        

        return cls(**args)


ConstantRef: typing.Literal["hey"] = "hey"


class NullableFieldsStructF:
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



