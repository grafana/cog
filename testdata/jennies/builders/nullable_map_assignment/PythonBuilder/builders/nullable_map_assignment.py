import typing
from ..cog import builder as cogbuilder
from ..models import nullable_map_assignment


class SomeStruct(cogbuilder.Builder[nullable_map_assignment.SomeStruct]):
    _internal: nullable_map_assignment.SomeStruct

    def __init__(self):
        self._internal = nullable_map_assignment.SomeStruct()

    def build(self) -> nullable_map_assignment.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def config(self, config: dict[str, str]) -> typing.Self:    
        self._internal.config = config
    
        return self
    
