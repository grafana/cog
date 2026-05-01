import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):
    _internal: sandbox.SomeStruct

    def __init__(self) -> None:
        self._internal = sandbox.SomeStruct()

    def build(self) -> sandbox.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def data(self, key: sandbox.StringEnum, value: str) -> typing.Self:    
        if self._internal.data is None:
            self._internal.data = {}
        
        self._internal.data[key] = value
    
        return self
    


class SomeStructWithDefaultEnum(cogbuilder.Builder[sandbox.SomeStructWithDefaultEnum]):
    _internal: sandbox.SomeStructWithDefaultEnum

    def __init__(self) -> None:
        self._internal = sandbox.SomeStructWithDefaultEnum()

    def build(self) -> sandbox.SomeStructWithDefaultEnum:
        """
        Builds the object.
        """
        return self._internal    
    
    def data(self, key: sandbox.StringEnumWithDefault, value: str) -> typing.Self:    
        if self._internal.data is None:
            self._internal.data = {}
        
        self._internal.data[key] = value
    
        return self
    
