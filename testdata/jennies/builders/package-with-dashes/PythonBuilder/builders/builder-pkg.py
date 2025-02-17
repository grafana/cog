import typing
from ..cog import builder as cogbuilder
from ..models import with-dashes


class SomeNiceBuilder(cogbuilder.Builder[with-dashes.SomeStruct]):
    _internal: with-dashes.SomeStruct

    def __init__(self):
        self._internal = with-dashes.SomeStruct()

    def build(self) -> with-dashes.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
