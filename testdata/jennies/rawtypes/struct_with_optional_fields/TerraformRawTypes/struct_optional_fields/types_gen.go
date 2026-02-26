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
"somestruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"FieldRef": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
"FieldString": types.StringType,
"Operator": types.StringType,
"FieldArrayOfStrings": types.ListType{
 ElemType: types.StringType,
},
"FieldAnonymousStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
},
},
"someotherstruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
"structoptionalfieldssomestructfieldanonymousstruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
}