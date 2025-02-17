import typing
from ..cog import builder as cogbuilder
from ..models import map_of_builders


class Panel(cogbuilder.Builder[map_of_builders.Panel]):
    _internal: map_of_builders.Panel

    def __init__(self):
        self._internal = map_of_builders.Panel()

    def build(self) -> map_of_builders.Panel:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    


class Dashboard(cogbuilder.Builder[map_of_builders.Dashboard]):
    _internal: map_of_builders.Dashboard

    def __init__(self):
        self._internal = map_of_builders.Dashboard()

    def build(self) -> map_of_builders.Dashboard:
        """
        Builds the object.
        """
        return self._internal    
    
    def panels(self, panels: dict[str, cogbuilder.Builder[map_of_builders.Panel]]) -> typing.Self:    
        panels_resources = { key1: val1.build() for (key1, val1) in panels.items() }
        self._internal.panels = panels_resources
    
        return self
    
