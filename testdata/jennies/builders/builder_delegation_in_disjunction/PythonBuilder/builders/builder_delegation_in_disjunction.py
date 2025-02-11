import typing
from ..cog import builder as cogbuilder
from ..models import builder_delegation_in_disjunction


class DashboardLink(cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink]):
    _internal: builder_delegation_in_disjunction.DashboardLink

    def __init__(self):
        self._internal = builder_delegation_in_disjunction.DashboardLink()

    def build(self) -> builder_delegation_in_disjunction.DashboardLink:
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
    


class ExternalLink(cogbuilder.Builder[builder_delegation_in_disjunction.ExternalLink]):
    _internal: builder_delegation_in_disjunction.ExternalLink

    def __init__(self):
        self._internal = builder_delegation_in_disjunction.ExternalLink()

    def build(self) -> builder_delegation_in_disjunction.ExternalLink:
        """
        Builds the object.
        """
        return self._internal    
    
    def url(self, url: str) -> typing.Self:    
        self._internal.url = url
    
        return self
    


class Dashboard(cogbuilder.Builder[builder_delegation_in_disjunction.Dashboard]):
    _internal: builder_delegation_in_disjunction.Dashboard

    def __init__(self):
        self._internal = builder_delegation_in_disjunction.Dashboard()

    def build(self) -> builder_delegation_in_disjunction.Dashboard:
        """
        Builds the object.
        """
        return self._internal    
    
    def single_link_or_string(self, single_link_or_string: typing.Union[cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink], str]) -> typing.Self:    
        """
        will be expanded to cog.Builder<DashboardLink> | string
        """
            
        single_link_or_string_resource = single_link_or_string.build() if hasattr(single_link_or_string, 'build') and callable(single_link_or_string.build) else single_link_or_string
        assert isinstance(single_link_or_string_resource, builder_delegation_in_disjunction.DashboardLink) or isinstance(single_link_or_string_resource, str)
        self._internal.single_link_or_string = single_link_or_string_resource
    
        return self
    
    def links_or_strings(self, links_or_strings: list[typing.Union[cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink], str]]) -> typing.Self:    
        """
        will be expanded to [](cog.Builder<DashboardLink> | string)
        """
            
        links_or_strings_resources = [r1.build() if hasattr(r1, 'build') and callable(r1.build) else r1 for r1 in links_or_strings]
        self._internal.links_or_strings = links_or_strings_resources
    
        return self
    
    def disjunction_of_builders(self, disjunction_of_builders: typing.Union[cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink], cogbuilder.Builder[builder_delegation_in_disjunction.ExternalLink]]) -> typing.Self:    
        disjunction_of_builders_resource = disjunction_of_builders.build()
        self._internal.disjunction_of_builders = disjunction_of_builders_resource
    
        return self
    
