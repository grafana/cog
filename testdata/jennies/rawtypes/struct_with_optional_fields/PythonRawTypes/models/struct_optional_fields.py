import typing


class StructOptionalFieldsSomeStructFieldAnonymousStruct:
    field_any: object

    def __init__(self, field_any: object = None):
        self.field_any = field_any

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldAny": self.field_any,
        }
        return payload


class SomeStruct:
    field_ref: typing.Optional['SomeOtherStruct']
    field_string: typing.Optional[str]
    operator: typing.Optional[typing.Literal[">", "<"]]
    field_array_of_strings: typing.Optional[list[str]]
    field_anonymous_struct: typing.Optional['StructOptionalFieldsSomeStructFieldAnonymousStruct']

    def __init__(self, field_ref: typing.Optional['SomeOtherStruct'] = None, field_string: typing.Optional[str] = None, operator: typing.Optional[typing.Literal[">", "<"]] = None, field_array_of_strings: typing.Optional[list[str]] = None, field_anonymous_struct: typing.Optional['StructOptionalFieldsSomeStructFieldAnonymousStruct'] = None):
        self.field_ref = field_ref
        self.field_string = field_string
        self.operator = operator
        self.field_array_of_strings = field_array_of_strings
        self.field_anonymous_struct = field_anonymous_struct

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
        }
        if self.field_ref is not None:
            payload["FieldRef"] = self.field_ref.to_json()
        if self.field_string is not None:
            payload["FieldString"] = self.field_string
        if self.operator is not None:
            payload["Operator"] = self.operator
        if self.field_array_of_strings is not None:
            payload["FieldArrayOfStrings"] = self.field_array_of_strings
        if self.field_anonymous_struct is not None:
            payload["FieldAnonymousStruct"] = self.field_anonymous_struct.to_json()
        return payload


class SomeOtherStruct:
    field_any: object

    def __init__(self, field_any: object = None):
        self.field_any = field_any

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldAny": self.field_any,
        }
        return payload
