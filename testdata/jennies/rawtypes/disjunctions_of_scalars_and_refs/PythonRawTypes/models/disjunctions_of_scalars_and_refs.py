import typing


DisjunctionOfScalarsAndRefs: typing.TypeAlias = typing.Union[typing.Literal["a"], bool, list[str], 'MyRefA', 'MyRefB']


class MyRefA:
    foo: str

    def __init__(self, foo: str = "") -> None:
        self.foo = foo

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "foo": self.foo,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "foo" in data:
            args["foo"] = data["foo"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, MyRefA):
            return False
        if self.foo != other.foo:
            return False
        return True


class MyRefB:
    bar: int

    def __init__(self, bar: int = 0) -> None:
        self.bar = bar

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "bar": self.bar,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "bar" in data:
            args["bar"] = data["bar"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, MyRefB):
            return False
        if self.bar != other.bar:
            return False
        return True



