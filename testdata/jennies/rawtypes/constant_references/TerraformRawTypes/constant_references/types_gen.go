package constant_references

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)



type ParentStruct struct {
 MyEnum types.String `tfsdk:"myEnum"`
 }

type Struct struct {
 MyValue types.String `tfsdk:"myValue"`
MyEnum types.String `tfsdk:"myEnum"`
 }

type StructA struct {
 MyEnum types.String `tfsdk:"myEnum"`
Other types.String `tfsdk:"other"`
 }

type StructB struct {
 MyEnum types.String `tfsdk:"myEnum"`
MyValue types.String `tfsdk:"myValue"`
 }

var SpecAttributes = map[string]schema.Attribute{
"parent_struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"my_enum": types.StringType,
},
},
"struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"my_value": types.StringType,
"my_enum": types.StringType,
},
},
"struct_a": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"my_enum": types.StringType,
"other": types.StringType,
},
},
"struct_b": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"my_enum": types.StringType,
"my_value": types.StringType,
},
},
}