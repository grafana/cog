import typing


class MyStruct:
    scalars: typing.Union[str, bool, float, int]
    same_kind: typing.Literal["a", "b", "c"]
    refs: typing.Union['StructA', 'StructB']
    mixed: typing.Union['StructA', str, int]

    def __init__(self, scalars: typing.Optional[typing.Union[str, bool, float, int]] = None, same_kind: typing.Optional[typing.Literal["a", "b", "c"]] = None, refs: typing.Optional[typing.Union['StructA', 'StructB']] = None, mixed: typing.Optional[typing.Union['StructA', str, int]] = None) -> None:
        self.scalars = scalars if scalars is not None else ""
        self.same_kind = same_kind if same_kind is not None else "a"
        self.refs = refs if refs is not None else StructA()
        self.mixed = mixed if mixed is not None else StructA()

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "scalars": self.scalars,
            "sameKind": self.same_kind,
            "refs": self.refs,
            "mixed": self.mixed,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "scalars" in data:
            args["scalars"] = data["scalars"]
        if "sameKind" in data:
            args["same_kind"] = data["sameKind"]
        if "refs" in data:
            args["refs"] = data["refs"]
        if "mixed" in data:
            args["mixed"] = data["mixed"]        

        return cls(**args)


class StructA:
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


class StructB:
    type_val: int

    def __init__(self, type_val: int = 0) -> None:
        self.type_val = type_val

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "type": self.type_val,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "type" in data:
            args["type_val"] = data["type"]        

        return cls(**args)



