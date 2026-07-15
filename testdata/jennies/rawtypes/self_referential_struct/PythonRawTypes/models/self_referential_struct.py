import typing


class Node:
    """
    Node represents a node in a singly-linked list.
    The next field points to the following node, or is absent if this is the last node.
    """

    value: str
    next_val: typing.Optional['Node']

    def __init__(self, value: str = "", next_val: typing.Optional['Node'] = None) -> None:
        self.value = value
        self.next_val = next_val

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "value": self.value,
        }
        if self.next_val is not None:
            payload["next"] = self.next_val
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "value" in data:
            args["value"] = data["value"]
        if "next" in data:
            args["next_val"] = Node.from_json(data["next"])        

        return cls(**args)
