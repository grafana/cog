package constraints

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"id": schema.Int64Attribute{
 Required: true,
},

"maybe_id": schema.Int64Attribute{
 Optional: true,
},

"title": schema.StringAttribute{
 Required: true,
},

"ref_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"labels": schema.MapAttribute{
 ElementType: types.StringType,
},

"tags": schema.ListAttribute{
 ElementType: types.StringType,
},

},
},

},
},
"ref_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"labels": schema.MapAttribute{
 ElementType: types.StringType,
},

"tags": schema.ListAttribute{
 ElementType: types.StringType,
},

},
},
}