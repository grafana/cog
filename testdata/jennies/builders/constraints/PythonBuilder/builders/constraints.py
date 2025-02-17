import typing
from ..cog import builder as cogbuilder
from ..models import constraints


class SomeStruct(cogbuilder.Builder[constraints.SomeStruct]):
    _internal: constraints.SomeStruct

    def __init__(self):
        self._internal = constraints.SomeStruct()

    def build(self) -> constraints.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def id(self, id_val: int) -> typing.Self:    
        if not id_val >= 5:
            raise ValueError("id_val must be >= 5")
        if not id_val < 10:
            raise ValueError("id_val must be < 10")
        self._internal.id_val = id_val
    
        return self
    
    def title(self, title: str) -> typing.Self:    
        if not len(title) >= 1:
            raise ValueError("len(title) must be >= 1")
        self._internal.title = title
    
        return self
    
