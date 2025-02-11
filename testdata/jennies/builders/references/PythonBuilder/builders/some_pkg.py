import typing
from ..cog import builder as cogbuilder
from ..models import some_pkg
from ..models import other_pkg


class Person(cogbuilder.Builder[some_pkg.Person]):
    _internal: some_pkg.Person

    def __init__(self):
        self._internal = some_pkg.Person()

    def build(self) -> some_pkg.Person:
        """
        Builds the object.
        """
        return self._internal    
    
    def name(self, name: other_pkg.Name) -> typing.Self:    
        self._internal.name = name
    
        return self
    
