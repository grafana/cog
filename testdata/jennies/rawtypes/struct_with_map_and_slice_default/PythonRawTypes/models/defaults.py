import typing


class SomeStruct:
    options: typing.Optional[dict[str, object]]
    items: typing.Optional[list[str]]
    extra: object

    def __init__(self, options: typing.Optional[dict[str, object]] = None, items: typing.Optional[list[str]] = None, extra: object = {}) -> None:
        self.options = options if options is not None else {}
        self.items = items if items is not None else []
        self.extra = extra

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "extra": self.extra,
        }
        if self.options is not None:
            payload["options"] = self.options
        if self.items is not None:
            payload["items"] = self.items
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "options" in data:
            args["options"] = data["options"]
        if "items" in data:
            args["items"] = data["items"]
        if "extra" in data:
            args["extra"] = data["extra"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, SomeStruct):
            return False
        if self.options != other.options:
            return False
        if self.items != other.items:
            return False
        if self.extra != other.extra:
            return False
        return True
