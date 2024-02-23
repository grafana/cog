import typing
from ..cog import builder as cogbuilder
from ..models import nullable_map_assignment


class SomeStruct(cogbuilder.Builder[nullable_map_assignment.SomeStruct]):    
    __internal: nullable_map_assignment.SomeStruct

    def __init__(self):
        self.__internal = nullable_map_assignment.SomeStruct()

    def build(self) -> nullable_map_assignment.SomeStruct:
        return self.__internal    
    
    def config(self, config: dict[str, str]) -> typing.Self:        
        self.__internal.config = config
    
        return self
    