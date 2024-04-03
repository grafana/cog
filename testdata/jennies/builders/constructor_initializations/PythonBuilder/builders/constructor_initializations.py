from ..cog import builder as cogbuilder
from ..models import constructor_initializations
import typing


class SomePanel(cogbuilder.Builder[constructor_initializations.SomePanel]):    
    __internal: constructor_initializations.SomePanel

    def __init__(self):
        self.__internal = constructor_initializations.SomePanel()        
        self.__internal.type_val = "panel_type"        
        self.__internal.cursor = constructor_initializations.CursorMode.TOOLTIP

    def build(self) -> constructor_initializations.SomePanel:
        return self.__internal    
    
    def title(self, title: str) -> typing.Self:        
        self.__internal.title = title
    
        return self
    