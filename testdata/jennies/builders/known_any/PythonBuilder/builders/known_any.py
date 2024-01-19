import typing
from ..cog import builder as cogbuilder
from ..models import known_any


class SomeStruct(cogbuilder.Builder[known_any.SomeStruct]):    
    __internal: known_any.SomeStruct

    def __init__(self):
        self.__internal = known_any.SomeStruct()

    def build(self) -> known_any.SomeStruct:
        return self.__internal    
    
    def title(self, title: str) -> typing.Self:        
        if self.__internal.config is None:
            self.__internal.config = known_any.Config()
        
        assert isinstance(self.__internal.config, known_any.Config)
        
        self.__internal.config.title = title
    
        return self
    