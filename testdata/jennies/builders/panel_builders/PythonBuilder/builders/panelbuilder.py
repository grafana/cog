import typing
from ..cog import builder as cogbuilder
from ..models import panelbuilder


class Panel(cogbuilder.Builder[panelbuilder.Panel]):
    _internal: panelbuilder.Panel

    def __init__(self):
        self._internal = panelbuilder.Panel()

    def build(self) -> panelbuilder.Panel:
        """
        Builds the object.
        """
        return self._internal    
    
    def only_from_this_dashboard(self, only_from_this_dashboard: bool) -> typing.Self:    
        self._internal.only_from_this_dashboard = only_from_this_dashboard
    
        return self
    
    def only_in_time_range(self, only_in_time_range: bool) -> typing.Self:    
        self._internal.only_in_time_range = only_in_time_range
    
        return self
    
    def tags(self, tags: list[str]) -> typing.Self:    
        self._internal.tags = tags
    
        return self
    
    def limit(self, limit: int) -> typing.Self:    
        self._internal.limit = limit
    
        return self
    
    def show_user(self, show_user: bool) -> typing.Self:    
        self._internal.show_user = show_user
    
        return self
    
    def show_time(self, show_time: bool) -> typing.Self:    
        self._internal.show_time = show_time
    
        return self
    
    def show_tags(self, show_tags: bool) -> typing.Self:    
        self._internal.show_tags = show_tags
    
        return self
    
    def navigate_to_panel(self, navigate_to_panel: bool) -> typing.Self:    
        self._internal.navigate_to_panel = navigate_to_panel
    
        return self
    
    def navigate_before(self, navigate_before: str) -> typing.Self:    
        self._internal.navigate_before = navigate_before
    
        return self
    
    def navigate_after(self, navigate_after: str) -> typing.Self:    
        self._internal.navigate_after = navigate_after
    
        return self
    
