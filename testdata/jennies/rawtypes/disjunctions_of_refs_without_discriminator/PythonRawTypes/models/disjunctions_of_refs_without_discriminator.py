import typing


DisjunctionWithoutDiscriminator: typing.TypeAlias = typing.Union['TypeA', 'TypeB']


class TypeA:
    field_a: str

    def __init__(self, field_a: str = "") -> None:
        self.field_a = field_a

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "fieldA": self.field_a,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "fieldA" in data:
            args["field_a"] = data["fieldA"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, TypeA):
            return False
        if self.field_a != other.field_a:
            return False
        return True


class TypeB:
    field_b: int

    def __init__(self, field_b: int = 0) -> None:
        self.field_b = field_b

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "fieldB": self.field_b,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "fieldB" in data:
            args["field_b"] = data["fieldB"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, TypeB):
            return False
        if self.field_b != other.field_b:
            return False
        return True



