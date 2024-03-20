from ..cog import variants as cogvariants
import typing
from ..cog import runtime as cogruntime


class Query(cogvariants.Dataquery):
    expr: str
    instant: typing.Optional[bool]

    def __init__(self, expr: str = "", instant: typing.Optional[bool] = None):
        self.expr = expr
        self.instant = instant

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "expr": self.expr,
        }
        if self.instant is not None:
            payload["instant"] = self.instant
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "expr" in data:
            args["expr"] = data["expr"]
        if "instant" in data:
            args["instant"] = data["instant"]        

        return cls(**args)


def variant_config():
    return cogruntime.DataqueryConfig(
        identifier="prometheus",
        from_json_hook=Query.from_json,
    )
