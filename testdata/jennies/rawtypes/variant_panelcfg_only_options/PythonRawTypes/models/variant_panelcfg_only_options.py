import typing


class Options:
    content: str

    def __init__(self, content: str = "") -> None:
        self.content = content

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "content": self.content,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "content" in data:
            args["content"] = data["content"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, Options):
            return False
        if self.content != other.content:
            return False
        return True
