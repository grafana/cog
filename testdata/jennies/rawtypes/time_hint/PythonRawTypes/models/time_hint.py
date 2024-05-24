import typing


ObjTime: typing.TypeAlias = str


class ObjWithTimeField:
    registered_at: str

    def __init__(self, registered_at: str = ""):
        self.registered_at = registered_at

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "registeredAt": self.registered_at,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "registeredAt" in data:
            args["registered_at"] = data["registeredAt"]        

        return cls(**args)



