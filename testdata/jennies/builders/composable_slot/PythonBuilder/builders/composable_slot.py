from ..cog import builder as cogbuilder
from ..cog import variants as cogvariants
from ..models import composable_slot
import typing


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
        targets_resources = [r1.build() for r1 in targets]
        self.__internal.targets = targets_resources
    
        return self
    