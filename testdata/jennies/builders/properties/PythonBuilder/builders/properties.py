import typing
from ..cog import builder as cogbuilder
from ..models import properties


class SomeStruct(cogbuilder.Builder[properties.SomeStruct]):
    _internal: properties.SomeStruct
    __some_builder_property: str = ""

    def __init__(self):
        self._internal = properties.SomeStruct()

    def build(self) -> properties.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def id(self, id_val: int) -> typing.Self:    
        self._internal.id_val = id_val
    
        return self
    
