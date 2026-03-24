package constraints

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	int64validator "github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

type SomeStruct struct {
 Id types.Int64 `tfsdk:"id"`
MaybeId types.Int64 `tfsdk:"maybeId"`
GreaterThanZero types.Int64 `tfsdk:"greaterThanZero"`
Negative types.Int64 `tfsdk:"negative"`
Title types.String `tfsdk:"title"`
Labels types.Map `tfsdk:"labels"`
Tags types.List `tfsdk:"tags"`
Regex types.String `tfsdk:"regex"`
NegativeRegex types.String `tfsdk:"negativeRegex"`
 }

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"id": schema.Int64Attribute{
 Required: true,
Validators: []validator.Int64{
int64validator.AtLeast(5),
int64validator.AtMost(9),
},

},

"maybe_id": schema.Int64Attribute{
 Optional: true,
Validators: []validator.Int64{
int64validator.AtLeast(5),
int64validator.AtMost(9),
},

},

"greater_than_zero": schema.Int64Attribute{
 Required: true,
Validators: []validator.Int64{
int64validator.AtMost(2),
},

},

"negative": schema.Int64Attribute{
 Required: true,
Validators: []validator.Int64{
int64validator.AtLeast(-19),
int64validator.AtMost(9),
},

},

"title": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.LengthAtLeast(1),
},

},

"labels": schema.MapAttribute{
 ElementType: types.StringType,
},

"tags": schema.ListAttribute{
 ElementType: types.StringType,
},

"regex": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
unknown},

},

"negative_regex": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
unknown},

},

},
},
}