import typing
from ..cog import builder as cogbuilder
from ..models import with-dashes


class SomeNiceBuilder(cogbuilder.Builder[with-dashes.SomeStruct]):    
    __internal: with-dashes.SomeStruct

    def __init__(self):
        self.__internal = with-dashes.SomeStruct()

    def build(self) -> with-dashes.SomeStruct:
        return self.__internal    
    
    def title(self, title: str) -> typing.Self:        
        self.__internal.title = title
    
        return self
    