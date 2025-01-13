import typing


class Options:
    content: str

    def __init__(self, content: str = ""):
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
