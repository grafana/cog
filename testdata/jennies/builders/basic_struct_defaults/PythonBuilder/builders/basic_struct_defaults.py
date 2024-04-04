from ..models import basic_struct_defaults
from ..cog import builder as cogbuilder
import typing


class SomeStruct(cogbuilder.Builder[basic_struct_defaults.SomeStruct]):    
    __internal: basic_struct_defaults.SomeStruct

    def __init__(self):
        self.__internal = basic_struct_defaults.SomeStruct()

    def build(self) -> basic_struct_defaults.SomeStruct:
        return self.__internal    
    
    def id_val(self, id_val: int) -> typing.Self:        
        self.__internal.id_val = id_val
    
        return self
    
    def uid(self, uid: str) -> typing.Self:        
        self.__internal.uid = uid
    
        return self
    
    def tags(self, tags: list[str]) -> typing.Self:        
        self.__internal.tags = tags
    
        return self
    
    def live_now(self, live_now: bool) -> typing.Self:        
        self.__internal.live_now = live_now
    
        return self
    