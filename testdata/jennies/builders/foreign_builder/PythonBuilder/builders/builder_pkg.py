import typing
from ..cog import builder as cogbuilder
from ..models import some_pkg


class SomeNiceBuilder(cogbuilder.Builder[some_pkg.SomeStruct]):
    _internal: some_pkg.SomeStruct

    def __init__(self):
        self._internal = some_pkg.SomeStruct()

    def build(self) -> some_pkg.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
