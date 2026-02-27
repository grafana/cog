package struct_optional_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type SomeStruct struct {
 FieldRef SomeOtherStruct `tfsdk:"FieldRef"`
FieldString types.String `tfsdk:"FieldString"`
Operator types.String `tfsdk:"Operator"`
FieldArrayOfStrings types.List `tfsdk:"FieldArrayOfStrings"`
FieldAnonymousStruct StructOptionalFieldsSomeStructFieldAnonymousStruct `tfsdk:"FieldAnonymousStruct"`
 }

type SomeOtherStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type StructOptionalFieldsSomeStructFieldAnonymousStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }



var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"field_ref": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
"field_string": types.StringType,
"operator": types.StringType,
"field_array_of_strings": types.ListType{
 ElemType: types.StringType,
},
"field_anonymous_struct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
},
},
"some_other_struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"field_any": types.ObjectType{},
},
},
"struct_optional_fields_some_struct_field_anonymous_struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"field_any": types.ObjectType{},
},
},
}