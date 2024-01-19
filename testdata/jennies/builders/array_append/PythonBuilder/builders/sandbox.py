import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):    
    __internal: sandbox.SomeStruct

    def __init__(self):
        self.__internal = sandbox.SomeStruct()

    def build(self) -> sandbox.SomeStruct:
        return self.__internal    
    
    def tags(self, tags: str) -> typing.Self:        
        if self.__internal.tags is None:
            self.__internal.tags = []
        
        self.__internal.tags.append(tags)
    
        return self
    