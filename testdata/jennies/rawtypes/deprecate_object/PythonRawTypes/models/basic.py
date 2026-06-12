import warnings
import typing


@warnings.deprecated("This object is deprecated, use NewStruct instead.")
class SomeStruct:
    field_string: str

    def __init__(self, field_string: str = "") -> None:
        self.field_string = field_string

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldString": self.field_string,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "FieldString" in data:
            args["field_string"] = data["FieldString"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, SomeStruct):
            return False
        if self.field_string != other.field_string:
            return False
        return True
