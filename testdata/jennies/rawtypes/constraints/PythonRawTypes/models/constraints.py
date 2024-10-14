import typing


class SomeStruct:
    id_val: int
    maybe_id: typing.Optional[int]
    title: str
    ref_struct: typing.Optional['RefStruct']

    def __init__(self, id_val: int = 0, maybe_id: typing.Optional[int] = None, title: str = "", ref_struct: typing.Optional['RefStruct'] = None):
        self.id_val = id_val
        self.maybe_id = maybe_id
        self.title = title
        self.ref_struct = ref_struct

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "id": self.id_val,
            "title": self.title,
        }
        if self.maybe_id is not None:
            payload["maybeId"] = self.maybe_id
        if self.ref_struct is not None:
            payload["refStruct"] = self.ref_struct
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "id" in data:
            args["id_val"] = data["id"]
        if "maybeId" in data:
            args["maybe_id"] = data["maybeId"]
        if "title" in data:
            args["title"] = data["title"]
        if "refStruct" in data:
            args["ref_struct"] = RefStruct.from_json(data["refStruct"])        

        return cls(**args)


class RefStruct:
    labels: dict[str, str]
    tags: list[str]

    def __init__(self, labels: typing.Optional[dict[str, str]] = None, tags: typing.Optional[list[str]] = None):
        self.labels = labels if labels is not None else {}
        self.tags = tags if tags is not None else []

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "labels": self.labels,
            "tags": self.tags,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "labels" in data:
            args["labels"] = data["labels"]
        if "tags" in data:
            args["tags"] = data["tags"]        

        return cls(**args)



