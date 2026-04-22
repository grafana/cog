import typing


class SomeStruct:
    field_any: object

    def __init__(self, field_any: object = None) -> None:
        self.field_any = field_any

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldAny": self.field_any,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "FieldAny" in data:
            args["field_any"] = data["FieldAny"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, SomeStruct):
            return False
        if self.field_any != other.field_any:
            return False
        return True


# Refresh rate or disabled.
RefreshRate: typing.TypeAlias = typing.Union[str, bool]



