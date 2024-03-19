import typing


class SomeStruct:
    field_bool: bool
    field_string: str
    field_string_with_constant_value: typing.Literal["auto"]
    field_float32: float
    field_int32: int

    def __init__(self, field_bool: bool = True, field_string: str = "foo", field_float32: float = 42.42, field_int32: int = 42):
        self.field_bool = field_bool
        self.field_string = field_string
        self.field_string_with_constant_value = "auto"
        self.field_float32 = field_float32
        self.field_int32 = field_int32

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "fieldBool": self.field_bool,
            "fieldString": self.field_string,
            "FieldStringWithConstantValue": self.field_string_with_constant_value,
            "FieldFloat32": self.field_float32,
            "FieldInt32": self.field_int32,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "fieldBool" in data:
            args["field_bool"] = data["fieldBool"]
        if "fieldString" in data:
            args["field_string"] = data["fieldString"]
        if "FieldFloat32" in data:
            args["field_float32"] = data["FieldFloat32"]
        if "FieldInt32" in data:
            args["field_int32"] = data["FieldInt32"]        

        return cls(**args)
