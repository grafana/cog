import typing


class StructComplexFieldsSomeStructFieldAnonymousStruct:
    field_any: object

    def __init__(self, field_any: object = None):
        self.field_any = field_any

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldAny": self.field_any,
        }
        return payload


class SomeStruct:
    """
    This struct does things.
    """

    field_ref: 'SomeOtherStruct'
    field_disjunction_of_scalars: typing.Union[str, bool]
    field_mixed_disjunction: typing.Union[str, 'SomeOtherStruct']
    field_disjunction_with_null: typing.Optional[str]
    operator: typing.Literal[">", "<"]
    field_array_of_strings: list[str]
    field_map_of_string_to_string: dict[str, str]
    field_anonymous_struct: 'StructComplexFieldsSomeStructFieldAnonymousStruct'

    def __init__(self, field_ref: typing.Optional['SomeOtherStruct'] = None, field_disjunction_of_scalars: typing.Union[str, bool] = "", field_mixed_disjunction: typing.Union[str, 'SomeOtherStruct'] = "", field_disjunction_with_null: typing.Optional[str] = None, operator: typing.Optional[typing.Literal[">", "<"]] = None, field_array_of_strings: typing.Optional[list[str]] = None, field_map_of_string_to_string: typing.Optional[dict[str, str]] = None, field_anonymous_struct: typing.Optional['StructComplexFieldsSomeStructFieldAnonymousStruct'] = None):
        self.field_ref = field_ref if field_ref is not None else SomeOtherStruct()
        self.field_disjunction_of_scalars = field_disjunction_of_scalars
        self.field_mixed_disjunction = field_mixed_disjunction
        self.field_disjunction_with_null = field_disjunction_with_null
        self.operator = operator if operator is not None else ">"
        self.field_array_of_strings = field_array_of_strings if field_array_of_strings is not None else []
        self.field_map_of_string_to_string = field_map_of_string_to_string if field_map_of_string_to_string is not None else {}
        self.field_anonymous_struct = field_anonymous_struct if field_anonymous_struct is not None else StructComplexFieldsSomeStructFieldAnonymousStruct()

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldRef": None if self.field_ref is None else self.field_ref.to_json(),
            "FieldDisjunctionOfScalars": self.field_disjunction_of_scalars,
            "FieldMixedDisjunction": self.field_mixed_disjunction,
            "FieldDisjunctionWithNull": self.field_disjunction_with_null,
            "Operator": self.operator,
            "FieldArrayOfStrings": self.field_array_of_strings,
            "FieldMapOfStringToString": self.field_map_of_string_to_string,
            "FieldAnonymousStruct": None if self.field_anonymous_struct is None else self.field_anonymous_struct.to_json(),
        }
        return payload


class SomeOtherStruct:
    field_any: object

    def __init__(self, field_any: object = None):
        self.field_any = field_any

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "FieldAny": self.field_any,
        }
        return payload
