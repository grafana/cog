import enum
import typing


class Enum(enum.StrEnum):
    VALUE_A = "ValueA"
    VALUE_B = "ValueB"
    VALUE_C = "ValueC"


class ParentStruct:
    my_enum: 'Enum'

    def __init__(self, my_enum: typing.Optional['Enum'] = None):
        self.my_enum = my_enum if my_enum is not None else Enum.VALUE_A

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "myEnum": self.my_enum,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "myEnum" in data:
            args["my_enum"] = data["myEnum"]        

        return cls(**args)


class Struct:
    my_value: str
    my_enum: 'Enum'

    def __init__(self, my_value: str = "", my_enum: typing.Optional['Enum'] = None):
        self.my_value = my_value
        self.my_enum = my_enum if my_enum is not None else Enum.VALUE_A

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "myValue": self.my_value,
            "myEnum": self.my_enum,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "myValue" in data:
            args["my_value"] = data["myValue"]
        if "myEnum" in data:
            args["my_enum"] = data["myEnum"]        

        return cls(**args)


class StructA:
    my_enum: str

    def __init__(self, ):
        self.my_enum = Enum.VALUE_A

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "myEnum": self.my_enum,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        return cls(**args)


class StructB:
    my_enum: str
    my_value: str

    def __init__(self, my_value: str = ""):
        self.my_enum = Enum.VALUE_B
        self.my_value = my_value

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "myEnum": self.my_enum,
            "myValue": self.my_value,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "myValue" in data:
            args["my_value"] = data["myValue"]        

        return cls(**args)



