from ..cog import builder as cogbuilder
from ..models import sandbox
import typing


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):    
    __internal: sandbox.SomeStruct

    def __init__(self, title: str):
        self.__internal = sandbox.SomeStruct()        
        self.__internal.title = title

    def build(self) -> sandbox.SomeStruct:
        return self.__internal    
    
    def title(self, title: str) -> typing.Self:        
        self.__internal.title = title
    
        return self
    