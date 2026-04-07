import typing
from ..cog import builder as cogbuilder
from ..models import with_dashes


class SomeNiceBuilder(cogbuilder.Builder[with_dashes.SomeStruct]):
    _internal: with_dashes.SomeStruct

    def __init__(self) -> None:
        self._internal = with_dashes.SomeStruct()

    def build(self) -> with_dashes.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
