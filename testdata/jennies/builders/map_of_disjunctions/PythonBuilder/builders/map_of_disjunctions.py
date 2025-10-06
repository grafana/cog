import typing
from ..cog import builder as cogbuilder
from ..models import map_of_disjunctions


class Element(cogbuilder.Builder[map_of_disjunctions.Element]):
    _internal: map_of_disjunctions.Element

    def __init__(self):
        self._internal = map_of_disjunctions.Element()

    def build(self) -> map_of_disjunctions.Element:
        """
        Builds the object.
        """
        return self._internal    
    
    def panel(self, panel: cogbuilder.Builder[map_of_disjunctions.Panel]) -> typing.Self:    
        panel_resource = panel.build()
        self._internal.panel = panel_resource
    
        return self
    
    def library_panel(self, library_panel: cogbuilder.Builder[map_of_disjunctions.LibraryPanel]) -> typing.Self:    
        library_panel_resource = library_panel.build()
        self._internal.library_panel = library_panel_resource
    
        return self
    


class Panel(cogbuilder.Builder[map_of_disjunctions.Panel]):
    _internal: map_of_disjunctions.Panel

    def __init__(self):
        self._internal = map_of_disjunctions.Panel()        
        self._internal.kind = "Panel"

    def build(self) -> map_of_disjunctions.Panel:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    


class LibraryPanel(cogbuilder.Builder[map_of_disjunctions.LibraryPanel]):
    _internal: map_of_disjunctions.LibraryPanel

    def __init__(self):
        self._internal = map_of_disjunctions.LibraryPanel()        
        self._internal.kind = "Library"

    def build(self) -> map_of_disjunctions.LibraryPanel:
        """
        Builds the object.
        """
        return self._internal    
    
    def text(self, text: str) -> typing.Self:    
        self._internal.text = text
    
        return self
    


class Dashboard(cogbuilder.Builder[map_of_disjunctions.Dashboard]):
    _internal: map_of_disjunctions.Dashboard

    def __init__(self):
        self._internal = map_of_disjunctions.Dashboard()

    def build(self) -> map_of_disjunctions.Dashboard:
        """
        Builds the object.
        """
        return self._internal    
    
    def panels(self, panels: dict[str, cogbuilder.Builder[map_of_disjunctions.Element]]) -> typing.Self:    
        panels_resources = { key1: val1.build() for (key1, val1) in panels.items() }
        self._internal.panels = panels_resources
    
        return self
    


class PanelOrLibraryPanel(cogbuilder.Builder[map_of_disjunctions.PanelOrLibraryPanel]):
    _internal: map_of_disjunctions.PanelOrLibraryPanel

    def __init__(self):
        self._internal = map_of_disjunctions.PanelOrLibraryPanel()

    def build(self) -> map_of_disjunctions.PanelOrLibraryPanel:
        """
        Builds the object.
        """
        return self._internal    
    
    def panel(self, panel: cogbuilder.Builder[map_of_disjunctions.Panel]) -> typing.Self:    
        panel_resource = panel.build()
        self._internal.panel = panel_resource
    
        return self
    
    def library_panel(self, library_panel: cogbuilder.Builder[map_of_disjunctions.LibraryPanel]) -> typing.Self:    
        library_panel_resource = library_panel.build()
        self._internal.library_panel = library_panel_resource
    
        return self
    
