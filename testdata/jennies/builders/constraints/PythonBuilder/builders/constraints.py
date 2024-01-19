import typing
from ..cog import builder as cogbuilder
from ..models import constraints


class SomeStruct(cogbuilder.Builder[constraints.SomeStruct]):    
    __internal: constraints.SomeStruct

    def __init__(self):
        self.__internal = constraints.SomeStruct()

    def build(self) -> constraints.SomeStruct:
        return self.__internal    
    
    def id_val(self, id_val: int) -> typing.Self:        
        if not id_val >= 5:
            raise ValueError("id_val must be >= 5")
        if not id_val < 10:
            raise ValueError("id_val must be < 10")
        self.__internal.id_val = id_val
    
        return self
    
    def title(self, title: str) -> typing.Self:        
        if not len(title) >= 1:
            raise ValueError("len(title) must be >= 1")
        self.__internal.title = title
    
        return self
    