from ..cog import builder as cogbuilder
from ..models import properties
import typing


class SomeStruct(cogbuilder.Builder[properties.SomeStruct]):    
    __internal: properties.SomeStruct
    __some_builder_property: str = ""

    def __init__(self):
        self.__internal = properties.SomeStruct()

    def build(self) -> properties.SomeStruct:
        return self.__internal    
    
    def id_val(self, id_val: int) -> typing.Self:        
        self.__internal.id_val = id_val
    
        return self
    