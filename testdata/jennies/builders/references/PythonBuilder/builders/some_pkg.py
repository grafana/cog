from ..cog import builder as cogbuilder
from ..models import other_pkg
from ..models import some_pkg
import typing


class Person(cogbuilder.Builder[some_pkg.Person]):    
    __internal: some_pkg.Person

    def __init__(self):
        self.__internal = some_pkg.Person()

    def build(self) -> some_pkg.Person:
        return self.__internal    
    
    def name(self, name: other_pkg.Name) -> typing.Self:        
        self.__internal.name = name
    
        return self
    