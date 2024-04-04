from ..models import basic_struct
from ..cog import builder as cogbuilder
import typing


class SomeStruct(cogbuilder.Builder[basic_struct.SomeStruct]):    
    """
    SomeStruct, to hold data.
    """
    
    __internal: basic_struct.SomeStruct

    def __init__(self):
        self.__internal = basic_struct.SomeStruct()

    def build(self) -> basic_struct.SomeStruct:
        return self.__internal    
    
    def id_val(self, id_val: int) -> typing.Self:    
        """
        id identifies something. Weird, right?
        """
            
        self.__internal.id_val = id_val
    
        return self
    
    def uid(self, uid: str) -> typing.Self:        
        self.__internal.uid = uid
    
        return self
    
    def tags(self, tags: list[str]) -> typing.Self:        
        self.__internal.tags = tags
    
        return self
    
    def live_now(self, live_now: bool) -> typing.Self:    
        """
        This thing could be live.
        Or maybe not.
        """
            
        self.__internal.live_now = live_now
    
        return self
    