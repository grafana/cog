import typing


# Refresh rate or disabled.
RefreshRate: typing.TypeAlias = typing.Union[str, bool]


StringOrNull: typing.TypeAlias = typing.Optional[str]


class SomeStruct:
    type: typing.Literal["some-struct"]
    field_any: object

    def __init__(self, field_any: object = None):
        self.type = "some-struct"
        self.field_any = field_any

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "Type": self.type,
            "FieldAny": self.field_any,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "FieldAny" in data:
            args["field_any"] = data["FieldAny"]        

        return cls(**args)


BoolOrRef: typing.TypeAlias = typing.Union[bool, 'SomeStruct']


class SomeOtherStruct:
    type: typing.Literal["some-other-struct"]
    foo: bytes

    def __init__(self, foo: bytes = ""):
        self.type = "some-other-struct"
        self.foo = foo

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "Type": self.type,
            "Foo": self.foo,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "Foo" in data:
            args["foo"] = data["Foo"]        

        return cls(**args)


class YetAnotherStruct:
    type: typing.Literal["yet-another-struct"]
    bar: int

    def __init__(self, bar: int = 0):
        self.type = "yet-another-struct"
        self.bar = bar

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "Type": self.type,
            "Bar": self.bar,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "Bar" in data:
            args["bar"] = data["Bar"]        

        return cls(**args)


SeveralRefs: typing.TypeAlias = typing.Union['SomeStruct', 'SomeOtherStruct', 'YetAnotherStruct']



