import typing


# List of tags, maybe?
ArrayOfStrings = list[str]


class SomeStruct:
    field_any: object

    def __init__(self, field_any: object = None):
        self.field_any = field_any

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldAny": self.field_any,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args = {
            "field_any": data["FieldAny"],
        }
        return cls(**args)


ArrayOfRefs = list['SomeStruct']


ArrayOfArrayOfNumbers = list[list[int]]



