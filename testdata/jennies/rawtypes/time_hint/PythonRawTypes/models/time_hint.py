import typing


ObjTime: typing.TypeAlias = str


class ObjWithTimeField:
    registered_at: str
    duration: str

    def __init__(self, registered_at: str = "", duration: str = "") -> None:
        self.registered_at = registered_at
        self.duration = duration

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "registeredAt": self.registered_at,
            "duration": self.duration,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "registeredAt" in data:
            args["registered_at"] = data["registeredAt"]
        if "duration" in data:
            args["duration"] = data["duration"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, ObjWithTimeField):
            return False
        if self.registered_at != other.registered_at:
            return False
        if self.duration != other.duration:
            return False
        return True



