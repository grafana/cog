import typing
from ..cog import builder as cogbuilder
from ..models import dataquery_variant_builder


class LokiBuilder(cogbuilder.Builder[dataquery_variant_builder.Loki]):
    _internal: dataquery_variant_builder.Loki

    def __init__(self):
        self._internal = dataquery_variant_builder.Loki()

    def build(self) -> dataquery_variant_builder.Loki:
        """
        Builds the object.
        """
        return self._internal    
    
    def expr(self, expr: str) -> typing.Self:    
        self._internal.expr = expr
    
        return self
    
