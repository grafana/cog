import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):    
    __internal: sandbox.SomeStruct

    def __init__(self):
        self.__internal = sandbox.SomeStruct()

    def build(self) -> sandbox.SomeStruct:
        return self.__internal    
    
    def time(self, from_val: str, to: str) -> typing.Self:        
        if self.__internal.time is None:
            self.__internal.time = "unknown"
        
        assert isinstance(self.__internal.time, unknown)
        
        self.__internal.time.from_val = from_val    
        if self.__internal.time is None:
            self.__internal.time = "unknown"
        
        assert isinstance(self.__internal.time, unknown)
        
        self.__internal.time.to = to
    
        return self
    