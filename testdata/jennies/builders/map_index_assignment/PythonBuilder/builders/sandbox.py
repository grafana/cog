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
    
    def annotations(self, key: str, value: str) -> typing.Self:    
        if self._internal.annotations is None:
            self._internal.annotations = {}
        
        self._internal.annotations[key] = value
    
        return self
    
