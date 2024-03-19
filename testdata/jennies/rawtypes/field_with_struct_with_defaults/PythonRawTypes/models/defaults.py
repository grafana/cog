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
        args: dict[str, typing.Any] = {}
        
        if "stringVal" in data:
            args["string_val"] = data["stringVal"]
        if "intVal" in data:
            args["int_val"] = data["intVal"]        

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
            "allFields": self.all_fields,
            "partialFields": self.partial_fields,
            "emptyFields": self.empty_fields,
            "complexField": self.complex_field,
            "partialComplexField": self.partial_complex_field,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "allFields" in data:
            args["all_fields"] = NestedStruct.from_json(data["allFields"])
        if "partialFields" in data:
            args["partial_fields"] = NestedStruct.from_json(data["partialFields"])
        if "emptyFields" in data:
            args["empty_fields"] = NestedStruct.from_json(data["emptyFields"])
        if "complexField" in data:
            args["complex_field"] = DefaultsStructComplexField.from_json(data["complexField"])
        if "partialComplexField" in data:
            args["partial_complex_field"] = DefaultsStructPartialComplexField.from_json(data["partialComplexField"])        

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
        args: dict[str, typing.Any] = {}
        
        if "nestedVal" in data:
            args["nested_val"] = data["nestedVal"]        

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
            "nested": self.nested,
            "array": self.array,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "uid" in data:
            args["uid"] = data["uid"]
        if "nested" in data:
            args["nested"] = DefaultsStructComplexFieldNested.from_json(data["nested"])
        if "array" in data:
            args["array"] = data["array"]        

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
        args: dict[str, typing.Any] = {}
        
        if "uid" in data:
            args["uid"] = data["uid"]
        if "intVal" in data:
            args["int_val"] = data["intVal"]        

        return cls(**args)



