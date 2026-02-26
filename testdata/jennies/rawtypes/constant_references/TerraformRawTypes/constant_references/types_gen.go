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
"enum": "parentstruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"myEnum": unknown,
},
"struct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"myValue": types.StringType,
"myEnum": unknown,
},
"structa": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"myEnum": unknown,
"other": unknown,
},
"structb": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"myEnum": unknown,
"myValue": types.StringType,
},
}