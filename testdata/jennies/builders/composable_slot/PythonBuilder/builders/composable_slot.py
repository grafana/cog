import typing
from ..cog import builder as cogbuilder
from ..models import composable_slot
from ..cog import variants as cogvariants


class LokiBuilder(cogbuilder.Builder[composable_slot.Dashboard]):    
    __internal: composable_slot.Dashboard

    def __init__(self):
        self.__internal = composable_slot.Dashboard()

    def build(self) -> composable_slot.Dashboard:
        return self.__internal    
    
    def target(self, target: cogbuilder.Builder[cogvariants.Dataquery]) -> typing.Self:        
        target_resource = target.build()
        self.__internal.target = target_resource
    
        return self
    
    def targets(self, targets: list[cogbuilder.Builder[cogvariants.Dataquery]]) -> typing.Self:        
        targets_resources = [r.build() for r in targets]
        self.__internal.targets = targets_resources
    
        return self
    