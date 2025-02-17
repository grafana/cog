import typing
from ..cog import builder as cogbuilder
from ..models import struct_with_defaults


class NestedStruct(cogbuilder.Builder[struct_with_defaults.NestedStruct]):
    _internal: struct_with_defaults.NestedStruct

    def __init__(self):
        self._internal = struct_with_defaults.NestedStruct()

    def build(self) -> struct_with_defaults.NestedStruct:
        """
        Builds the object.
        """
        return self._internal    
    
    def string_val(self, string_val: str) -> typing.Self:    
        self._internal.string_val = string_val
    
        return self
    
    def int_val(self, int_val: int) -> typing.Self:    
        self._internal.int_val = int_val
    
        return self
    


class Struct(cogbuilder.Builder[struct_with_defaults.Struct]):
    _internal: struct_with_defaults.Struct

    def __init__(self):
        self._internal = struct_with_defaults.Struct()

    def build(self) -> struct_with_defaults.Struct:
        """
        Builds the object.
        """
        return self._internal    
    
    def all_fields(self, all_fields: cogbuilder.Builder[struct_with_defaults.NestedStruct]) -> typing.Self:    
        all_fields_resource = all_fields.build()
        self._internal.all_fields = all_fields_resource
    
        return self
    
    def partial_fields(self, partial_fields: cogbuilder.Builder[struct_with_defaults.NestedStruct]) -> typing.Self:    
        partial_fields_resource = partial_fields.build()
        self._internal.partial_fields = partial_fields_resource
    
        return self
    
    def empty_fields(self, empty_fields: cogbuilder.Builder[struct_with_defaults.NestedStruct]) -> typing.Self:    
        empty_fields_resource = empty_fields.build()
        self._internal.empty_fields = empty_fields_resource
    
        return self
    
    def complex_field(self, complex_field: unknown) -> typing.Self:    
        self._internal.complex_field = complex_field
    
        return self
    
    def partial_complex_field(self, partial_complex_field: unknown) -> typing.Self:    
        self._internal.partial_complex_field = partial_complex_field
    
        return self
    
