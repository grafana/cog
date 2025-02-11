import typing
from ..cog import builder as cogbuilder
from ..models import constructor_initializations


class SomePanel(cogbuilder.Builder[constructor_initializations.SomePanel]):
    _internal: constructor_initializations.SomePanel

    def __init__(self):
        self._internal = constructor_initializations.SomePanel()        
        self._internal.type_val = "panel_type"        
        self._internal.cursor = constructor_initializations.CursorMode.TOOLTIP

    def build(self) -> constructor_initializations.SomePanel:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
