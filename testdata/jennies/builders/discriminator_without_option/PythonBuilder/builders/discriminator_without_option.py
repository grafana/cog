import typing
from ..cog import builder as cogbuilder
from ..models import discriminator_without_option


class NoShowFieldOption(cogbuilder.Builder[discriminator_without_option.NoShowFieldOption]):
    _internal: discriminator_without_option.NoShowFieldOption

    def __init__(self):
        self._internal = discriminator_without_option.NoShowFieldOption()

    def build(self) -> discriminator_without_option.NoShowFieldOption:
        """
        Builds the object.
        """
        return self._internal    
    
    def text(self, text: str) -> typing.Self:    
        self._internal.text = text
    
        return self
    


class ShowFieldOption(cogbuilder.Builder[discriminator_without_option.ShowFieldOption]):
    _internal: discriminator_without_option.ShowFieldOption

    def __init__(self):
        self._internal = discriminator_without_option.ShowFieldOption()

    def build(self) -> discriminator_without_option.ShowFieldOption:
        """
        Builds the object.
        """
        return self._internal    
    
    def field(self, field: discriminator_without_option.AnEnum) -> typing.Self:    
        self._internal.field = field
    
        return self
    
    def text(self, text: str) -> typing.Self:    
        self._internal.text = text
    
        return self
