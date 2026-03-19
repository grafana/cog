import typing


class SomeStruct:
    id_val: int
    maybe_id: typing.Optional[int]
    greater_than_zero: int
    negative: int
    title: str
    labels: dict[str, str]
    tags: list[str]

    def __init__(self, id_val: int = 0, maybe_id: typing.Optional[int] = None, greater_than_zero: int = 0, negative: int = 0, title: str = "", labels: typing.Optional[dict[str, str]] = None, tags: typing.Optional[list[str]] = None) -> None:
        self.id_val = id_val
        self.maybe_id = maybe_id
        self.greater_than_zero = greater_than_zero
        self.negative = negative
        self.title = title
        self.labels = labels if labels is not None else {}
        self.tags = tags if tags is not None else []

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "id": self.id_val,
            "greaterThanZero": self.greater_than_zero,
            "negative": self.negative,
            "title": self.title,
            "labels": self.labels,
            "tags": self.tags,
        }
        if self.maybe_id is not None:
            payload["maybeId"] = self.maybe_id
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "id" in data:
            args["id_val"] = data["id"]
        if "maybeId" in data:
            args["maybe_id"] = data["maybeId"]
        if "greaterThanZero" in data:
            args["greater_than_zero"] = data["greaterThanZero"]
        if "negative" in data:
            args["negative"] = data["negative"]
        if "title" in data:
            args["title"] = data["title"]
        if "labels" in data:
            args["labels"] = data["labels"]
        if "tags" in data:
            args["tags"] = data["tags"]        

        return cls(**args)
