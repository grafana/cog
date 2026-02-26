package constraints

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type SomeStruct struct {
 Id types.Int64 `tfsdk:"id"`
MaybeId types.Int64 `tfsdk:"maybeId"`
Title types.String `tfsdk:"title"`
RefStruct RefStruct `tfsdk:"refStruct"`
 }

type RefStruct struct {
 Labels types.Map `tfsdk:"labels"`
Tags types.List `tfsdk:"tags"`
 }

var SpecAttributes = map[string]schema.Attribute{
"somestruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"id": types.Int64Type,
"maybeId": types.Int64Type,
"title": types.StringType,
"refStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"labels": types.MapType{
 ElemType: types.StringType,
},
"tags": types.ListType{
 ElemType: types.StringType,
},
},
,
},
"refstruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"labels": types.MapType{
 ElemType: types.StringType,
},
"tags": types.ListType{
 ElemType: types.StringType,
},
},
}