package disjunctions_of_scalars_and_refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type DisjunctionOfScalarsAndRefs = StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB

type MyRefA struct {
 Foo types.String `tfsdk:"foo"`
 }

type MyRefB struct {
 Bar types.Int64 `tfsdk:"bar"`
 }

type StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB struct {
 String types.String `tfsdk:"String"`
Bool types.Bool `tfsdk:"Bool"`
ArrayOfString types.List `tfsdk:"ArrayOfString"`
MyRefA MyRefA `tfsdk:"MyRefA"`
MyRefB MyRefB `tfsdk:"MyRefB"`
 }

var SpecAttributes = map[string]schema.Attribute{
"disjunctionofscalarsandrefs": "myrefa": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"foo": types.StringType,
},
"myrefb": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"bar": types.Int64Type,
},
"stringorboolorarrayofstringormyrefaormyrefb": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"String": types.StringType,
"Bool": types.BoolType,
"ArrayOfString": types.ListType{
 ElemType: types.StringType,
},
"MyRefA": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"foo": types.StringType,
},
,
"MyRefB": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"bar": types.Int64Type,
},
,
},
}