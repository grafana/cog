import typing
from ..cog import builder as cogbuilder
from ..models import builder_delegation


class DashboardLink(cogbuilder.Builder[builder_delegation.DashboardLink]):
    _internal: builder_delegation.DashboardLink

    def __init__(self):
        self._internal = builder_delegation.DashboardLink()

    def build(self) -> builder_delegation.DashboardLink:
        """
        Builds the object.
        """
        return self._internal    
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
    def url(self, url: str) -> typing.Self:    
        self._internal.url = url
    
        return self
    


class Dashboard(cogbuilder.Builder[builder_delegation.Dashboard]):
    _internal: builder_delegation.Dashboard

    def __init__(self):
        self._internal = builder_delegation.Dashboard()

    def build(self) -> builder_delegation.Dashboard:
        """
        Builds the object.
        """
        return self._internal    
    
    def id(self, id_val: int) -> typing.Self:    
        self._internal.id_val = id_val
    
        return self
    
    def title(self, title: str) -> typing.Self:    
        self._internal.title = title
    
        return self
    
    def links(self, links: list[cogbuilder.Builder[builder_delegation.DashboardLink]]) -> typing.Self:    
        """
        will be expanded to []cog.Builder<DashboardLink>
        """
            
        links_resources = [r1.build() for r1 in links]
        self._internal.links = links_resources
    
        return self
    
    def links_of_links(self, links_of_links: list[list[cogbuilder.Builder[builder_delegation.DashboardLink]]]) -> typing.Self:    
        """
        will be expanded to [][]cog.Builder<DashboardLink>
        """
            
        links_of_links_resources = [[r2.build() for r2 in r1] for r1 in links_of_links]
        self._internal.links_of_links = links_of_links_resources
    
        return self
    
    def single_link(self, single_link: cogbuilder.Builder[builder_delegation.DashboardLink]) -> typing.Self:    
        """
        will be expanded to cog.Builder<DashboardLink>
        """
            
        single_link_resource = single_link.build()
        self._internal.single_link = single_link_resource
    
        return self
    
