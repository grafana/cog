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
    
    def tags(self, tags: str) -> typing.Self:    
        if self._internal.tags is None:
            self._internal.tags = []
        
        self._internal.tags.append(tags)
    
        return self
    
