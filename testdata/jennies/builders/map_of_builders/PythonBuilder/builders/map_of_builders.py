import typing
from ..cog import builder as cogbuilder
from ..models import map_of_builders


class Panel(cogbuilder.Builder[map_of_builders.Panel]):
    _internal: map_of_builders.Panel

    def __init__(self) -> None:
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

    def __init__(self) -> None:
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
    
    def facets(self, facets: dict[str, list[cogbuilder.Builder[map_of_builders.FacetValue]]]) -> typing.Self:    
        facets_resources = { key1: [r2.build() for r2 in val1] for (key1, val1) in facets.items() }
        self._internal.facets = facets_resources
    
        return self
    


class FacetValue(cogbuilder.Builder[map_of_builders.FacetValue]):
    _internal: map_of_builders.FacetValue

    def __init__(self) -> None:
        self._internal = map_of_builders.FacetValue()

    def build(self) -> map_of_builders.FacetValue:
        """
        Builds the object.
        """
        return self._internal    
    
    def value(self, value: str) -> typing.Self:    
        self._internal.value = value
    
        return self
    
    def count(self, count: int) -> typing.Self:    
        self._internal.count = count
    
        return self
    
