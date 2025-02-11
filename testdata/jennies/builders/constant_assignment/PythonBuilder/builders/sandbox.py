import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):
    _internal: sandbox.SomeStruct

    def __init__(self):
        self._internal = sandbox.SomeStruct()

    def build(self) -> sandbox.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def editable(self) -> typing.Self:    
        self._internal.editable = True
    
        return self
    
    def readonly(self) -> typing.Self:    
        self._internal.editable = False
    
        return self
    
    def auto_refresh(self) -> typing.Self:    
        self._internal.auto_refresh = True
    
        return self
    
    def no_auto_refresh(self) -> typing.Self:    
        self._internal.auto_refresh = False
    
        return self
    
