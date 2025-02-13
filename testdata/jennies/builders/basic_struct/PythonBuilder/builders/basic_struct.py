import typing
from ..cog import builder as cogbuilder
from ..models import basic_struct


class SomeStruct(cogbuilder.Builder[basic_struct.SomeStruct]):    
    """
    SomeStruct, to hold data.
    """
    
    _internal: basic_struct.SomeStruct

    def __init__(self):
        self._internal = basic_struct.SomeStruct()

    def build(self) -> basic_struct.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def id(self, id_val: int) -> typing.Self:    
        """
        id identifies something. Weird, right?
        """
            
        self._internal.id_val = id_val
    
        return self
    
    def uid(self, uid: str) -> typing.Self:    
        self._internal.uid = uid
    
        return self
    
    def tags(self, tags: list[str]) -> typing.Self:    
        self._internal.tags = tags
    
        return self
    
    def live_now(self, live_now: bool) -> typing.Self:    
        """
        This thing could be live.
        Or maybe not.
        """
            
        self._internal.live_now = live_now
    
        return self
    
