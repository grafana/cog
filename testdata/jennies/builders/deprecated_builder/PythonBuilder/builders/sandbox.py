import typing
from ..cog import builder as cogbuilder
from ..models import sandbox
import warnings


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):
    warnings.warn("This builder is deprecated. Don't use. Please.", DeprecationWarning)

    _internal: sandbox.SomeStruct

    def __init__(self) -> None:
        self._internal = sandbox.SomeStruct()

    def build(self) -> sandbox.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
