import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class SomeStruct(cogbuilder.Builder[sandbox.SomeStruct]):
    _internal: sandbox.SomeStruct

    def __init__(self):
        self._internal = sandbox.SomeStruct()

    def build(self) -> sandbox.SomeStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def time(self, from_val: str, to: str) -> typing.Self:    
        if self._internal.time is None:
            self._internal.time = "unknown"
        assert isinstance(self._internal.time, unknown)
        self._internal.time.from_val = from_val    
        self._internal.time.to = to
    
        return self
    
