import typing
from ..cog import builder as cogbuilder
from ..models import initialization_safeguards


class SomePanel(cogbuilder.Builder[initialization_safeguards.SomePanel]):    
    __internal: initialization_safeguards.SomePanel

    def __init__(self):
        self.__internal = initialization_safeguards.SomePanel()

    def build(self) -> initialization_safeguards.SomePanel:
        return self.__internal    
    
    def title(self, title: str) -> typing.Self:        
        self.__internal.title = title
    
        return self
    
    def show_legend(self, show: boolean) -> typing.Self:        
        if self.__internal.options is None:
            self.__internal.options = initialization_safeguards.Options()
        
        assert isinstance(self.__internal.options, initialization_safeguards.Options)
        
        if self.__internal.options.legend is None:
            self.__internal.options.legend = initialization_safeguards.LegendOptions()
        
        assert isinstance(self.__internal.options.legend, initialization_safeguards.LegendOptions)
        
        self.__internal.options.legend.show = show
    
        return self
    