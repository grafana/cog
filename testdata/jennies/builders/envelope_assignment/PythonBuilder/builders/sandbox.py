import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class Dashboard(cogbuilder.Builder[sandbox.Dashboard]):    
    __internal: sandbox.Dashboard

    def __init__(self):
        self.__internal = sandbox.Dashboard()

    def build(self) -> sandbox.Dashboard:
        return self.__internal    
    
    def with_variable(self, name: str, value: str) -> typing.Self:        
        if self.__internal.variables is None:
            self.__internal.variables = []
        
        self.__internal.variables.append(sandbox.Variable(
            name=name,
            value=value,
        ))
    
        return self
    