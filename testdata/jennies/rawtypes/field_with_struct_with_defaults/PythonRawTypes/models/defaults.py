import typing


class NestedStruct:
    string_val: str
    int_val: int

    def __init__(self, string_val: str = "", int_val: int = 0):
        self.string_val = string_val
        self.int_val = int_val

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "stringVal": self.string_val,
            "intVal": self.int_val,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args = {
            "string_val": data["stringVal"],
            "int_val": data["intVal"],
        }
        return cls(**args)


class Struct:
    all_fields: 'NestedStruct'
    partial_fields: 'NestedStruct'
    empty_fields: 'NestedStruct'
    complex_field: 'DefaultsStructComplexField'
    partial_complex_field: 'DefaultsStructPartialComplexField'

    def __init__(self, all_fields: typing.Optional['NestedStruct'] = None, partial_fields: typing.Optional['NestedStruct'] = None, empty_fields: typing.Optional['NestedStruct'] = None, complex_field: typing.Optional['DefaultsStructComplexField'] = None, partial_complex_field: typing.Optional['DefaultsStructPartialComplexField'] = None):
        self.all_fields = all_fields if all_fields is not None else NestedStruct(int_val=3, string_val="hello")
        self.partial_fields = partial_fields if partial_fields is not None else NestedStruct(int_val=3)
        self.empty_fields = empty_fields if empty_fields is not None else NestedStruct()
        self.complex_field = complex_field if complex_field is not None else DefaultsStructComplexField(array=["hello"], nested=DefaultsStructComplexFieldNested(nested_val="nested"), uid="myUID")
        self.partial_complex_field = partial_complex_field if partial_complex_field is not None else DefaultsStructPartialComplexField()

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "allFields": None if self.all_fields is None else self.all_fields.to_json(),
            "partialFields": None if self.partial_fields is None else self.partial_fields.to_json(),
            "emptyFields": None if self.empty_fields is None else self.empty_fields.to_json(),
            "complexField": None if self.complex_field is None else self.complex_field.to_json(),
            "partialComplexField": None if self.partial_complex_field is None else self.partial_complex_field.to_json(),
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args = {
            "all_fields": NestedStruct.from_json(data["allFields"]),
            "partial_fields": NestedStruct.from_json(data["partialFields"]),
            "empty_fields": NestedStruct.from_json(data["emptyFields"]),
            "complex_field": DefaultsStructComplexField.from_json(data["complexField"]),
            "partial_complex_field": DefaultsStructPartialComplexField.from_json(data["partialComplexField"]),
        }
        return cls(**args)


class DefaultsStructComplexFieldNested:
    nested_val: str

    def __init__(self, nested_val: str = ""):
        self.nested_val = nested_val

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "nestedVal": self.nested_val,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args = {
            "nested_val": data["nestedVal"],
        }
        return cls(**args)


class DefaultsStructComplexField:
    uid: str
    nested: 'DefaultsStructComplexFieldNested'
    array: list[str]

    def __init__(self, uid: str = "", nested: typing.Optional['DefaultsStructComplexFieldNested'] = None, array: typing.Optional[list[str]] = None):
        self.uid = uid
        self.nested = nested if nested is not None else DefaultsStructComplexFieldNested()
        self.array = array if array is not None else []

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "uid": self.uid,
            "nested": None if self.nested is None else self.nested.to_json(),
            "array": self.array,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args = {
            "uid": data["uid"],
            "nested": DefaultsStructComplexFieldNested.from_json(data["nested"]),
            "array": data["array"],
        }
        return cls(**args)


class DefaultsStructPartialComplexField:
    uid: str
    int_val: int

    def __init__(self, uid: str = "", int_val: int = 0):
        self.uid = uid
        self.int_val = int_val

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "uid": self.uid,
            "intVal": self.int_val,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args = {
            "uid": data["uid"],
            "int_val": data["intVal"],
        }
        return cls(**args)



