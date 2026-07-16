package constraints

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	int64validator "github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	regexp "regexp"
	listvalidator "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
)

type SomeStructModel struct {
	Id types.Int64 `tfsdk:"id"`
	MaybeId types.Int64 `tfsdk:"maybe_id"`
	GreaterThanZero types.Int64 `tfsdk:"greater_than_zero"`
	Negative types.Int64 `tfsdk:"negative"`
	Title types.String `tfsdk:"title"`
	Labels types.Map `tfsdk:"labels"`
	Tags types.List `tfsdk:"tags"`
	Regex types.String `tfsdk:"regex"`
	NegativeRegex types.String `tfsdk:"negative_regex"`
	MinMaxList types.List `tfsdk:"min_max_list"`
	UniqueList types.List `tfsdk:"unique_list"`
	FullConstraintList types.List `tfsdk:"full_constraint_list"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id": types.Int64Type,
		"maybe_id": types.Int64Type,
		"greater_than_zero": types.Int64Type,
		"negative": types.Int64Type,
		"title": types.StringType,
		"labels": types.MapType{
	ElemType: types.StringType,
},
		"tags": types.ListType{
	ElemType: types.StringType,
},
		"regex": types.StringType,
		"negative_regex": types.StringType,
		"min_max_list": types.ListType{
	ElemType: types.StringType,
},
		"unique_list": types.ListType{
	ElemType: types.StringType,
},
		"full_constraint_list": types.ListType{
	ElemType: types.Int64Type,
},
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
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
	Required: true,
},
	"tags": schema.ListAttribute{
	ElementType: types.StringType,
	Required: true,
},
	"regex": schema.StringAttribute{
	Required: true,
	Validators: []validator.String{
stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9_-]+$`), ""),
},
},
	"negative_regex": schema.StringAttribute{
	Required: true,
},
	"min_max_list": schema.ListAttribute{
	ElementType: types.StringType,
	Validators: []validator.List{
listvalidator.SizeAtLeast(1),
listvalidator.SizeAtMost(64),
},
	Required: true,
},
	"unique_list": schema.ListAttribute{
	ElementType: types.StringType,
	Validators: []validator.List{
listvalidator.UniqueValues(),
},
	Required: true,
},
	"full_constraint_list": schema.ListAttribute{
	ElementType: types.Int64Type,
	Validators: []validator.List{
listvalidator.SizeAtLeast(2),
listvalidator.SizeAtMost(10),
listvalidator.UniqueValues(),
},
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

