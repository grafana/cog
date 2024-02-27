import typing
from ..cog import builder as cogbuilder
from ..models import builder_delegation_in_disjunction


class DashboardLink(cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink]):    
    __internal: builder_delegation_in_disjunction.DashboardLink

    def __init__(self):
        self.__internal = builder_delegation_in_disjunction.DashboardLink()

    def build(self) -> builder_delegation_in_disjunction.DashboardLink:
        return self.__internal    
    
    def title(self, title: str) -> typing.Self:        
        self.__internal.title = title
    
        return self
    
    def url(self, url: str) -> typing.Self:        
        self.__internal.url = url
    
        return self
    

class ExternalLink(cogbuilder.Builder[builder_delegation_in_disjunction.ExternalLink]):    
    __internal: builder_delegation_in_disjunction.ExternalLink

    def __init__(self):
        self.__internal = builder_delegation_in_disjunction.ExternalLink()

    def build(self) -> builder_delegation_in_disjunction.ExternalLink:
        return self.__internal    
    
    def url(self, url: str) -> typing.Self:        
        self.__internal.url = url
    
        return self
    

class Dashboard(cogbuilder.Builder[builder_delegation_in_disjunction.Dashboard]):    
    __internal: builder_delegation_in_disjunction.Dashboard

    def __init__(self):
        self.__internal = builder_delegation_in_disjunction.Dashboard()

    def build(self) -> builder_delegation_in_disjunction.Dashboard:
        return self.__internal    
    
    def single_link_or_string(self, single_link_or_string: typing.Union[cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink], str]) -> typing.Self:    
        """
        will be expanded to cog.Builder<DashboardLink> | string
        """
            
        single_link_or_string_resource = single_link_or_string.build()
        self.__internal.single_link_or_string = single_link_or_string_resource
    
        return self
    
    def links_or_strings(self, links_or_strings: list[typing.Union[cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink], str]]) -> typing.Self:    
        """
        will be expanded to [](cog.Builder<DashboardLink> | string)
        """
            
        links_or_strings_resources = [r1.build() for r1 in links_or_strings]
        self.__internal.links_or_strings = links_or_strings_resources
    
        return self
    
    def disjunction_of_builders(self, disjunction_of_builders: typing.Union[cogbuilder.Builder[builder_delegation_in_disjunction.DashboardLink], cogbuilder.Builder[builder_delegation_in_disjunction.ExternalLink]]) -> typing.Self:        
        disjunction_of_builders_resource = disjunction_of_builders.build()
        self.__internal.disjunction_of_builders = disjunction_of_builders_resource
    
        return self
    