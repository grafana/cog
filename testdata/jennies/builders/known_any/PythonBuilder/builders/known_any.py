import typing
from ..cog import builder as cogbuilder
from ..models import known_any


class SomeStruct(cogbuilder.Builder[known_any.SomeStruct]):
    _internal: known_any.SomeStruct

    def __init__(self):
        self._internal = known_any.SomeStruct()

    def build(self) -> known_any.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        if self._internal.config is None:
            self._internal.config = known_any.Config()
        assert isinstance(self._internal.config, known_any.Config)
        self._internal.config.title = title
    
        return self
    
