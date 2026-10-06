import typing
from ..cog import builder as cogbuilder
from ..models import sandbox


class Dashboard(cogbuilder.Builder[sandbox.Dashboard]):
    _internal: sandbox.Dashboard

    def __init__(self) -> None:
        self._internal = sandbox.Dashboard()

    def build(self) -> sandbox.Dashboard:
        """
        Builds the object.
        """
        return self._internal    
    
    def variables(self, variables: list[sandbox.Variable]) -> typing.Self:    
        self._internal.variables = variables
    
        return self
    
    def with_variable(self, name: str, value: str) -> typing.Self:    
        if self._internal.variables is None:
            self._internal.variables = []
        
        self._internal.variables.append(sandbox.Variable(
            name=name,
            value=value,
        ))
    
        return self
    
