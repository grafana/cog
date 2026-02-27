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
"someStruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldRef": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
"fieldString": types.StringType,
"operator": types.StringType,
"fieldArrayOfStrings": types.ListType{
 ElemType: types.StringType,
},
"fieldAnonymousStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
},
},
"someOtherStruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
"structOptionalFieldsSomeStructFieldAnonymousStruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
}