package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	booldefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	float32default "github.com/hashicorp/terraform-plugin-framework/resource/schema/float32default"
	int32default "github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
)

type SomeStructModel struct {
	FieldBool types.Bool `tfsdk:"field_bool"`
	FieldString types.String `tfsdk:"field_string"`
	FieldStringWithConstantValue types.String `tfsdk:"field_string_with_constant_value"`
	FieldFloat32 types.Float32 `tfsdk:"field_float32"`
	FieldInt32 types.Int32 `tfsdk:"field_int32"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_bool": types.BoolType,
		"field_string": types.StringType,
		"field_string_with_constant_value": types.StringType,
		"field_float32": types.Float32Type,
		"field_int32": types.Int32Type,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_bool": schema.BoolAttribute{
	Optional: true,
	Default: booldefault.StaticBool(true),
	Computed: true,
},
	"field_string": schema.StringAttribute{
	Optional: true,
	Default: stringdefault.StaticString("foo"),
	Computed: true,
},
	"field_string_with_constant_value": schema.StringAttribute{
	Required: true,
	Default: stringdefault.StaticString("auto"),
	Computed: true,
},
	"field_float32": schema.Float32Attribute{
	Optional: true,
	Default: float32default.StaticFloat32(42.42),
	Computed: true,
},
	"field_int32": schema.Int32Attribute{
	Optional: true,
	Default: int32default.StaticInt32(42),
	Computed: true,
},
},
	Blocks: map[string]schema.Block{
},
}

