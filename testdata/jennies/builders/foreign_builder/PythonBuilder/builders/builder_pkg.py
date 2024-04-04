from ..cog import builder as cogbuilder
from ..models import some_pkg
import typing


class SomeNiceBuilder(cogbuilder.Builder[some_pkg.SomeStruct]):    
    __internal: some_pkg.SomeStruct

    def __init__(self):
        self.__internal = some_pkg.SomeStruct()

    def build(self) -> some_pkg.SomeStruct:
        return self.__internal    
    
    def title(self, title: str) -> typing.Self:        
        self.__internal.title = title
    
        return self
    