import typing


class SomeStruct:
    """
    This
    is
    a
    comment
    """

    # Anything can go in there.
    # Really, anything.
    field_any: object
    field_bool: bool
    field_bytes: bytes
    field_string: str
    field_string_with_constant_value: typing.Literal["auto"]
    field_float32: float
    field_float64: float
    field_uint8: int
    field_uint16: int
    field_uint32: int
    field_uint64: int
    field_int8: int
    field_int16: int
    field_int32: int
    field_int64: int

    def __init__(self, field_any: object = None, field_bool: bool = False, field_bytes: bytes = "", field_string: str = "", field_float32: float = 0, field_float64: float = 0, field_uint8: int = 0, field_uint16: int = 0, field_uint32: int = 0, field_uint64: int = 0, field_int8: int = 0, field_int16: int = 0, field_int32: int = 0, field_int64: int = 0):
        self.field_any = field_any
        self.field_bool = field_bool
        self.field_bytes = field_bytes
        self.field_string = field_string
        self.field_string_with_constant_value = "auto"
        self.field_float32 = field_float32
        self.field_float64 = field_float64
        self.field_uint8 = field_uint8
        self.field_uint16 = field_uint16
        self.field_uint32 = field_uint32
        self.field_uint64 = field_uint64
        self.field_int8 = field_int8
        self.field_int16 = field_int16
        self.field_int32 = field_int32
        self.field_int64 = field_int64

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldAny": self.field_any,
            "FieldBool": self.field_bool,
            "FieldBytes": self.field_bytes,
            "FieldString": self.field_string,
            "FieldStringWithConstantValue": self.field_string_with_constant_value,
            "FieldFloat32": self.field_float32,
            "FieldFloat64": self.field_float64,
            "FieldUint8": self.field_uint8,
            "FieldUint16": self.field_uint16,
            "FieldUint32": self.field_uint32,
            "FieldUint64": self.field_uint64,
            "FieldInt8": self.field_int8,
            "FieldInt16": self.field_int16,
            "FieldInt32": self.field_int32,
            "FieldInt64": self.field_int64,
        }
        return payload
