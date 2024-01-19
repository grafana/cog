import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):    
    __internal: sandbox.SomeStruct

    def __init__(self):
        self.__internal = sandbox.SomeStruct()

    def build(self) -> sandbox.SomeStruct:
        return self.__internal    
    
    def editable(self) -> typing.Self:        
        self.__internal.editable = True
    
        return self
    
    def readonly(self) -> typing.Self:        
        self.__internal.editable = False
    
        return self
    
    def auto_refresh(self) -> typing.Self:        
        self.__internal.auto_refresh = True
    
        return self
    
    def no_auto_refresh(self) -> typing.Self:        
        self.__internal.auto_refresh = False
    
        return self
    