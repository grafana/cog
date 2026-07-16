package disjunctions_of_scalars_and_refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

type DisjunctionOfScalarsAndRefsModel = types.Object
var DisjunctionOfScalarsAndRefsType = StringOrBoolOrArrayOfStringOrMyRefAOrMyRefBType


type MyRefAModel struct {
	Foo types.String `tfsdk:"foo"`
}
var MyRefAType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"foo": types.StringType,
	},
}
var MyRefASchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"foo": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type MyRefBModel struct {
	Bar types.Int64 `tfsdk:"bar"`
}
var MyRefBType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"bar": types.Int64Type,
	},
}
var MyRefBSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"bar": schema.Int64Attribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type StringOrBoolOrArrayOfStringOrMyRefAOrMyRefBModel struct {
	String types.String `tfsdk:"string"`
	Bool types.Bool `tfsdk:"bool"`
	ArrayOfString types.List `tfsdk:"array_of_string"`
	MyRefA types.Object `tfsdk:"my_ref_a"`
	MyRefB types.Object `tfsdk:"my_ref_b"`
}
var StringOrBoolOrArrayOfStringOrMyRefAOrMyRefBType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"string": types.StringType,
		"bool": types.BoolType,
		"array_of_string": types.ListType{
	ElemType: types.StringType,
},
		"my_ref_a": MyRefAType,
		"my_ref_b": MyRefBType,
	},
}
var StringOrBoolOrArrayOfStringOrMyRefAOrMyRefBSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"string": schema.StringAttribute{
	Optional: true,
	Default: stringdefault.StaticString("a"),
	Computed: true,
},
	"bool": schema.BoolAttribute{
	Optional: true,
},
	"array_of_string": schema.ListAttribute{
	ElementType: types.StringType,
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
	"my_ref_a": schema.SingleNestedBlock{
		Attributes: MyRefASchema.Attributes,
		Blocks: MyRefASchema.Blocks,
		Validators: MyRefASchema.Validators,
	},
	"my_ref_b": schema.SingleNestedBlock{
		Attributes: MyRefBSchema.Attributes,
		Blocks: MyRefBSchema.Blocks,
		Validators: MyRefBSchema.Validators,
	},
},
}

