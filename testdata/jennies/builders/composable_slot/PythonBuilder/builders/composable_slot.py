import typing
from ..cog import builder as cogbuilder
from ..models import composable_slot
from ..cog import variants as cogvariants


class LokiBuilder(cogbuilder.Builder[composable_slot.Dashboard]):
    _internal: composable_slot.Dashboard

    def __init__(self):
        self._internal = composable_slot.Dashboard()

    def build(self) -> composable_slot.Dashboard:
        """
        Builds the object.
        """
        return self._internal    
    
    def target(self, target: cogbuilder.Builder[cogvariants.Dataquery]) -> typing.Self:    
        target_resource = target.build()
        self._internal.target = target_resource
    
        return self
    
    def targets(self, targets: list[cogbuilder.Builder[cogvariants.Dataquery]]) -> typing.Self:    
        targets_resources = [r1.build() for r1 in targets]
        self._internal.targets = targets_resources
    
        return self
    
