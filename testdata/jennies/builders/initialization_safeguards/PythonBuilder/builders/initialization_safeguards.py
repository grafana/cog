import typing
from ..cog import builder as cogbuilder
from ..models import initialization_safeguards


class SomePanel(cogbuilder.Builder[initialization_safeguards.SomePanel]):
    _internal: initialization_safeguards.SomePanel

    def __init__(self):
        self._internal = initialization_safeguards.SomePanel()

    def build(self) -> initialization_safeguards.SomePanel:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
    def show_legend(self, show: bool) -> typing.Self:    
        if self._internal.options is None:
            self._internal.options = initialization_safeguards.Options()
        assert isinstance(self._internal.options, initialization_safeguards.Options)
        if self._internal.options.legend is None:
            self._internal.options.legend = initialization_safeguards.LegendOptions()
        assert isinstance(self._internal.options.legend, initialization_safeguards.LegendOptions)
        self._internal.options.legend.show = show
    
        return self
    
