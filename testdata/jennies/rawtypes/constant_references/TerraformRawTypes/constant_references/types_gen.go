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
"parentstruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"myEnum": types.StringType,
},
},
"struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"myValue": types.StringType,
"myEnum": types.StringType,
},
},
"structa": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"myEnum": types.StringType,
"other": types.StringType,
},
},
"structb": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"myEnum": types.StringType,
"myValue": types.StringType,
},
},
}