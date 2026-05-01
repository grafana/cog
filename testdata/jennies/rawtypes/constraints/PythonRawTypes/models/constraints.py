import typing


class SomeStruct:
    id_val: int
    maybe_id: typing.Optional[int]
    greater_than_zero: int
    negative: int
    title: str
    labels: dict[str, str]
    tags: list[str]
    regex: str
    negative_regex: str
    min_max_list: list[str]
    unique_list: list[str]
    full_constraint_list: list[int]

    def __init__(self, id_val: int = 0, maybe_id: typing.Optional[int] = None, greater_than_zero: int = 0, negative: int = 0, title: str = "", labels: typing.Optional[dict[str, str]] = None, tags: typing.Optional[list[str]] = None, regex: str = "", negative_regex: str = "", min_max_list: typing.Optional[list[str]] = None, unique_list: typing.Optional[list[str]] = None, full_constraint_list: typing.Optional[list[int]] = None) -> None:
        self.id_val = id_val
        self.maybe_id = maybe_id
        self.greater_than_zero = greater_than_zero
        self.negative = negative
        self.title = title
        self.labels = labels if labels is not None else {}
        self.tags = tags if tags is not None else []
        self.regex = regex
        self.negative_regex = negative_regex
        self.min_max_list = min_max_list if min_max_list is not None else []
        self.unique_list = unique_list if unique_list is not None else []
        self.full_constraint_list = full_constraint_list if full_constraint_list is not None else []

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "id": self.id_val,
            "greaterThanZero": self.greater_than_zero,
            "negative": self.negative,
            "title": self.title,
            "labels": self.labels,
            "tags": self.tags,
            "regex": self.regex,
            "negativeRegex": self.negative_regex,
            "minMaxList": self.min_max_list,
            "uniqueList": self.unique_list,
            "fullConstraintList": self.full_constraint_list,
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
        if "regex" in data:
            args["regex"] = data["regex"]
        if "negativeRegex" in data:
            args["negative_regex"] = data["negativeRegex"]
        if "minMaxList" in data:
            args["min_max_list"] = data["minMaxList"]
        if "uniqueList" in data:
            args["unique_list"] = data["uniqueList"]
        if "fullConstraintList" in data:
            args["full_constraint_list"] = data["fullConstraintList"]        

        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, SomeStruct):
            return False
        if self.id_val != other.id_val:
            return False
        if self.maybe_id != other.maybe_id:
            return False
        if self.greater_than_zero != other.greater_than_zero:
            return False
        if self.negative != other.negative:
            return False
        if self.title != other.title:
            return False
        if self.labels != other.labels:
            return False
        if self.tags != other.tags:
            return False
        if self.regex != other.regex:
            return False
        if self.negative_regex != other.negative_regex:
            return False
        if self.min_max_list != other.min_max_list:
            return False
        if self.unique_list != other.unique_list:
            return False
        if self.full_constraint_list != other.full_constraint_list:
            return False
        return True
