import typing


ConstantRefString: typing.Literal["AString"] = "AString"


class MyStruct:
    a_string: str
    opt_string: str

    def __init__(self, ) -> None:
        self.a_string = ConstantRefString
        self.opt_string = ConstantRefString

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "aString": self.a_string,
        }
        if self.opt_string is not None:
            payload["optString"] = self.opt_string
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        return cls(**args)

    def __eq__(self, other: object) -> bool:
        if not isinstance(other, MyStruct):
            return False
        if self.a_string != other.a_string:
            return False
        if self.opt_string != other.opt_string:
            return False
        return True



